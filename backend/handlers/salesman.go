package handlers

import (
	"fmt"
	"log"
	"path/filepath"
	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/middleware"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== SALESMAN ====================

func SalesmanDashboard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	var listings []models.PropertyListing
	query := database.DB.Model(&models.PropertyListing{})

	if role == models.RoleSalesman {
		query = query.Where("salesman_id = ?", userID)
	} else {
		tenantID := middleware.GetTenantID(c)
		query = query.Where("tenant_id = ?", tenantID)
	}

	query.Find(&listings)

	// Aggregate
	total := len(listings)
	statusCount := map[string]int{}
	activeCount := 0
	for _, l := range listings {
		statusCount[l.Status]++
		if l.IsQuotaCounted() {
			activeCount++
		}
	}

	// Quota
	var subscription models.Subscription
	tenantID := middleware.GetTenantID(c)
	if tenantID != uuid.Nil {
		database.DB.Where("tenant_id = ?", tenantID).First(&subscription)
	}

	utils.OK(c, gin.H{
		"total_listings": total,
		"status_breakdown": gin.H{
			"draft":    statusCount[models.ListingStatusDraft],
			"pending":  statusCount[models.ListingStatusPending],
			"approved": statusCount[models.ListingStatusApproved],
			"rejected": statusCount[models.ListingStatusRejected],
			"sold":     statusCount[models.ListingStatusSold],
			"rented":   statusCount[models.ListingStatusRented],
			"inactive": statusCount[models.ListingStatusInactive],
		},
		"active_count": activeCount,
		"quota": gin.H{
			"used":      activeCount,
			"max":       subscription.MaxListingsPerSalesman,
			"remaining": subscription.MaxListingsPerSalesman - activeCount,
			"plan_type": subscription.PlanType,
		},
	})
}

