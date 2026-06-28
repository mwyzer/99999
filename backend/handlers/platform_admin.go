package handlers

import (
	"fmt"
	"time"

	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/middleware"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== PLATFORM ADMIN ====================

func AdminDashboard(c *gin.Context) {
	var totalTenants, activeTenants, suspendedTenants, totalUsers, totalListings, pendingReviews, pendingUpgrades int64

	database.DB.Model(&models.Tenant{}).Count(&totalTenants)
	database.DB.Model(&models.Tenant{}).Where("status = ?", models.TenantStatusActive).Count(&activeTenants)
	database.DB.Model(&models.Tenant{}).Where("status = ?", models.TenantStatusSuspended).Count(&suspendedTenants)
	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.PropertyListing{}).Count(&totalListings)
	database.DB.Model(&models.PropertyListing{}).Where("status = ?", models.ListingStatusPending).Count(&pendingReviews)
	database.DB.Model(&models.Subscription{}).Where("plan_type IN ?",
		[]string{models.PlanPendingUpgrade, models.PlanPendingFree, models.PlanPendingDisable}).Count(&pendingUpgrades)

	statusCount := map[string]int64{}
	statuses := []string{models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved,
		models.ListingStatusRejected, models.ListingStatusSold, models.ListingStatusRented, models.ListingStatusInactive}
	for _, s := range statuses {
		var c int64
		database.DB.Model(&models.PropertyListing{}).Where("status = ?", s).Count(&c)
		statusCount[s] = c
	}

	utils.OK(c, gin.H{
		"total_tenants":     totalTenants,
		"active_tenants":    activeTenants,
		"suspended_tenants": suspendedTenants,
		"pending_upgrades":  pendingUpgrades,
		"total_users":       totalUsers,
		"total_listings":    totalListings,
		"pending_reviews":   pendingReviews,
		"listings_by_status": gin.H{
			"draft":    statusCount[models.ListingStatusDraft],
			"pending":  statusCount[models.ListingStatusPending],
			"approved": statusCount[models.ListingStatusApproved],
			"rejected": statusCount[models.ListingStatusRejected],
			"sold":     statusCount[models.ListingStatusSold],
			"rented":   statusCount[models.ListingStatusRented],
			"inactive": statusCount[models.ListingStatusInactive],
		},
	})
}

func AdminListTenants(c *gin.Context) {
	p := utils.GetPagination(c)
	statusFilter := c.Query("status")

	query := database.DB.Model(&models.Tenant{}).
		Preload("Subscription")

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var total int64
	query.Count(&total)

	var tenants []models.Tenant
	query.Offset(p.Offset).Limit(p.PerPage).Order("created_at DESC").Find(&tenants)

	type item struct {
		ID               string `json:"id"`
		OrganizationName string `json:"organization_name"`
		SubdomainSlug    string `json:"subdomain_slug"`
		Phone            *string `json:"phone"`
		Status           string `json:"status"`
		PlanType         string `json:"plan_type"`
		TotalListings    int64  `json:"total_listings"`
		TotalUsers       int64  `json:"total_users"`
		CreatedAt        string `json:"created_at"`
	}

	result := make([]item, 0, len(tenants))
	for _, t := range tenants {
		var listings, users int64
		database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ?", t.ID).Count(&listings)
		database.DB.Model(&models.User{}).Where("tenant_id = ?", t.ID).Count(&users)

		planType := "free"
		if t.Subscription != nil {
			planType = t.Subscription.PlanType
		}

		result = append(result, item{
			ID: t.ID.String(), OrganizationName: t.OrganizationName,
			SubdomainSlug: t.SubdomainSlug, Phone: t.Phone, Status: t.Status,
			PlanType: planType, TotalListings: listings, TotalUsers: users,
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05+07:00"),
		})
	}

	utils.SuccessPaginated(c, 200, result, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func AdminCreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	// Check subdomain unique
	var existing models.Tenant
	if err := database.DB.Where("subdomain_slug = ?", req.SubdomainSlug).First(&existing).Error; err == nil {
		utils.Conflict(c, "RES_ALREADY_EXISTS", "Subdomain sudah digunakan.")
		return
	}

	// Check email unique (search the User table, not Tenant)
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.AdminEmail).First(&existingUser).Error; err == nil {
		utils.Conflict(c, "AUTH_EMAIL_REGISTERED", "Email admin sudah terdaftar.")
		return
	}

	planType := req.PlanType
	if planType == "" {
		planType = models.PlanFree
	}

	maxSalesmen := 5
	maxListings := 5
	if planType == models.PlanPremium {
		maxSalesmen = 999999
		maxListings = 999999
	}

	// Transaction
	tx := database.DB.Begin()

	tenant := models.Tenant{
		OrganizationName: req.OrganizationName,
		SubdomainSlug:    req.SubdomainSlug,
		Status:           models.TenantStatusActive,
	}
	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal membuat tenant.")
		return
	}

	hash, _ := utils.HashPassword(req.AdminPassword, AppConfig.BcryptCost)
	adminUser := models.User{
		TenantID:     &tenant.ID,
		Email:        req.AdminEmail,
		PasswordHash: hash,
		Name:         req.AdminName,
		Phone:        &req.AdminPhone,
		Role:         models.RoleTenantAdmin,
		Status:       models.UserStatusActive,
	}
	if err := tx.Create(&adminUser).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal membuat admin tenant.")
		return
	}

	sub := models.Subscription{
		TenantID:               tenant.ID,
		PlanType:               planType,
		MaxSalesmen:            maxSalesmen,
		MaxListingsPerSalesman: maxListings,
	}
	if err := tx.Create(&sub).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal membuat subscription.")
		return
	}

	tx.Commit()

	utils.Created(c, gin.H{
		"tenant": gin.H{
			"id":                tenant.ID,
			"organization_name": tenant.OrganizationName,
			"subdomain_slug":    tenant.SubdomainSlug,
			"status":            tenant.Status,
		},
		"admin": gin.H{
			"id":    adminUser.ID,
			"name":  adminUser.Name,
			"email": adminUser.Email,
			"role":  adminUser.Role,
		},
		"subscription": gin.H{
			"plan_type":                sub.PlanType,
			"max_salesmen":             sub.MaxSalesmen,
			"max_listings_per_salesman": sub.MaxListingsPerSalesman,
		},
	})
}

