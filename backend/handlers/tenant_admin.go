package handlers

import (
	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/middleware"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== TENANT ADMIN ====================

func TenantDashboard(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var totalListings, activeListings, totalSalesmen int64
	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ?", tenantID).Count(&totalListings)
	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ? AND status IN ?", tenantID,
		[]string{models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved}).Count(&activeListings)
	database.DB.Model(&models.User{}).Where("tenant_id = ? AND role = ? AND status = ?", tenantID, models.RoleSalesman, models.UserStatusActive).Count(&totalSalesmen)

	// Status breakdown
	statusCount := map[string]int64{}
	statuses := []string{models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved,
		models.ListingStatusRejected, models.ListingStatusSold, models.ListingStatusRented, models.ListingStatusInactive}
	for _, s := range statuses {
		var c int64
		database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ? AND status = ?", tenantID, s).Count(&c)
		statusCount[s] = c
	}

	var sub models.Subscription
	database.DB.Where("tenant_id = ?", tenantID).First(&sub)

	utils.OK(c, gin.H{
		"total_listings":  totalListings,
		"active_listings": activeListings,
		"total_salesmen":  totalSalesmen,
		"max_salesmen":    sub.MaxSalesmen,
		"status_breakdown": gin.H{
			"draft":    statusCount[models.ListingStatusDraft],
			"pending":  statusCount[models.ListingStatusPending],
			"approved": statusCount[models.ListingStatusApproved],
			"rejected": statusCount[models.ListingStatusRejected],
			"sold":     statusCount[models.ListingStatusSold],
			"rented":   statusCount[models.ListingStatusRented],
			"inactive": statusCount[models.ListingStatusInactive],
		},
		"plan": gin.H{
			"type":                     sub.PlanType,
			"max_salesmen":             sub.MaxSalesmen,
			"max_listings_per_salesman": sub.MaxListingsPerSalesman,
		},
	})
}

func TenantGetProfile(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var tenant models.Tenant
	if err := database.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		utils.NotFound(c, "RES_TENANT_NOT_FOUND", "Organisasi tidak ditemukan.")
		return
	}

	utils.OK(c, gin.H{
		"id":                tenant.ID,
		"organization_name": tenant.OrganizationName,
		"subdomain_slug":    tenant.SubdomainSlug,
		"logo_url":          tenant.LogoURL,
		"description":       tenant.Description,
		"phone":             tenant.Phone,
		"address":           tenant.Address,
		"status":            tenant.Status,
		"created_at":        tenant.CreatedAt,
	})
}

func TenantUpdateProfile(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	updates := map[string]interface{}{}
	if req.OrganizationName != nil {
		updates["organization_name"] = *req.OrganizationName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}

	if len(updates) > 0 {
		database.DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Updates(updates)
	}

	TenantGetProfile(c)
}

func TenantListSalesmen(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	p := utils.GetPagination(c)

	var salesmen []models.User
	query := database.DB.Where("tenant_id = ? AND role = ?", tenantID, models.RoleSalesman)

	var total int64
	query.Model(&models.User{}).Count(&total)
	query.Offset(p.Offset).Limit(p.PerPage).Find(&salesmen)

	type salesmanItem struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Email        string  `json:"email"`
		Phone        *string `json:"phone"`
		PhotoURL     *string `json:"photo_url"`
		Status       string  `json:"status"`
		ListingCount struct {
			Total  int `json:"total"`
			Active int `json:"active"`
		} `json:"listing_count"`
		CreatedAt string `json:"created_at"`
	}

	result := make([]salesmanItem, 0, len(salesmen))
	for _, s := range salesmen {
		var totalL, activeL int64
		database.DB.Model(&models.PropertyListing{}).Where("salesman_id = ?", s.ID).Count(&totalL)
		database.DB.Model(&models.PropertyListing{}).Where("salesman_id = ? AND status IN ?", s.ID,
			[]string{models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved}).Count(&activeL)

		item := salesmanItem{
			ID: s.ID.String(), Name: s.Name, Email: s.Email,
			Phone: s.Phone, PhotoURL: s.PhotoURL, Status: s.Status,
			ListingCount: struct {
				Total  int `json:"total"`
				Active int `json:"active"`
			}{Total: int(totalL), Active: int(activeL)},
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05+07:00"),
		}
		result = append(result, item)
	}

	utils.SuccessPaginated(c, 200, result, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func TenantAddSalesman(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.AddSalesmanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	// Check email unique
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Conflict(c, "AUTH_EMAIL_REGISTERED", "Email sudah terdaftar. Silakan gunakan email lain atau login.")
		return
	}

	// Check salesman limit
	var sub models.Subscription
	database.DB.Where("tenant_id = ?", tenantID).First(&sub)

	var currentCount int64
	database.DB.Model(&models.User{}).Where("tenant_id = ? AND role = ? AND status = ?", tenantID, models.RoleSalesman, models.UserStatusActive).Count(&currentCount)

	if int(currentCount) >= sub.MaxSalesmen {
		utils.Unprocessable(c, "BIZ_SALESMAN_LIMIT", "Jumlah salesman sudah mencapai batas. Upgrade ke Premium untuk menambah salesman.")
		return
	}

	hash, _ := utils.HashPassword(req.Password, AppConfig.BcryptCost)

	salesman := models.User{
		TenantID:     &tenantID,
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Phone:        &req.Phone,
		Role:         models.RoleSalesman,
		Status:       models.UserStatusActive,
	}

	if err := database.DB.Create(&salesman).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	utils.Created(c, gin.H{
		"id":         salesman.ID,
		"name":       salesman.Name,
		"email":      salesman.Email,
		"phone":      salesman.Phone,
		"role":       salesman.Role,
		"status":     salesman.Status,
		"created_at": salesman.CreatedAt,
	})
}

