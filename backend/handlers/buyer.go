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

// ==================== BUYER ====================

func BuyerListSaved(c *gin.Context) {
	p := utils.GetPagination(c)
	buyerID := middleware.GetUserID(c)

	var saved []models.SavedProperty
	query := database.DB.
		Preload("Listing.Photos", "sort_order = ?", 0).
		Preload("Listing.Salesman").
		Preload("Listing.Tenant").
		Where("buyer_id = ?", buyerID).
		Order("created_at DESC")

	var total int64
	query.Model(&models.SavedProperty{}).Count(&total)

	if err := query.Offset(p.Offset).Limit(p.PerPage).Find(&saved).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	cards := make([]dto.PropertyCardResponse, 0, len(saved))
	for _, s := range saved {
		if s.Listing != nil {
			card := toPropertyCard(*s.Listing)
			cards = append(cards, card)
		}
	}

	utils.SuccessPaginated(c, 200, cards, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func BuyerSaveProperty(c *gin.Context) {
	buyerID := middleware.GetUserID(c)
	propertyID, err := uuid.Parse(c.Param("propertyId"))
	if err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	// Verify listing exists and is approved
	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND status = ?", propertyID, models.ListingStatusApproved).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	// Check duplicate
	var existing models.SavedProperty
	if err := database.DB.Where("buyer_id = ? AND listing_id = ?", buyerID, propertyID).First(&existing).Error; err == nil {
		utils.Conflict(c, "RES_ALREADY_SAVED", "Properti sudah ada di daftar favorit Anda.")
		return
	}

	saved := models.SavedProperty{
		BuyerID:   buyerID,
		ListingID: propertyID,
	}

	if err := database.DB.Create(&saved).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	utils.Created(c, gin.H{"message": "Properti berhasil disimpan ke favorit."})
}

func BuyerRemoveSaved(c *gin.Context) {
	buyerID := middleware.GetUserID(c)
	propertyID, err := uuid.Parse(c.Param("propertyId"))
	if err != nil {
		utils.NotFound(c, "RES_NOT_FOUND", "Data tidak ditemukan.")
		return
	}

	result := database.DB.Where("buyer_id = ? AND listing_id = ?", buyerID, propertyID).Delete(&models.SavedProperty{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "RES_NOT_FOUND", "Data tidak ditemukan.")
		return
	}

	utils.OK(c, gin.H{"message": "Properti berhasil dihapus dari favorit."})
}