func SalesmanListListings(c *gin.Context) {
	p := utils.GetPagination(c)
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	statusFilter := c.Query("status")

	query := database.DB.Model(&models.PropertyListing{}).
		Preload("Photos", "sort_order = ?", 0)

	if role == models.RoleSalesman {
		query = query.Where("salesman_id = ?", userID)
	} else {
		tenantID := middleware.GetTenantID(c)
		query = query.Where("tenant_id = ?", tenantID)
	}

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	query = query.Order("created_at DESC")

	var total int64
	query.Count(&total)

	var listings []models.PropertyListing
	if err := query.Offset(p.Offset).Limit(p.PerPage).Find(&listings).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	type listingBrief struct {
		ID           string  `json:"id"`
		Title        string  `json:"title"`
		Price        string  `json:"price"`
		ListingType  string  `json:"listing_type"`
		PropertyType string  `json:"property_type"`
		SourceType   string  `json:"source_type"`
		City         *string `json:"city"`
		Status       string  `json:"status"`
		MainPhotoURL *string `json:"main_photo_url"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    *string `json:"updated_at"`
	}

	result := make([]listingBrief, 0, len(listings))
	for _, l := range listings {
		lb := listingBrief{
			ID: l.ID.String(), Title: l.Title, Price: formatPrice(l.Price),
			ListingType: l.ListingType, PropertyType: l.PropertyType,
			SourceType: l.SourceType, City: l.City, Status: l.Status,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05+07:00"),
		}
		if l.UpdatedAt != nil {
			s := l.UpdatedAt.Format("2006-01-02T15:04:05+07:00")
			lb.UpdatedAt = &s
		}
		if len(l.Photos) > 0 {
			lb.MainPhotoURL = l.Photos[0].MediumURL
		}
		result = append(result, lb)
	}

	utils.SuccessPaginated(c, 200, result, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func SalesmanCreateListing(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req dto.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	// Quota check
	if !checkQuota(userID, tenantID) {
		utils.Unprocessable(c, "BIZ_QUOTA_EXCEEDED", "Kuota listing Anda sudah penuh. Upgrade ke Premium untuk listing unlimited.")
		return
	}

	sourceType := req.SourceType
	if sourceType == "" {
		sourceType = models.SourceTypeRegular
	}

	listing := models.PropertyListing{
		TenantID:        tenantID,
		SalesmanID:      userID,
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		ListingType:     req.ListingType,
		PropertyType:    req.PropertyType,
		SourceType:      sourceType,
		Address:         req.Address,
		City:            req.City,
		Province:        req.Province,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		LandArea:        req.LandArea,
		BuildingArea:    req.BuildingArea,
		Bedrooms:        req.Bedrooms,
		Bathrooms:       req.Bathrooms,
		Floors:          req.Floors,
		CertificateType: req.CertificateType,
		Facilities:      models.Facilities(req.Facilities),
		Status:          models.ListingStatusDraft,
	}

	if err := database.DB.Create(&listing).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	utils.Created(c, gin.H{
		"id":            listing.ID,
		"title":         listing.Title,
		"price":         formatPrice(listing.Price),
		"listing_type":  listing.ListingType,
		"property_type": listing.PropertyType,
		"source_type":   listing.SourceType,
		"city":          listing.City,
		"status":        listing.Status,
		"created_at":    listing.CreatedAt,
	})
}

func SalesmanGetListing(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	query := database.DB.Preload("Photos", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") })

	if role == models.RoleSalesman {
		query = query.Where("id = ? AND salesman_id = ?", listingID, userID)
	} else {
		tenantID := middleware.GetTenantID(c)
		query = query.Where("id = ? AND tenant_id = ?", listingID, tenantID)
	}

	if err := query.First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	utils.OK(c, toPropertyDetail(listing))
}

func SalesmanUpdateListing(c *gin.Context) {
	userID := middleware.GetUserID(c)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ?", listingID, userID).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	if !listing.CanBeEdited() {
		utils.Unprocessable(c, "BIZ_LISTING_NOT_EDITABLE", "Listing dengan status "+listing.Status+" tidak dapat diedit. Hanya listing dengan status draft atau rejected yang dapat diedit.")
		return
	}

	var req dto.UpdateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	updates := buildListingUpdates(req)
	if len(updates) > 0 {
		database.DB.Model(&listing).Updates(updates)
	}

	utils.OK(c, gin.H{
		"id":         listing.ID,
		"title":      listing.Title,
		"price":      formatPrice(listing.Price),
		"status":     listing.Status,
		"updated_at": listing.UpdatedAt,
	})
}

func SalesmanDeleteListing(c *gin.Context) {
	userID := middleware.GetUserID(c)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ?", listingID, userID).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	if !listing.CanBeDeleted() {
		utils.Unprocessable(c, "BIZ_LISTING_NOT_DELETABLE", "Listing dengan status "+listing.Status+" tidak dapat dihapus.")
		return
	}

	database.DB.Delete(&listing)
	utils.OK(c, gin.H{"message": "Listing berhasil dihapus."})
}

func SalesmanSubmitListing(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ?", listingID, userID).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	if !listing.CanBeSubmitted() {
		utils.Unprocessable(c, "BIZ_LISTING_NOT_SUBMITTABLE", "Hanya listing dengan status draft atau rejected yang dapat diajukan.")
		return
	}

	// Quota check (re-submit after reject doesn't increase count, but we still verify)
	if listing.Status == models.ListingStatusDraft {
		if !checkQuota(userID, tenantID) {
			utils.Unprocessable(c, "BIZ_QUOTA_EXCEEDED", "Kuota listing Anda sudah penuh. Upgrade ke Premium untuk listing unlimited.")
			return
		}
	}

	database.DB.Model(&listing).Update("status", models.ListingStatusPending)

	utils.OK(c, gin.H{
		"id":      listing.ID,
		"status":  models.ListingStatusPending,
		"message": "Listing berhasil diajukan untuk review.",
	})
}

func SalesmanDeactivateListing(c *gin.Context) {
	userID := middleware.GetUserID(c)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ? AND status = ?", listingID, userID, models.ListingStatusApproved).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	database.DB.Model(&listing).Update("status", models.ListingStatusInactive)

	utils.OK(c, gin.H{
		"id":      listing.ID,
		"status":  models.ListingStatusInactive,
		"message": "Listing berhasil dinonaktifkan.",
	})
}

func SalesmanMarkSold(c *gin.Context) {
	markListingStatus(c, models.ListingStatusSold, "terjual")
}

func SalesmanMarkRented(c *gin.Context) {
	markListingStatus(c, models.ListingStatusRented, "tersewa")
}

func SalesmanUploadPhotos(c *gin.Context) {
	userID := middleware.GetUserID(c)
	listingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	// Verify listing ownership
	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ?", listingID, userID).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	// Only draft and rejected can upload photos
	if listing.Status != models.ListingStatusDraft && listing.Status != models.ListingStatusRejected {
		utils.Unprocessable(c, "BIZ_LISTING_NOT_EDITABLE", "Foto hanya dapat diunggah saat listing berstatus draft atau rejected.")
		return
	}

	// Count existing photos
	var existingCount int64
	database.DB.Model(&models.PropertyPhoto{}).Where("listing_id = ?", listingID).Count(&existingCount)

	maxPhotos := AppConfig.MaxPhotosPerListing
	if maxPhotos <= 0 {
		maxPhotos = 10
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "Format request tidak valid. Gunakan multipart/form-data.")
		return
	}

	files := form.File["photos"]
	if len(files) == 0 {
		utils.BadRequest(c, "Tidak ada file foto yang dikirim.")
		return
	}

	if int(existingCount)+len(files) > maxPhotos {
		utils.Unprocessable(c, "BIZ_MAX_PHOTOS_EXCEEDED",
			fmt.Sprintf("Maksimal %d foto per listing. Saat ini sudah ada %d foto.", maxPhotos, existingCount))
		return
	}

	uploadDir := AppConfig.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Get current max sort_order
	var maxSort int
	database.DB.Model(&models.PropertyPhoto{}).
		Where("listing_id = ?", listingID).
		Select("COALESCE(MAX(sort_order), -1)").
		Scan(&maxSort)

	photos := make([]gin.H, 0, len(files))
	subDir := "listings/" + listingID.String()

	for i, fh := range files {
		// Validate file
		if err := utils.ValidateImageFile(fh, AppConfig.MaxUploadSizeMB); err != nil {
			utils.Unprocessable(c, "VAL_FILE_INVALID",
				fmt.Sprintf("Foto #%d (%s): %s", i+1, fh.Filename, err.Error()))
			return
		}

		// Save original
		origPath, origURL, err := utils.SaveUploadedFile(fh, uploadDir, subDir+"/original")
		if err != nil {
			utils.InternalError(c, fmt.Sprintf("Gagal menyimpan foto #%d. Silakan coba lagi.", i+1))
			return
		}

		// Strip EXIF metadata for privacy (log warning on failure, non-blocking)
		if err := utils.StripEXIF(origPath); err != nil {
			log.Printf("[WARN] Gagal strip EXIF untuk %s: %v", origPath, err)
		}

		// Generate thumbnail (150x150)
		thumbFilename := "thumb_" + filepath.Base(origPath)
		thumbPath := filepath.Join(uploadDir, subDir, "thumbnail", thumbFilename)
		thumbURL := ""
		if err := utils.GenerateThumbnail(origPath, thumbPath, 150, 150); err == nil {
			thumbURL = "/uploads/" + subDir + "/thumbnail/" + thumbFilename
		}

		// Generate medium (800x600)
		mediumFilename := "med_" + filepath.Base(origPath)
		mediumPath := filepath.Join(uploadDir, subDir, "medium", mediumFilename)
		mediumURL := ""
		if err := utils.GenerateThumbnail(origPath, mediumPath, 800, 600); err == nil {
			mediumURL = "/uploads/" + subDir + "/medium/" + mediumFilename
		}

		// Generate watermarked copy
		watermarkFilename := "wm_" + filepath.Base(origPath)
		watermarkPath := filepath.Join(uploadDir, subDir, "watermarked", watermarkFilename)
		watermarkedURL := origURL // fallback to /uploads/... URL from SaveUploadedFile
		if err := utils.GenerateWatermarked(origPath, watermarkPath, AppConfig.PlatformName); err == nil {
			watermarkedURL = "/uploads/" + subDir + "/watermarked/" + watermarkFilename
		} else {
			log.Printf("[WARN] Gagal generate watermark untuk %s: %v", origPath, err)
		}

		sortOrder := maxSort + i + 1
		photo := models.PropertyPhoto{
			ListingID:      listingID,
			FileName:       &fh.Filename,
			OriginalURL:    origURL,
			ThumbnailURL:   &thumbURL,
			MediumURL:      &mediumURL,
			WatermarkedURL: watermarkedURL,
			SortOrder:      sortOrder,
		}

		if err := database.DB.Create(&photo).Error; err != nil {
			utils.InternalError(c, fmt.Sprintf("Gagal menyimpan data foto #%d.", i+1))
			return
		}

		photos = append(photos, gin.H{
			"id":              photo.ID,
			"file_name":       photo.FileName,
			"original_url":    photo.OriginalURL,
			"thumbnail_url":   photo.ThumbnailURL,
			"medium_url":      photo.MediumURL,
			"watermarked_url": photo.WatermarkedURL,
			"sort_order":      photo.SortOrder,
		})
	}

	utils.Created(c, gin.H{
		"message":      fmt.Sprintf("%d foto berhasil diunggah.", len(photos)),
		"photos_count": len(photos),
		"photos":       photos,
	})
}

func SalesmanDeletePhoto(c *gin.Context) {
	userID := middleware.GetUserID(c)
	listingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}
	photoID, err := uuid.Parse(c.Param("photoId"))
	if err != nil {
		utils.NotFound(c, "RES_PHOTO_NOT_FOUND", "Foto tidak ditemukan.")
		return
	}

	// Verify listing ownership
	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ?", listingID, userID).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	// Only draft and rejected can delete photos
	if listing.Status != models.ListingStatusDraft && listing.Status != models.ListingStatusRejected {
		utils.Unprocessable(c, "BIZ_LISTING_NOT_EDITABLE", "Foto hanya dapat dihapus saat listing berstatus draft atau rejected.")
		return
	}

	// Find the photo
	var photo models.PropertyPhoto
	if err := database.DB.Where("id = ? AND listing_id = ?", photoID, listingID).First(&photo).Error; err != nil {
		utils.NotFound(c, "RES_PHOTO_NOT_FOUND", "Foto tidak ditemukan.")
		return
	}

	uploadDir := AppConfig.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Delete files from disk
	_ = utils.DeleteFile(uploadDir, photo.OriginalURL)
	if photo.ThumbnailURL != nil {
		_ = utils.DeleteFile(uploadDir, *photo.ThumbnailURL)
	}
	if photo.MediumURL != nil {
		_ = utils.DeleteFile(uploadDir, *photo.MediumURL)
	}
	_ = utils.DeleteFile(uploadDir, photo.WatermarkedURL)

	// Delete from DB
	database.DB.Delete(&photo)

	utils.OK(c, gin.H{"message": "Foto berhasil dihapus."})
}

func SalesmanReorderPhotos(c *gin.Context) {
	var req dto.ReorderPhotosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	for i, photoID := range req.PhotoIDs {
		database.DB.Model(&models.PropertyPhoto{}).Where("id = ?", photoID).Update("sort_order", i)
	}

	utils.OK(c, gin.H{"message": "Urutan foto berhasil diubah."})
}

func SalesmanGetQuota(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var subscription models.Subscription
	database.DB.Where("tenant_id = ?", tenantID).First(&subscription)

	var activeCount int64
	database.DB.Model(&models.PropertyListing{}).
		Where("salesman_id = ? AND status IN ?", userID, []string{
			models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved,
		}).Count(&activeCount)

	utils.OK(c, gin.H{
		"used":      activeCount,
		"max":       subscription.MaxListingsPerSalesman,
		"remaining": subscription.MaxListingsPerSalesman - int(activeCount),
		"plan_type": subscription.PlanType,
		"status_breakdown": gin.H{
			"draft":    0,
			"pending":  0,
			"approved": 0,
		},
	})
}

// ==================== HELPERS ====================

func checkQuota(salesmanID, tenantID uuid.UUID) bool {
	var subscription models.Subscription
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&subscription).Error; err != nil {
		return false
	}

	var count int64
	database.DB.Model(&models.PropertyListing{}).
		Where("salesman_id = ? AND status IN ?", salesmanID, []string{
			models.ListingStatusDraft, models.ListingStatusPending, models.ListingStatusApproved,
		}).Count(&count)

	return int(count) < subscription.MaxListingsPerSalesman
}

func markListingStatus(c *gin.Context, status, label string) {
	userID := middleware.GetUserID(c)
	listingID, _ := uuid.Parse(c.Param("id"))

	var listing models.PropertyListing
	if err := database.DB.Where("id = ? AND salesman_id = ? AND status = ?", listingID, userID, models.ListingStatusApproved).First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	database.DB.Model(&listing).Update("status", status)

	utils.OK(c, gin.H{
		"id":      listing.ID,
		"status":  status,
		"message": "Listing berhasil ditandai sebagai " + label + ".",
	})
}

func buildListingUpdates(req dto.UpdateListingRequest) map[string]interface{} {
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.ListingType != nil {
		updates["listing_type"] = *req.ListingType
	}
	if req.PropertyType != nil {
		updates["property_type"] = *req.PropertyType
	}
	if req.SourceType != nil {
		updates["source_type"] = *req.SourceType
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.Province != nil {
		updates["province"] = *req.Province
	}
	if req.LandArea != nil {
		updates["land_area"] = *req.LandArea
	}
	if req.BuildingArea != nil {
		updates["building_area"] = *req.BuildingArea
	}
	if req.Bedrooms != nil {
		updates["bedrooms"] = *req.Bedrooms
	}
	if req.Bathrooms != nil {
		updates["bathrooms"] = *req.Bathrooms
	}
	if req.Floors != nil {
		updates["floors"] = *req.Floors
	}
	if req.CertificateType != nil {
		updates["certificate_type"] = *req.CertificateType
	}
	if req.Facilities != nil {
		updates["facilities"] = models.Facilities(req.Facilities)
	}
	return updates
}