func TenantRemoveSalesman(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	salesmanID, _ := uuid.Parse(c.Param("id"))

	result := database.DB.Model(&models.User{}).
		Where("id = ? AND tenant_id = ? AND role = ?", salesmanID, tenantID, models.RoleSalesman).
		Update("status", models.UserStatusInactive)

	if result.RowsAffected == 0 {
		utils.NotFound(c, "RES_USER_NOT_FOUND", "Salesman tidak ditemukan.")
		return
	}

	utils.OK(c, gin.H{"message": "Salesman berhasil dinonaktifkan."})
}

func TenantListListings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	p := utils.GetPagination(c)
	statusFilter := c.Query("status")
	salesmanFilter := c.Query("salesman_id")

	query := database.DB.Model(&models.PropertyListing{}).
		Preload("Salesman").
		Where("tenant_id = ?", tenantID)

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}
	if salesmanFilter != "" {
		query = query.Where("salesman_id = ?", salesmanFilter)
	}

	query = query.Order("created_at DESC")

	var total int64
	query.Count(&total)

	var listings []models.PropertyListing
	query.Offset(p.Offset).Limit(p.PerPage).Find(&listings)

	type item struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Price        string `json:"price"`
		ListingType  string `json:"listing_type"`
		PropertyType string `json:"property_type"`
		Status       string `json:"status"`
		Salesman     struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"salesman"`
		CreatedAt string `json:"created_at"`
	}

	result := make([]item, 0, len(listings))
	for _, l := range listings {
		i := item{
			ID: l.ID.String(), Title: l.Title, Price: formatPrice(l.Price),
			ListingType: l.ListingType, PropertyType: l.PropertyType, Status: l.Status,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05+07:00"),
		}
		if l.Salesman != nil {
			i.Salesman.ID = l.Salesman.ID.String()
			i.Salesman.Name = l.Salesman.Name
		}
		result = append(result, i)
	}

	utils.SuccessPaginated(c, 200, result, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func TenantGetSubscription(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var sub models.Subscription
	database.DB.Where("tenant_id = ?", tenantID).First(&sub)

	var salesmenUsed int64
	database.DB.Model(&models.User{}).Where("tenant_id = ? AND role = ? AND status = ?", tenantID, models.RoleSalesman, models.UserStatusActive).Count(&salesmenUsed)

	var activeListings int64
	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ? AND status IN ?", tenantID,
		[]string{models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved}).Count(&activeListings)

	utils.OK(c, gin.H{
		"plan_type":                sub.PlanType,
		"max_salesmen":             sub.MaxSalesmen,
		"max_listings_per_salesman": sub.MaxListingsPerSalesman,
		"start_date": sub.StartDate,
		"end_date":   sub.EndDate,
		"usage": gin.H{
			"salesmen_used":        salesmenUsed,
			"salesmen_max":         sub.MaxSalesmen,
			"total_active_listings": activeListings,
		},
	})
}

func TenantRequestUpgrade(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.RequestUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field plan_type diperlukan (free, premium, atau disable).", nil)
		return
	}

	var sub models.Subscription
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&sub).Error; err != nil {
		utils.NotFound(c, "SUBSCRIPTION_NOT_FOUND", "Langganan tidak ditemukan.")
		return
	}

	// Check for any existing pending request
	if sub.PlanType == models.PlanPendingUpgrade ||
		sub.PlanType == models.PlanPendingFree ||
		sub.PlanType == models.PlanPendingDisable {
		utils.Conflict(c, "REQUEST_PENDING", "Anda sudah memiliki permintaan yang sedang diproses.")
		return
	}

	// Map request to target pending state
	var targetPlan string
	var message string
	switch req.PlanType {
	case "premium":
		if sub.PlanType == models.PlanPremium {
			utils.Conflict(c, "ALREADY_PREMIUM", "Anda sudah menggunakan paket Premium.")
			return
		}
		targetPlan = models.PlanPendingUpgrade
		message = "Permintaan upgrade ke Premium telah dikirim. Menunggu persetujuan Admin."
	case "free":
		if sub.PlanType == models.PlanFree {
			utils.Conflict(c, "ALREADY_FREE", "Anda sudah menggunakan paket Free.")
			return
		}
		targetPlan = models.PlanPendingFree
		message = "Permintaan downgrade ke Free telah dikirim. Menunggu persetujuan Admin."
	case "disable":
		targetPlan = models.PlanPendingDisable
		message = "Permintaan penonaktifan akun telah dikirim. Menunggu persetujuan Admin."
	}

	if err := database.DB.Model(&sub).Update("plan_type", targetPlan).Error; err != nil {
		utils.InternalError(c, "Gagal mengirim permintaan.")
		return
	}

	utils.OK(c, gin.H{
		"message": message,
		"status":  targetPlan,
	})
}