func AdminGetTenant(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	var tenant models.Tenant
	if err := database.DB.Preload("Subscription").First(&tenant, "id = ?", id).Error; err != nil {
		utils.NotFound(c, "RES_TENANT_NOT_FOUND", "Organisasi tidak ditemukan.")
		return
	}

	var totalUsers, totalSalesmen, totalListings, pendingListings int64
	database.DB.Model(&models.User{}).Where("tenant_id = ?", tenant.ID).Count(&totalUsers)
	database.DB.Model(&models.User{}).Where("tenant_id = ? AND role = ?", tenant.ID, models.RoleSalesman).Count(&totalSalesmen)
	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ?", tenant.ID).Count(&totalListings)
	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ? AND status = ?", tenant.ID, models.ListingStatusPending).Count(&pendingListings)

	data := gin.H{
		"id":                tenant.ID,
		"organization_name": tenant.OrganizationName,
		"subdomain_slug":    tenant.SubdomainSlug,
		"logo_url":          tenant.LogoURL,
		"description":       tenant.Description,
		"phone":             tenant.Phone,
		"address":           tenant.Address,
		"status":            tenant.Status,
		"created_at":        tenant.CreatedAt,
		"stats": gin.H{
			"total_users":      totalUsers,
			"total_salesmen":   totalSalesmen,
			"total_listings":   totalListings,
			"pending_listings": pendingListings,
		},
	}

	if tenant.Subscription != nil {
		data["subscription"] = gin.H{
			"plan_type":                tenant.Subscription.PlanType,
			"max_salesmen":             tenant.Subscription.MaxSalesmen,
			"max_listings_per_salesman": tenant.Subscription.MaxListingsPerSalesman,
			"start_date": tenant.Subscription.StartDate,
		}
	}

	utils.OK(c, data)
}

func AdminSuspendTenant(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	result := database.DB.Model(&models.Tenant{}).Where("id = ? AND status = ?", id, models.TenantStatusActive).Update("status", models.TenantStatusSuspended)
	if result.RowsAffected == 0 {
		utils.Unprocessable(c, "BIZ_INVALID_STATUS_TRANSITION", "Tenant tidak dapat dinonaktifkan.")
		return
	}

	utils.OK(c, gin.H{
		"id":      id,
		"status":  models.TenantStatusSuspended,
		"message": "Tenant berhasil dinonaktifkan. Semua pengguna tidak dapat login.",
	})
}

func AdminActivateTenant(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	result := database.DB.Model(&models.Tenant{}).Where("id = ? AND status = ?", id, models.TenantStatusSuspended).Update("status", models.TenantStatusActive)
	if result.RowsAffected == 0 {
		utils.Unprocessable(c, "BIZ_INVALID_STATUS_TRANSITION", "Tenant tidak dapat diaktifkan.")
		return
	}

	utils.OK(c, gin.H{
		"id":      id,
		"status":  models.TenantStatusActive,
		"message": "Tenant berhasil diaktifkan kembali.",
	})
}

