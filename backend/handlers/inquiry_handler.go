package handlers

import (
	"property-hub-backend/database"
	"property-hub-backend/middleware"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== BUYER INQUIRIES ====================

func BuyerCreateInquiry(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		PropertyID string `json:"property_id" binding:"required"`
		Message    string `json:"message" binding:"required,min=1,max=2000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	propID, err := uuid.Parse(req.PropertyID)
	if err != nil {
		utils.NotFound(c, "RES_NOT_FOUND", "Properti tidak ditemukan.")
		return
	}

	// Verify property exists and is approved
	var prop models.PropertyListing
	if err := database.DB.Where("id = ? AND status = ?", propID, models.ListingStatusApproved).First(&prop).Error; err != nil {
		utils.NotFound(c, "RES_NOT_FOUND", "Properti tidak ditemukan.")
		return
	}

	inquiry := models.Inquiry{
		PropertyID: propID,
		BuyerID:    userID,
		Message:    &req.Message,
		Status:     models.InquiryStatusUnread,
	}
	if err := database.DB.Create(&inquiry).Error; err != nil {
		utils.InternalError(c, "Gagal mengirim pertanyaan.")
		return
	}

	utils.Created(c, gin.H{
		"id":          inquiry.ID,
		"property_id": inquiry.PropertyID,
		"message":     inquiry.Message,
		"status":      inquiry.Status,
		"created_at":  inquiry.CreatedAt,
	})
}

func BuyerListInquiries(c *gin.Context) {
	userID := middleware.GetUserID(c)
	p := utils.GetPagination(c)

	var inquiries []models.Inquiry
	query := database.DB.Preload("Property").Where("buyer_id = ?", userID).Order("created_at DESC")

	var total int64
	query.Model(&models.Inquiry{}).Count(&total)
	query.Offset(p.Offset).Limit(p.PerPage).Find(&inquiries)

	items := make([]gin.H, 0, len(inquiries))
	for _, i := range inquiries {
		item := gin.H{
			"id":          i.ID,
			"property_id": i.PropertyID,
			"message":     i.Message,
			"status":      i.Status,
			"created_at":  i.CreatedAt,
		}
		if i.Property != nil {
			item["property_title"] = i.Property.Title
		}
		items = append(items, item)
	}

	utils.SuccessPaginated(c, 200, items, utils.CalculateMeta(p.Page, p.PerPage, total))
}

// ==================== SALESMAN INQUIRIES ====================

func SalesmanListInquiries(c *gin.Context) {
	userID := middleware.GetUserID(c)
	p := utils.GetPagination(c)

	var inquiries []models.Inquiry
	query := database.DB.
		Preload("Property").
		Preload("Buyer").
		Joins("JOIN property_listings ON property_listings.id = inquiries.property_id").
		Where("property_listings.salesman_id = ?", userID).
		Order("inquiries.created_at DESC")

	if status := c.Query("status"); status != "" {
		query = query.Where("inquiries.status = ?", status)
	}

	var total int64
	query.Model(&models.Inquiry{}).Count(&total)
	query.Offset(p.Offset).Limit(p.PerPage).Find(&inquiries)

	items := make([]gin.H, 0, len(inquiries))
	for _, i := range inquiries {
		item := gin.H{
			"id":          i.ID,
			"property_id": i.PropertyID,
			"message":     i.Message,
			"status":      i.Status,
			"created_at":  i.CreatedAt,
		}
		if i.Property != nil {
			item["property_title"] = i.Property.Title
		}
		if i.Buyer != nil {
			item["buyer"] = gin.H{"id": i.Buyer.ID, "name": i.Buyer.Name, "email": i.Buyer.Email}
		}
		items = append(items, item)
	}

	utils.SuccessPaginated(c, 200, items, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func SalesmanUpdateInquiry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.NotFound(c, "RES_NOT_FOUND", "Inquiry tidak ditemukan.")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=read replied closed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Status tidak valid.", nil)
		return
	}

	result := database.DB.Model(&models.Inquiry{}).Where("id = ?", id).Update("status", req.Status)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "RES_NOT_FOUND", "Inquiry tidak ditemukan.")
		return
	}

	utils.OK(c, gin.H{"id": id, "status": req.Status, "message": "Status inquiry diperbarui."})
}

// ==================== TENANT INQUIRIES ====================

func TenantListInquiries(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	p := utils.GetPagination(c)

	var inquiries []models.Inquiry
	query := database.DB.
		Preload("Property").
		Preload("Buyer").
		Joins("JOIN property_listings ON property_listings.id = inquiries.property_id").
		Where("property_listings.tenant_id = ?", tenantID).
		Order("inquiries.created_at DESC")

	if status := c.Query("status"); status != "" {
		query = query.Where("inquiries.status = ?", status)
	}
	if propID := c.Query("property_id"); propID != "" {
		query = query.Where("inquiries.property_id = ?", propID)
	}

	var total int64
	query.Model(&models.Inquiry{}).Count(&total)
	query.Offset(p.Offset).Limit(p.PerPage).Find(&inquiries)

	items := make([]gin.H, 0, len(inquiries))
	for _, i := range inquiries {
		item := gin.H{
			"id":          i.ID,
			"property_id": i.PropertyID,
			"message":     i.Message,
			"status":      i.Status,
			"created_at":  i.CreatedAt,
		}
		if i.Property != nil {
			item["property_title"] = i.Property.Title
		}
		if i.Buyer != nil {
			item["buyer"] = gin.H{"id": i.Buyer.ID, "name": i.Buyer.Name}
		}
		items = append(items, item)
	}

	utils.SuccessPaginated(c, 200, items, utils.CalculateMeta(p.Page, p.PerPage, total))
}

// ==================== PROPERTY VIEW LOG ====================

func LogPropertyView(c *gin.Context) {
	propID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return // silently ignore invalid IDs
	}

	var userID *uuid.UUID
	if uid, exists := c.Get("userID"); exists {
		if u, ok := uid.(uuid.UUID); ok {
			userID = &u
		}
	}

	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")

	view := models.PropertyView{
		PropertyID: propID,
		UserID:     userID,
		IPAddress:  &ip,
		UserAgent:  &ua,
	}
	database.DB.Create(&view)

	c.JSON(200, gin.H{"success": true})
}
