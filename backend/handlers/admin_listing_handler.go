package handlers

import (
	"time"

	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/middleware"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminListAllListings lists all listings with optional tenant filter
func AdminListAllListings(c *gin.Context) {
	p := utils.GetPagination(c)

	var listings []models.PropertyListing
	query := database.DB.
		Preload("Photos", "sort_order = ?", 0).
		Preload("Salesman").
		Preload("Tenant").
		Order("created_at DESC")

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.PropertyListing{}).Count(&total)
	query.Offset(p.Offset).Limit(p.PerPage).Find(&listings)

	cards := make([]dto.PropertyCardResponse, 0, len(listings))
	for _, l := range listings {
		cards = append(cards, toPropertyCard(l))
	}

	utils.SuccessPaginated(c, 200, cards, utils.CalculateMeta(p.Page, p.PerPage, total))
}

// AdminChangePlanByID changes subscription by plan_id (new v2 endpoint)
func AdminChangePlanByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.NotFound(c, "RES_TENANT_NOT_FOUND", "ID tenant tidak valid.")
		return
	}

	var req struct {
		PlanID string `json:"plan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field plan_id diperlukan.", nil)
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		utils.NotFound(c, "RES_NOT_FOUND", "Paket tidak ditemukan.")
		return
	}

	// Verify plan exists
	var plan models.SubscriptionPlan
	if err := database.DB.First(&plan, "id = ?", planID).Error; err != nil {
		utils.NotFound(c, "RES_NOT_FOUND", "Paket tidak ditemukan.")
		return
	}

	// Verify subscription exists
	var sub models.Subscription
	if err := database.DB.Where("tenant_id = ?", id).First(&sub).Error; err != nil {
		utils.NotFound(c, "RES_SUBSCRIPTION_NOT_FOUND", "Subscription tidak ditemukan.")
		return
	}

	now := time.Now()
	oldPlan := sub.PlanType

	result := database.DB.Model(&sub).Updates(map[string]interface{}{
		"plan_type":                 plan.Slug,
		"max_salesmen":              plan.MaxSalesmen,
		"max_listings_per_salesman": plan.MaxListingsPerSalesman,
		"updated_at":                &now,
	})
	if result.Error != nil {
		utils.InternalError(c, "Gagal mengubah paket.")
		return
	}

	// Audit
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	oldVals := models.JSONMap{"plan_type": oldPlan}
	newVals := models.JSONMap{"plan_type": plan.Slug}
	database.DB.Create(&models.AuditLog{
		UserID:     &userID,
		UserRole:   role,
		Action:     models.AuditActionUpdate,
		EntityType: "subscription",
		EntityID:   sub.ID.String(),
		OldValues:  &oldVals,
		NewValues:  &newVals,
	})

	utils.OK(c, gin.H{
		"tenant_id":                 id,
		"plan_type":                 plan.Slug,
		"plan_name":                 plan.Name,
		"max_salesmen":              plan.MaxSalesmen,
		"max_listings_per_salesman": plan.MaxListingsPerSalesman,
		"message":                   "Paket berhasil diubah ke " + plan.Name + ".",
	})
}
