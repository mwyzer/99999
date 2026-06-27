package handlers

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== PUBLIC PROPERTY ENDPOINTS ====================

func ListProperties(c *gin.Context) {
	p := utils.GetPagination(c)

	var listings []models.PropertyListing
	query := database.DB.
		Preload("Photos", "sort_order = ?", 0).
		Preload("Salesman").
		Preload("Tenant").
		Where("status = ?", models.ListingStatusApproved)

	// Filters
	if pt := c.Query("property_type"); pt != "" {
		query = query.Where("property_type = ?", pt)
	}
	if st := c.Query("source_type"); st != "" {
		query = query.Where("source_type = ?", st)
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

	var total int64
	query.Model(&models.PropertyListing{}).Count(&total)

	if err := query.Offset(p.Offset).Limit(p.PerPage).Find(&listings).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	cards := make([]dto.PropertyCardResponse, 0, len(listings))
	for _, l := range listings {
		card := toPropertyCard(l)
		cards = append(cards, card)
	}

	utils.SuccessPaginated(c, 200, cards, utils.CalculateMeta(p.Page, p.PerPage, total))
}

func GetPropertyDetail(c *gin.Context) {
	id := c.Param("id")

	var listing models.PropertyListing
	if err := database.DB.
		Preload("Photos", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Salesman").
		Preload("Tenant").
		Where("id = ? AND status = ?", id, models.ListingStatusApproved).
		First(&listing).Error; err != nil {
		utils.NotFound(c, "RES_LISTING_NOT_FOUND", "Listing properti tidak ditemukan.")
		return
	}

	utils.OK(c, toPropertyDetail(listing))
}

func FeaturedProperties(c *gin.Context) {
	city := c.Query("city")
	limit := 6
	if l := c.Query("limit"); l != "" {
		// parse limit, clamp to 12
	}

	var listings []models.PropertyListing
	query := database.DB.
		Preload("Photos", "sort_order = ?", 0).
		Preload("Salesman").
		Preload("Tenant").
		Where("status = ?", models.ListingStatusApproved)

	if city != "" {
		query = query.Where("city = ?", city)
	}

	if err := query.Order("created_at DESC").Limit(limit).Find(&listings).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	cards := make([]dto.PropertyCardResponse, 0, len(listings))
	for _, l := range listings {
		cards = append(cards, toPropertyCard(l))
	}

	loc := city
	if loc == "" {
		loc = "Semua Kota"
	}

	utils.OK(c, gin.H{
		"location":   loc,
		"properties": cards,
	})
}

func NearbyProperties(c *gin.Context) {
	lat := c.Query("latitude")
	lng := c.Query("longitude")
	radius := c.DefaultQuery("radius_km", "10")
	limitStr := c.DefaultQuery("limit", "10")

	if lat == "" || lng == "" {
		utils.BadRequest(c, "Parameter latitude dan longitude diperlukan.")
		return
	}

	// Parse limit
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
		limit = l
	}

	// Use earthdistance extension for geospatial query
	// earth_box creates a bounding box, earth_distance calculates actual distance
	query := database.DB.
		Preload("Photos", "sort_order = ?", 0).
		Preload("Salesman").
		Preload("Tenant").
		Table("property_listings").
		Select(`property_listings.*, 
			earth_distance(
				ll_to_earth(?, ?), 
				ll_to_earth(COALESCE(latitude, 0), COALESCE(longitude, 0))
			) AS distance`, lat, lng).
		Where("status = ?", models.ListingStatusApproved).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL").
		Where(`earth_box(ll_to_earth(?, ?), ? * 1000) @> ll_to_earth(latitude, longitude)`, lat, lng, radius).
		Order("distance ASC").
		Limit(limit)

	var results []struct {
		models.PropertyListing
		Distance float64 `json:"distance"`
	}

	if err := query.Scan(&results).Error; err != nil {
		// Fallback: if earthdistance query fails (e.g., extension not loaded),
		// return empty array gracefully
		log.Printf("[WARN] Nearby query failed: %v", err)
		utils.OK(c, []dto.NearbyPropertyResponse{})
		return
	}

	properties := make([]dto.NearbyPropertyResponse, 0, len(results))
	for _, r := range results {
		np := dto.NearbyPropertyResponse{
			ID:           r.ID.String(),
			Title:        r.Title,
			Price:        formatPrice(r.Price),
			ListingType:  r.ListingType,
			PropertyType: r.PropertyType,
			City:         r.City,
			DistanceKm:   math.Round(r.Distance*10) / 10, // Round to 1 decimal
		}
		if len(r.Photos) > 0 {
			np.MainPhotoURL = r.Photos[0].MediumURL
			if np.MainPhotoURL == nil {
				np.MainPhotoURL = &r.Photos[0].WatermarkedURL
			}
		}
		if r.Salesman != nil {
			np.Salesman = dto.AgentBrief{
				ID: r.Salesman.ID.String(), Name: r.Salesman.Name,
				PhotoURL: r.Salesman.PhotoURL, Phone: r.Salesman.Phone,
			}
		}
		if r.Tenant != nil {
			np.Tenant = dto.TenantBrief{
				ID: r.Tenant.ID.String(), Name: r.Tenant.OrganizationName,
				LogoURL: r.Tenant.LogoURL,
			}
		}
		properties = append(properties, np)
	}

	utils.OK(c, properties)
}

func ListCities(c *gin.Context) {
	var cities []struct {
		City     string `json:"city"`
		Province string `json:"province"`
	}

	database.DB.Model(&models.PropertyListing{}).
		Select("DISTINCT city, province").
		Where("status = ? AND city IS NOT NULL", models.ListingStatusApproved).
		Order("city ASC").
		Find(&cities)

	utils.OK(c, cities)
}

// ==================== MAPPING HELPERS ====================

func toPropertyCard(l models.PropertyListing) dto.PropertyCardResponse {
	card := dto.PropertyCardResponse{
		ID:           l.ID.String(),
		Title:        l.Title,
		Price:        formatPrice(l.Price),
		ListingType:  l.ListingType,
		PropertyType: l.PropertyType,
		SourceType:   l.SourceType,
		City:         l.City,
		Province:     l.Province,
		Status:       l.Status,
		CreatedAt:    l.CreatedAt,
	}

	if l.LandArea != nil {
		s := formatDecimal(*l.LandArea)
		card.LandArea = &s
	}
	if l.BuildingArea != nil {
		s := formatDecimal(*l.BuildingArea)
		card.BuildingArea = &s
	}
	card.Bedrooms = l.Bedrooms
	card.Bathrooms = l.Bathrooms

	// Main photo — prefer WatermarkedURL (always present, not null)
	if len(l.Photos) > 0 {
		card.MainPhotoURL = &l.Photos[0].WatermarkedURL
	}

	// Salesman
	if l.Salesman != nil {
		card.Salesman = dto.AgentBrief{
			ID: l.Salesman.ID.String(), Name: l.Salesman.Name,
			PhotoURL: l.Salesman.PhotoURL, Phone: l.Salesman.Phone,
		}
	}

	// Tenant
	if l.Tenant != nil {
		card.Tenant = dto.TenantBrief{
			ID: l.Tenant.ID.String(), Name: l.Tenant.OrganizationName,
			LogoURL: l.Tenant.LogoURL,
		}
	}

	return card
}

func toPropertyDetail(l models.PropertyListing) gin.H {
	data := gin.H{
		"id":               l.ID,
		"title":            l.Title,
		"description":      l.Description,
		"price":            formatPrice(l.Price),
		"listing_type":     l.ListingType,
		"property_type":    l.PropertyType,
		"source_type":      l.SourceType,
		"address":          l.Address,
		"city":             l.City,
		"province":         l.Province,
		"latitude":         formatDecimalPtr(l.Latitude),
		"longitude":        formatDecimalPtr(l.Longitude),
		"land_area":        formatDecimalPtr(l.LandArea),
		"building_area":    formatDecimalPtr(l.BuildingArea),
		"bedrooms":         l.Bedrooms,
		"bathrooms":        l.Bathrooms,
		"floors":           l.Floors,
		"certificate_type": l.CertificateType,
		"facilities":       l.Facilities,
		"status":           l.Status,
		"created_at":       l.CreatedAt,
	}

	// Photos
	photos := make([]gin.H, 0)
	for _, p := range l.Photos {
		photos = append(photos, gin.H{
			"id":              p.ID,
			"original_url":    p.OriginalURL,
			"thumbnail_url":   p.ThumbnailURL,
			"medium_url":      p.MediumURL,
			"watermarked_url": p.WatermarkedURL,
			"sort_order":      p.SortOrder,
		})
	}
	data["photos"] = photos

	// Salesman
	if l.Salesman != nil {
		data["salesman"] = gin.H{
			"id":        l.Salesman.ID,
			"name":      l.Salesman.Name,
			"photo_url": l.Salesman.PhotoURL,
			"phone":     l.Salesman.Phone,
		}
	}

	// Tenant
	if l.Tenant != nil {
		data["tenant"] = gin.H{
			"id":       l.Tenant.ID,
			"name":     l.Tenant.OrganizationName,
			"logo_url": l.Tenant.LogoURL,
			"phone":    l.Tenant.Phone,
		}
	}

	return data
}

func formatPrice(price float64) string {
	return formatDecimal(price)
}

func formatDecimal(val float64) string {
	return fmt.Sprintf("%.2f", val)
}

func formatDecimalPtr(val *float64) *string {
	if val == nil {
		return nil
	}
	s := formatDecimal(*val)
	return &s
}