func AdminDeleteTenant(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))

	var tenant models.Tenant
	if err := database.DB.First(&tenant, "id = ?", id).Error; err != nil {
		utils.NotFound(c, "RES_TENANT_NOT_FOUND", "Tenant tidak ditemukan.")
		return
	}

	tx := database.DB.Begin()

	// 1. Delete listings first (FK listings→salesman is ON DELETE RESTRICT)
	if err := tx.Unscoped().Where("tenant_id = ?", id).Delete(&models.PropertyListing{}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus listing.")
		return
	}

	// 2. Delete users (salesmen + tenant admin) — soft delete
	if err := tx.Where("tenant_id = ?", id).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus pengguna.")
		return
	}

	// 3. Delete junction / dependent rows (GORM AutoMigrate uses NO ACTION FKs, not CASCADE)
	if err := tx.Where("tenant_id = ?", id).Delete(&models.TenantUser{}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus relasi tenant-user.")
		return
	}
	if err := tx.Where("tenant_id = ?", id).Delete(&models.Subscription{}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus subscription.")
		return
	}
	if err := tx.Where("tenant_id = ?", id).Delete(&models.TenantSubscription{}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus tenant subscription.")
		return
	}
	if err := tx.Where("tenant_id = ?", id).Delete(&models.AuditLog{}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus audit log.")
		return
	}

	// 4. Soft-delete tenant (sets deleted_at, recoverable)
	if err := tx.Delete(&tenant).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Gagal menghapus tenant.")
		return
	}

	tx.Commit()

	utils.OK(c, gin.H{
		"id":      id,
		"message": "Tenant dan semua data terkait berhasil dihapus.",
	})
}

func AdminChangePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.NotFound(c, "RES_TENANT_NOT_FOUND", "ID tenant tidak valid.")
		return
	}

	var req dto.ChangePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	// Verify tenant exists and is active
	var tenant models.Tenant
	if err := database.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		utils.NotFound(c, "RES_TENANT_NOT_FOUND", "Tenant tidak ditemukan.")
		return
	}

	// Verify subscription exists
	var sub models.Subscription
	if err := database.DB.Where("tenant_id = ?", id).First(&sub).Error; err != nil {
		utils.NotFound(c, "RES_SUBSCRIPTION_NOT_FOUND", "Subscription tenant tidak ditemukan.")
		return
	}

	// Don't allow change if already on the same plan
	if sub.PlanType == req.PlanType {
		utils.Unprocessable(c, "BIZ_SAME_PLAN", fmt.Sprintf("Tenant sudah menggunakan paket %s.", req.PlanType))
		return
	}

	maxSalesmen := 5
	maxListings := 5
	if req.PlanType == models.PlanPremium {
		maxSalesmen = 999999
		maxListings = 999999
	}

	now := time.Now()
	oldPlan := sub.PlanType

	result := database.DB.Model(&sub).Updates(map[string]interface{}{
		"plan_type":                 req.PlanType,
		"max_salesmen":              maxSalesmen,
		"max_listings_per_salesman": maxListings,
		"updated_at":                &now,
	})
	if result.Error != nil {
		utils.InternalError(c, "Gagal mengubah paket tenant. Silakan coba lagi.")
		return
	}
	if result.RowsAffected == 0 {
		utils.InternalError(c, "Gagal mengubah paket tenant. Subscription tidak ditemukan.")
		return
	}

	// Audit log
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	entityID := sub.ID.String()
	oldVals := models.JSONMap{"plan_type": oldPlan}
	newVals := models.JSONMap{"plan_type": req.PlanType}
	database.DB.Create(&models.AuditLog{
		UserID:     &userID,
		UserRole:   role,
		Action:     models.AuditActionUpdate,
		EntityType: "subscription",
		EntityID:   entityID,
		OldValues:  &oldVals,
		NewValues:  &newVals,
	})

	utils.OK(c, gin.H{
		"tenant_id":                 id,
		"plan_type":                 req.PlanType,
		"max_salesmen":              maxSalesmen,
		"max_listings_per_salesman": maxListings,
		"message":                   "Paket tenant berhasil diubah dari " + oldPlan + " ke " + req.PlanType + ".",
	})
}

