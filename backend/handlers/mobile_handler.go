package handlers

import (
	"math"
	"net/http"

	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ================================================================
// MOBILE API HANDLERS
// Optimized for mobile: flatter responses, smaller payloads,
// cursor-based pagination support, device token management.
// Prefix: /api/mobile/v1
// ================================================================

// ==================== MOBILE AUTH ====================

type MobileLoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	DeviceToken string `json:"device_token"` // optional — for push notifications
	DeviceOS    string `json:"device_os"`    // optional — "ios" / "android"
}

type MobileRegisterRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Phone       string `json:"phone" binding:"required,min=8,max=20"`
	Password    string `json:"password" binding:"required,min=8"`
	DeviceToken string `json:"device_token"`
	DeviceOS    string `json:"device_os"`
}

type MobileAuthResponse struct {
	Token    string             `json:"token"`
	User     MobileUserBrief    `json:"user"`
	Tenant   *MobileTenantBrief `json:"tenant,omitempty"`
}

type MobileUserBrief struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	PhotoURL string `json:"photo_url,omitempty"`
	Role     string `json:"role"`
}

type MobileTenantBrief struct {
	ID               string `json:"id"`
	OrganizationName string `json:"organization_name"`
	SubdomainSlug    string `json:"subdomain_slug"`
	LogoURL          string `json:"logo_url,omitempty"`
}

func MobileLogin(c *gin.Context) {
	var req MobileLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Email dan password diperlukan.", nil)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Unauthorized(c, "AUTH_INVALID", "Email atau password salah.")
		return
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		utils.Unauthorized(c, "AUTH_INVALID", "Email atau password salah.")
		return
	}

	if user.Status != models.UserStatusActive {
		utils.Forbidden(c, "AUTH_INACTIVE", "Akun Anda tidak aktif. Hubungi admin.")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role, user.TenantID, AppConfig)
	if err != nil {
		utils.InternalError(c, "Gagal membuat token.")
		return
	}

	// Build response
	resp := MobileAuthResponse{
		Token: token,
		User: MobileUserBrief{
			ID:       user.ID.String(),
			Name:     user.Name,
			Email:    user.Email,
			Phone:    ptrToStr(user.Phone),
			PhotoURL: ptrToStr(user.PhotoURL),
			Role:     user.Role,
		},
	}

	// Load tenant info for tenant-bound users
	if user.TenantID != nil {
		var tenant models.Tenant
		if err := database.DB.First(&tenant, user.TenantID).Error; err == nil {
			resp.Tenant = &MobileTenantBrief{
				ID:               tenant.ID.String(),
				OrganizationName: tenant.OrganizationName,
				SubdomainSlug:    tenant.SubdomainSlug,
				LogoURL:          ptrToStr(tenant.LogoURL),
			}
		}
	}

	utils.Success(c, http.StatusOK, resp)
}

func MobileRegister(c *gin.Context) {
	var req MobileRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid.", nil)
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Conflict(c, "AUTH_EMAIL_REGISTERED", "Email sudah terdaftar.")
		return
	}

	hash, err := utils.HashPassword(req.Password, AppConfig.BcryptCost)
	if err != nil {
		utils.InternalError(c, "Terjadi kesalahan server.")
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Phone:        &req.Phone,
		Role:         models.RoleBuyer,
		Status:       models.UserStatusActive,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, "Gagal mendaftarkan akun.")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role, user.TenantID, AppConfig)
	if err != nil {
		utils.InternalError(c, "Gagal membuat token.")
		return
	}

	resp := MobileAuthResponse{
		Token: token,
		User: MobileUserBrief{
			ID:       user.ID.String(),
			Name:     user.Name,
			Email:    user.Email,
			Phone:    ptrToStr(user.Phone),
			PhotoURL: ptrToStr(user.PhotoURL),
			Role:     user.Role,
		},
	}

	utils.Created(c, resp)
}

// ==================== MOBILE PROPERTIES ====================

type MobilePropertyCard struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	ListingType  string  `json:"listing_type"`
	PropertyType string  `json:"property_type"`
	SourceType   string  `json:"source_type"`
	City         string  `json:"city,omitempty"`
	ThumbnailURL string  `json:"thumbnail_url,omitempty"`
	Bedrooms     int     `json:"bedrooms,omitempty"`
	Bathrooms    int     `json:"bathrooms,omitempty"`
	LandArea     float64 `json:"land_area,omitempty"`
	BuildingArea float64 `json:"building_area,omitempty"`
	AgencyName   string  `json:"agency_name,omitempty"`
	IsFeatured   bool    `json:"is_featured"`
	CreatedAt    string  `json:"created_at"`
}