func AdminListPending(c *gin.Context) {
	p := utils.GetPagination(c)

	var listings []models.PropertyListing
	query := database.DB.
		Preload("Photos", "sort_order = ?", 0).
		Preload("Salesman").
		Preload("Tenant").
		Where("status = ?", models.ListingStatusPending).
		Order("created_at ASC")

	var total int64
	query.Model(&models.PropertyListing{}).Count(&total)
	query.Offset(p.Offset).Limit(p.PerPage).Find(&listings)

	type item struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Price        string `json:"price"`
		ListingType  string `json:"listing_type"`
		PropertyType string `json:"property_type"`
		SourceType   string `json:"source_type"`
		City         *string `json:"city"`
		Status       string `json:"status"`
		Tenant       struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"tenant"`
		Salesman struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"salesman"`
		Photos    []gin.H `json:"photos"`
		CreatedAt string  `json:"created_at"`
	}

	result := make([]item, 0, len(listings))
	for _, l := range listings {
		i := item{
			ID: l.ID.String(), Title: l.Title, Price: formatPrice(l.Price),
			ListingType: l.ListingType, PropertyType: l.PropertyType,
			SourceType: l.SourceType, City: l.City, Status: l.Status,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05+07:00"),
		}
		if l.Tenant != nil {
			i.Tenant.ID = l.Tenant.ID.String()
			i.Tenant.Name = l.Tenant.OrganizationName
		}
		if l.Salesman != nil {
			i.Salesman.ID = l.Salesman.ID.String()
			i.Salesman.Name = l.Salesman.Name
		}
		i.Photos = make([]gin.H, 0)
		for _, ph := range l.Photos {
			if ph.ThumbnailURL != nil {
				i.Photos = append(i.Photos, gin.H{"thumbnail_url": *ph.ThumbnailURL})
			}
		}
		result = append(result, i)
	}

	utils.SuccessPaginated(c, 200, result, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func AdminApproveListing(c *gin.Context) {
	adminID := c.MustGet("userID").(uuid.UUID)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND status = ?", listingID, models.ListingStatusPending).First(&listing).Error; err != nil {
		utils.Unprocessable(c, "BIZ_INVALID_STATUS_TRANSITION", "Hanya listing dengan status pending yang dapat disetujui.")
		return
	}

	now := time.Now()
	database.DB.Model(&listing).Updates(map[string]interface{}{
		"status":      models.ListingStatusApproved,
		"approved_by": adminID,
		"approved_at": now,
	})

	utils.OK(c, gin.H{
		"id":          listing.ID,
		"status":      models.ListingStatusApproved,
		"approved_by": adminID,
		"approved_at": now,
		"message":     "Listing berhasil disetujui. Sekarang tampil di halaman publik.",
	})
}

func AdminRejectListing(c *gin.Context) {
	listingID, _ := uuid.Parse(c.Param("id"))

	var req dto.RejectListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND status = ?", listingID, models.ListingStatusPending).First(&listing).Error; err != nil {
		utils.Unprocessable(c, "BIZ_INVALID_STATUS_TRANSITION", "Hanya listing dengan status pending yang dapat ditolak.")
		return
	}

	database.DB.Model(&listing).Updates(map[string]interface{}{
		"status":        models.ListingStatusRejected,
		"reject_reason": req.Reason,
	})

	utils.OK(c, gin.H{
		"id":            listing.ID,
		"status":        models.ListingStatusRejected,
		"reject_reason": req.Reason,
		"message":       "Listing berhasil ditolak. Salesman dapat mengedit dan mengajukan ulang.",
	})
}

func AdminAuditLogs(c *gin.Context) {
	p := utils.GetPagination(c)

	query := database.DB.Model(&models.AuditLog{}).Preload("User")

	if uid := c.Query("user_id"); uid != "" {
		query = query.Where("user_id = ?", uid)
	}
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if etype := c.Query("entity_type"); etype != "" {
		query = query.Where("entity_type = ?", etype)
	}
	if from := c.Query("from"); from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		query = query.Where("created_at <= ?", to+"T23:59:59")
	}

	query = query.Order("created_at DESC")

	var total int64
	query.Count(&total)

	var logs []models.AuditLog
	query.Offset(p.Offset).Limit(p.PerPage).Find(&logs)

	type logItem struct {
		ID         string          `json:"id"`
		User       *gin.H          `json:"user"`
		Action     string          `json:"action"`
		EntityType string          `json:"entity_type"`
		EntityID   string          `json:"entity_id"`
		OldValues  *models.JSONMap `json:"old_values"`
		NewValues  *models.JSONMap `json:"new_values"`
		IPAddress  *string         `json:"ip_address"`
		CreatedAt  string          `json:"created_at"`
	}

	result := make([]logItem, 0, len(logs))
	for _, l := range logs {
		item := logItem{
			ID: l.ID.String(), Action: l.Action,
			EntityType: l.EntityType, EntityID: l.EntityID,
			OldValues: l.OldValues, NewValues: l.NewValues,
			IPAddress: l.IPAddress,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05+07:00"),
		}
		if l.User != nil {
			item.User = &gin.H{
				"id":   l.User.ID,
				"name": l.User.Name,
				"role": l.User.Role,
			}
		}
		result = append(result, item)
	}

	utils.SuccessPaginated(c, 200, result, utils.CalculateMeta(p.Page, p.PerPage, total))
}