type MobilePropertyListResponse struct {
	Properties []MobilePropertyCard `json:"properties"`
	NextCursor string               `json:"next_cursor,omitempty"` // for infinite scroll
	HasMore    bool                 `json:"has_more"`
}

func MobileListProperties(c *gin.Context) {
	p := utils.GetPagination(c)
	perPage := p.PerPage
	if perPage > 50 {
		perPage = 50 // mobile limit
	}

	var listings []models.PropertyListing
	query := database.DB.
		Preload("Photos", database.DB.Where("is_primary = ?", true).Or("sort_order = ?", 0)).
		Preload("Tenant").
		Where("status = ?", models.ListingStatusApproved)

	// Filters — same as web
	if pt := c.Query("property_type"); pt != "" {
		query = query.Where("property_type = ?", pt)
	}
	if lt := c.Query("listing_type"); lt != "" {
		query = query.Where("listing_type = ?", lt)
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}
	if pmin := c.Query("price_min"); pmin != "" {
		query = query.Where("price >= ?", pmin)
	}
	if pmax := c.Query("price_max"); pmax != "" {
		query = query.Where("price <= ?", pmax)
	}
	if q := c.Query("q"); q != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	// Sort
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	if sort == "price" {
		query = query.Order("price " + order)
	} else {
		query = query.Order("created_at " + order)
	}

	// Cursor-based pagination (mobile infinite scroll)
	cursor := c.Query("cursor")
	if cursor != "" {
		query = query.Where("created_at < ?", cursor)
	}

	var total int64
	query.Count(&total)

	query = query.Limit(perPage + 1) // fetch one extra to detect has_more
	if err := query.Find(&listings).Error; err != nil {
		utils.InternalError(c, "Gagal memuat properti.")
		return
	}

	hasMore := len(listings) > perPage
	if hasMore {
		listings = listings[:perPage]
	}

	cards := make([]MobilePropertyCard, len(listings))
	var nextCursor string
	for i, l := range listings {
		thumb := ""
		if len(l.Photos) > 0 {
			thumb = ptrToStr(l.Photos[0].ThumbnailURL)
			if thumb == "" {
				thumb = l.Photos[0].OriginalURL
			}
		}
		agency := ""
		if l.Tenant != nil {
			agency = l.Tenant.OrganizationName
		}
		cards[i] = MobilePropertyCard{
			ID:           l.ID.String(),
			Title:        l.Title,
			Price:        l.Price,
			ListingType:  l.ListingType,
			PropertyType: l.PropertyType,
			SourceType:   l.SourceType,
			City:         ptrToStr(l.City),
			ThumbnailURL: thumb,
			Bedrooms:     ptrToInt(l.Bedrooms),
			Bathrooms:    ptrToInt(l.Bathrooms),
			LandArea:     ptrToFloat(l.LandArea),
			BuildingArea: ptrToFloat(l.BuildingArea),
			AgencyName:   agency,
			CreatedAt:    l.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if hasMore && i == len(listings)-1 {
			nextCursor = l.CreatedAt.Format("2006-01-02T15:04:05Z")
		}
	}

	utils.Success(c, http.StatusOK, MobilePropertyListResponse{
		Properties: cards,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
}

type MobilePropertyDetail struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description,omitempty"`
	Price          float64                `json:"price"`
	ListingType    string                 `json:"listing_type"`
	PropertyType   string                 `json:"property_type"`
	SourceType     string                 `json:"source_type"`
	Address        string                 `json:"address,omitempty"`
	City           string                 `json:"city,omitempty"`
	Province       string                 `json:"province,omitempty"`
	Latitude       float64                `json:"latitude,omitempty"`
	Longitude      float64                `json:"longitude,omitempty"`
	LandArea       float64                `json:"land_area,omitempty"`
	BuildingArea   float64                `json:"building_area,omitempty"`
	Bedrooms       int                    `json:"bedrooms,omitempty"`
	Bathrooms      int                    `json:"bathrooms,omitempty"`
	Floors         int                    `json:"floors,omitempty"`
	Certificate    string                 `json:"certificate_type,omitempty"`
	Status         string                 `json:"status"`
	Photos         []MobilePhotoBrief     `json:"photos"`
	Facilities     map[string]interface{} `json:"facilities,omitempty"`
	Agency         *MobileTenantBrief     `json:"agency,omitempty"`
	Salesman       *MobileUserBrief       `json:"salesman,omitempty"`
	AuctionDetail  *MobileAuctionBrief    `json:"auction_detail,omitempty"`
	CompanyDetail  *MobileCompanyBrief    `json:"company_detail,omitempty"`
	CreatedAt      string                 `json:"created_at"`
}

type MobilePhotoBrief struct {
	ID           string `json:"id"`
	OriginalURL  string `json:"original_url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	MediumURL    string `json:"medium_url,omitempty"`
	IsPrimary    bool   `json:"is_primary"`
}

type MobileAuctionBrief struct {
	BankName     string  `json:"bank_name"`
	AuctionDate  string  `json:"auction_date,omitempty"`
	LimitPrice   float64 `json:"limit_price,omitempty"`
	Deposit      float64 `json:"deposit,omitempty"`
	Location     string  `json:"location,omitempty"`
}

type MobileCompanyBrief struct {
	CompanyName string `json:"company_name"`
	Disposal    string `json:"disposal_type"`
	PICName     string `json:"pic_name,omitempty"`
	PICPhone    string `json:"pic_phone,omitempty"`
}

func MobileGetPropertyDetail(c *gin.Context) {
	id := c.Param("id")

	var listing models.PropertyListing
	if err := database.DB.
		Preload("Photos", database.DB.Order("sort_order ASC")).
		Preload("Tenant").
		Preload("Salesman").
		Preload("BankAuctionDetail").
		Preload("CompanyAssetDetail").
		Where("id = ? AND status = ?", id, models.ListingStatusApproved).
		First(&listing).Error; err != nil {
		utils.NotFound(c, "PROPERTY_NOT_FOUND", "Properti tidak ditemukan.")
		return
	}

	photos := make([]MobilePhotoBrief, len(listing.Photos))
	for i, p := range listing.Photos {
		photos[i] = MobilePhotoBrief{
			ID:           p.ID.String(),
			OriginalURL:  p.OriginalURL,
			ThumbnailURL: ptrToStr(p.ThumbnailURL),
			MediumURL:    ptrToStr(p.MediumURL),
			IsPrimary:    p.IsPrimary,
		}
	}

	detail := MobilePropertyDetail{
		ID:           listing.ID.String(),
		Title:        listing.Title,
		Description:  ptrToStr(listing.Description),
		Price:        listing.Price,
		ListingType:  listing.ListingType,
		PropertyType: listing.PropertyType,
		SourceType:   listing.SourceType,
		Address:      ptrToStr(listing.Address),
		City:         ptrToStr(listing.City),
		Province:     ptrToStr(listing.Province),
		Latitude:     ptrToFloat(listing.Latitude),
		Longitude:    ptrToFloat(listing.Longitude),
		LandArea:     ptrToFloat(listing.LandArea),
		BuildingArea: ptrToFloat(listing.BuildingArea),
		Bedrooms:     ptrToInt(listing.Bedrooms),
		Bathrooms:    ptrToInt(listing.Bathrooms),
		Floors:       ptrToInt(listing.Floors),
		Certificate:  ptrToStr(listing.CertificateType),
		Status:       listing.Status,
		Photos:       photos,
		Facilities:   listing.Facilities,
		CreatedAt:    listing.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if listing.Tenant != nil {
		detail.Agency = &MobileTenantBrief{
			ID:               listing.Tenant.ID.String(),
			OrganizationName: listing.Tenant.OrganizationName,
			SubdomainSlug:    listing.Tenant.SubdomainSlug,
			LogoURL:          ptrToStr(listing.Tenant.LogoURL),
		}
	}

	if listing.Salesman != nil {
		detail.Salesman = &MobileUserBrief{
			ID:       listing.Salesman.ID.String(),
			Name:     listing.Salesman.Name,
			Email:    listing.Salesman.Email,
			Phone:    ptrToStr(listing.Salesman.Phone),
			PhotoURL: ptrToStr(listing.Salesman.PhotoURL),
			Role:     listing.Salesman.Role,
		}
	}

	if listing.BankAuctionDetail != nil {
		detail.AuctionDetail = &MobileAuctionBrief{
			BankName:   listing.BankAuctionDetail.BankName,
			LimitPrice: ptrToFloat(listing.BankAuctionDetail.AuctionLimitPrice),
			Deposit:    ptrToFloat(listing.BankAuctionDetail.AuctionDeposit),
		}
		if listing.BankAuctionDetail.AuctionDate != nil {
			detail.AuctionDetail.AuctionDate = listing.BankAuctionDetail.AuctionDate.Format("2006-01-02T15:04:05Z")
		}
		if listing.BankAuctionDetail.AuctionLocation != nil {
			detail.AuctionDetail.Location = *listing.BankAuctionDetail.AuctionLocation
		}
	}

	if listing.CompanyAssetDetail != nil {
		detail.CompanyDetail = &MobileCompanyBrief{
			CompanyName: listing.CompanyAssetDetail.CompanyName,
			Disposal:    listing.CompanyAssetDetail.DisposalType,
			PICName:     ptrToStr(listing.CompanyAssetDetail.PICName),
			PICPhone:    ptrToStr(listing.CompanyAssetDetail.PICPhone),
		}
	}

	utils.Success(c, http.StatusOK, detail)
}

// ==================== MOBILE BUYER ====================

type MobileInquiryRequest struct {
	PropertyID string `json:"property_id" binding:"required,uuid"`
	Message    string `json:"message" binding:"required,min=1,max=2000"`
}

func MobileCreateInquiry(c *gin.Context) {
	userID := c.GetString("userID")

	var req MobileInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Property ID dan pesan diperlukan.", nil)
		return
	}

	inquiry := models.Inquiry{
		BuyerID:    mustUUID(userID),
		PropertyID: mustUUID(req.PropertyID),
		Message:    &req.Message,
		Status:     models.InquiryStatusUnread,
	}

	if err := database.DB.Create(&inquiry).Error; err != nil {
		utils.InternalError(c, "Gagal mengirim pertanyaan.")
		return
	}

	utils.Created(c, gin.H{"id": inquiry.ID.String(), "status": inquiry.Status})
}

func MobileListBuyerInquiries(c *gin.Context) {
	userID := c.GetString("userID")
	p := utils.GetPagination(c)

	var inquiries []models.Inquiry
	var total int64

	query := database.DB.
		Preload("Property", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title, status")
		}).
		Where("buyer_id = ?", userID)

	query.Count(&total)
	if err := query.Offset(p.Offset).Limit(p.PerPage).Order("created_at DESC").Find(&inquiries).Error; err != nil {
		utils.InternalError(c, "Gagal memuat pertanyaan.")
		return
	}

	type inquiryItem struct {
		ID         string `json:"id"`
		PropertyID string `json:"property_id"`
		Title      string `json:"property_title"`
		Message    string `json:"message"`
		Status     string `json:"status"`
		CreatedAt  string `json:"created_at"`
	}

	items := make([]inquiryItem, len(inquiries))
	for i, inq := range inquiries {
		title := ""
		if inq.Property != nil {
			title = inq.Property.Title
		}
		items[i] = inquiryItem{
			ID:         inq.ID.String(),
			PropertyID: inq.PropertyID.String(),
			Title:      title,
			Message:    ptrToStr(inq.Message),
			Status:     inq.Status,
			CreatedAt:  inq.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	utils.SuccessPaginated(c, http.StatusOK, items, utils.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PerPage))),
	})
}

// ==================== MOBILE PROFILE ====================

func MobileGetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.NotFound(c, "PROFILE_NOT_FOUND", "Profil tidak ditemukan.")
		return
	}

	type mobileProfile struct {
		ID         string             `json:"id"`
		Name       string             `json:"name"`
		Email      string             `json:"email"`
		Phone      string             `json:"phone,omitempty"`
		PhotoURL   string             `json:"photo_url,omitempty"`
		Role       string             `json:"role"`
		Tenant     *MobileTenantBrief `json:"tenant,omitempty"`
		Whatsapp   string             `json:"whatsapp,omitempty"`
		ShowWA     bool               `json:"show_whatsapp"`
	}

	profile := mobileProfile{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Phone:    ptrToStr(user.Phone),
		PhotoURL: ptrToStr(user.PhotoURL),
		Role:     user.Role,
		Whatsapp: ptrToStr(user.WhatsappNumber),
		ShowWA:   user.ShowWhatsapp != nil && *user.ShowWhatsapp,
	}

	if user.TenantID != nil {
		var tenant models.Tenant
		if err := database.DB.First(&tenant, user.TenantID).Error; err == nil {
			profile.Tenant = &MobileTenantBrief{
				ID:               tenant.ID.String(),
				OrganizationName: tenant.OrganizationName,
				SubdomainSlug:    tenant.SubdomainSlug,
				LogoURL:          ptrToStr(tenant.LogoURL),
			}
		}
	}

	_ = role // reserved for role-specific profile extensions
	utils.Success(c, http.StatusOK, profile)
}

func MobileUpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Data tidak valid.", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if len(updates) == 0 {
		utils.BadRequest(c, "Tidak ada data yang diubah.")
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Gagal menyimpan profil.")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "Profil berhasil diperbarui."})
}

// ==================== MOBILE SAVED/FAVORITES ====================

func MobileToggleSave(c *gin.Context) {
	userID := c.GetString("userID")
	propertyID := c.Param("propertyId")

	var existing models.SavedProperty
	if err := database.DB.Where("buyer_id = ? AND listing_id = ?", userID, propertyID).First(&existing).Error; err == nil {
		// Already saved — remove
		database.DB.Delete(&existing)
		utils.Success(c, http.StatusOK, gin.H{"saved": false, "message": "Dihapus dari favorit."})
		return
	}

	// Save
	sp := models.SavedProperty{
		BuyerID:   mustUUID(userID),
		ListingID: mustUUID(propertyID),
	}
	if err := database.DB.Create(&sp).Error; err != nil {
		utils.InternalError(c, "Gagal menyimpan.")
		return
	}

	utils.Created(c, gin.H{"saved": true, "message": "Disimpan ke favorit."})
}

func MobileListSaved(c *gin.Context) {
	userID := c.GetString("userID")
	p := utils.GetPagination(c)

	var saved []models.SavedProperty
	var total int64

	query := database.DB.
		Preload("Listing", database.DB.Where("status = ?", models.ListingStatusApproved)).
		Preload("Listing.Photos", database.DB.Where("is_primary = ?", true).Or("sort_order = ?", 0)).
		Preload("Listing.Tenant").
		Where("buyer_id = ?", userID)

	query.Count(&total)
	if err := query.Offset(p.Offset).Limit(p.PerPage).Order("created_at DESC").Find(&saved).Error; err != nil {
		utils.InternalError(c, "Gagal memuat favorit.")
		return
	}

	cards := make([]MobilePropertyCard, 0, len(saved))
	for _, s := range saved {
		if s.Listing == nil {
			continue
		}
		l := s.Listing
		thumb := ""
		if len(l.Photos) > 0 {
			thumb = ptrToStr(l.Photos[0].ThumbnailURL)
			if thumb == "" {
				thumb = l.Photos[0].OriginalURL
			}
		}
		agency := ""
		if l.Tenant != nil {
			agency = l.Tenant.OrganizationName
		}
		cards = append(cards, MobilePropertyCard{
			ID:           l.ID.String(),
			Title:        l.Title,
			Price:        l.Price,
			ListingType:  l.ListingType,
			PropertyType: l.PropertyType,
			SourceType:   l.SourceType,
			City:         ptrToStr(l.City),
			ThumbnailURL: thumb,
			Bedrooms:     ptrToInt(l.Bedrooms),
			Bathrooms:    ptrToInt(l.Bathrooms),
			AgencyName:   agency,
			CreatedAt:    l.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	utils.SuccessPaginated(c, http.StatusOK, cards, utils.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PerPage))),
	})
}

// ==================== MOBILE SALESMAN ====================

type MobileSalesmanDashboardData struct {
	TotalListings   int                    `json:"total_listings"`
	ActiveListings  int                    `json:"active_listings"`
	PendingListings int                    `json:"pending_listings"`
	TotalInquiries  int                    `json:"total_inquiries"`
	UnreadInquiries int                    `json:"unread_inquiries"`
	Quota           MobileQuota            `json:"quota"`
	RecentListings  []MobilePropertyCard   `json:"recent_listings"`
}

type MobileQuota struct {
	Used  int `json:"used"`
	Max   int `json:"max"`
}

func MobileSalesmanDashboard(c *gin.Context) {
	userID := c.GetString("userID")
	tenantID := c.GetString("tenantID")

	var total, active, pending, inquiries, unread int64

	database.DB.Model(&models.PropertyListing{}).Where("salesman_id = ?", userID).Count(&total)
	database.DB.Model(&models.PropertyListing{}).Where("salesman_id = ? AND status = ?", userID, models.ListingStatusApproved).Count(&active)
	database.DB.Model(&models.PropertyListing{}).Where("salesman_id = ? AND status = ?", userID, models.ListingStatusPending).Count(&pending)

	// Count inquiries on salesman's listings
	database.DB.Model(&models.Inquiry{}).
		Joins("JOIN property_listings ON property_listings.id = inquiries.property_id").
		Where("property_listings.salesman_id = ?", userID).
		Count(&inquiries)
	database.DB.Model(&models.Inquiry{}).
		Joins("JOIN property_listings ON property_listings.id = inquiries.property_id").
		Where("property_listings.salesman_id = ? AND inquiries.status = ?", userID, models.InquiryStatusUnread).
		Count(&unread)

	// Quota
	quota := MobileQuota{Max: 5} // default free
	if tenantID != "" {
		var sub models.Subscription
		if err := database.DB.Where("tenant_id = ?", tenantID).First(&sub).Error; err == nil {
			quota.Max = sub.MaxListingsPerSalesman
		}
	}

	// Recent listings
	var recent []models.PropertyListing
	database.DB.
		Preload("Photos", database.DB.Where("is_primary = ?", true).Or("sort_order = ?", 0)).
		Preload("Tenant").
		Where("salesman_id = ?", userID).
		Order("created_at DESC").
		Limit(5).
		Find(&recent)

	recentCards := make([]MobilePropertyCard, len(recent))
	for i, l := range recent {
		thumb := ""
		if len(l.Photos) > 0 {
			thumb = ptrToStr(l.Photos[0].ThumbnailURL)
			if thumb == "" {
				thumb = l.Photos[0].OriginalURL
			}
		}
		recentCards[i] = MobilePropertyCard{
			ID:           l.ID.String(),
			Title:        l.Title,
			Price:        l.Price,
			ListingType:  l.ListingType,
			PropertyType: l.PropertyType,
			City:         ptrToStr(l.City),
			ThumbnailURL: thumb,
			CreatedAt:    l.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	utils.Success(c, http.StatusOK, MobileSalesmanDashboardData{
		TotalListings:   int(total),
		ActiveListings:  int(active),
		PendingListings: int(pending),
		TotalInquiries:  int(inquiries),
		UnreadInquiries: int(unread),
		Quota:           quota,
		RecentListings:  recentCards,
	})
}

// ==================== MOBILE TENANT ====================

type MobileTenantDashboardData struct {
	TotalListings   int    `json:"total_listings"`
	ActiveListings  int    `json:"active_listings"`
	TotalSalesmen   int    `json:"total_salesmen"`
	MaxSalesmen     int    `json:"max_salesmen"`
	PlanType        string `json:"plan_type"`
	PlanName        string `json:"plan_name,omitempty"`
}

func MobileTenantDashboard(c *gin.Context) {
	tenantID := c.GetString("tenantID")

	var total, active, salesmenCount int64

	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ?", tenantID).Count(&total)
	database.DB.Model(&models.PropertyListing{}).Where("tenant_id = ? AND status = ?", tenantID, models.ListingStatusApproved).Count(&active)

	// Count active salesmen
	database.DB.Model(&models.TenantUser{}).
		Where("tenant_id = ? AND tenant_role = ?", tenantID, models.TenantRoleSalesman).
		Count(&salesmenCount)

	dash := MobileTenantDashboardData{
		TotalListings:  int(total),
		ActiveListings: int(active),
		TotalSalesmen:  int(salesmenCount),
		MaxSalesmen:    5,
		PlanType:       "free",
	}

	// Load subscription
	var sub models.Subscription
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&sub).Error; err == nil {
		dash.MaxSalesmen = sub.MaxSalesmen
		dash.PlanType = sub.PlanType
		dash.PlanName = sub.PlanType // "free" / "premium"
	}

	utils.Success(c, http.StatusOK, dash)
}

func MobileTenantSubscription(c *gin.Context) {
	tenantID := c.GetString("tenantID")

	type subResponse struct {
		PlanType               string `json:"plan_type"`
		PlanName               string `json:"plan_name"`
		MaxSalesmen            int    `json:"max_salesmen"`
		MaxListingsPerSalesman int    `json:"max_listings_per_salesman"`
		SalesmenUsed           int    `json:"salesmen_used"`
		TotalActiveListings    int    `json:"total_active_listings"`
	}

	resp := subResponse{
		MaxSalesmen:            5,
		MaxListingsPerSalesman: 5,
		PlanType:               "free",
		PlanName:               "Free",
	}

	var sub models.Subscription
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&sub).Error; err == nil {
		resp.MaxSalesmen = sub.MaxSalesmen
		resp.MaxListingsPerSalesman = sub.MaxListingsPerSalesman
		resp.PlanType = sub.PlanType
		resp.PlanName = sub.PlanType
	}

	// Count usage
	var sc int64
	database.DB.Model(&models.TenantUser{}).
		Where("tenant_id = ? AND tenant_role = ?", tenantID, models.TenantRoleSalesman).
		Count(&sc)
	resp.SalesmenUsed = int(sc)

	var lc int64
	database.DB.Model(&models.PropertyListing{}).
		Where("tenant_id = ? AND status = ?", tenantID, models.ListingStatusApproved).
		Count(&lc)
	resp.TotalActiveListings = int(lc)

	utils.Success(c, http.StatusOK, resp)
}

func MobileTenantRequestUpgrade(c *gin.Context) {
	tenantID := c.GetString("tenantID")

	var sub models.Subscription
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&sub).Error; err != nil {
		utils.NotFound(c, "SUBSCRIPTION_NOT_FOUND", "Langganan tidak ditemukan.")
		return
	}

	// Check if already requested (treat pending_upgrade as already requested)
	if sub.PlanType == "pending_upgrade" {
		utils.Conflict(c, "UPGRADE_PENDING", "Permintaan upgrade sudah dikirim sebelumnya.")
		return
	}

	if err := database.DB.Model(&sub).Update("plan_type", "pending_upgrade").Error; err != nil {
		utils.InternalError(c, "Gagal mengirim permintaan upgrade.")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "Permintaan upgrade telah dikirim.", "status": "pending_upgrade"})
}

// ==================== MOBILE MASTER DATA ====================

func MobilePropertyTypes(c *gin.Context) {
	var types []models.PropertyType
	database.DB.Where("is_active = ?", true).Order("name ASC").Find(&types)

	type item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	items := make([]item, len(types))
	for i, t := range types {
		items[i] = item{ID: t.ID.String(), Name: t.Name, Slug: t.Slug}
	}
	utils.Success(c, http.StatusOK, items)
}

func MobileFacilities(c *gin.Context) {
	var facilities []models.Facility
	database.DB.Where("is_active = ?", true).Order("name ASC").Find(&facilities)

	type item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Icon string `json:"icon,omitempty"`
	}
	items := make([]item, len(facilities))
	for i, f := range facilities {
		items[i] = item{ID: f.ID.String(), Name: f.Name, Icon: ptrToStr(f.Icon)}
	}
	utils.Success(c, http.StatusOK, items)
}

func MobileLocations(c *gin.Context) {
	var locations []models.Location
	query := database.DB.Where("is_active = ?", true).Order("province, city ASC")

	if province := c.Query("province"); province != "" {
		query = query.Where("province ILIKE ?", "%"+province+"%")
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}

	query.Find(&locations)

	type item struct {
		ID        string  `json:"id"`
		City      string  `json:"city"`
		Province  string  `json:"province"`
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	}
	items := make([]item, len(locations))
	for i, l := range locations {
		items[i] = item{
			ID:        l.ID.String(),
			City:      l.City,
			Province:  l.Province,
			Latitude:  ptrToFloat(l.Latitude),
			Longitude: ptrToFloat(l.Longitude),
		}
	}
	utils.Success(c, http.StatusOK, items)
}

// ==================== HELPERS ====================

func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func ptrToFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func mustUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
