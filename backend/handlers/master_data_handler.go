package handlers

import (
	"property-hub-backend/database"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
)

// ==================== PROPERTY TYPES (public + admin) ====================

func ListPropertyTypes(c *gin.Context) {
	var types []models.PropertyType
	database.DB.Where("is_active = ?", true).Order("name ASC").Find(&types)
	utils.OK(c, types)
}

func AdminListPropertyTypes(c *gin.Context) {
	var types []models.PropertyType
	database.DB.Order("name ASC").Find(&types)
	utils.OK(c, types)
}

func AdminCreatePropertyType(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	pt := models.PropertyType{Name: req.Name, Slug: req.Slug}
	if req.Description != "" {
		pt.Description = &req.Description
	}
	if err := database.DB.Create(&pt).Error; err != nil {
		utils.InternalError(c, "Gagal membuat tipe properti.")
		return
	}
	utils.Created(c, pt)
}

func AdminUpdatePropertyType(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		IsActive    *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	database.DB.Model(&models.PropertyType{}).Where("id = ?", id).Updates(updates)
	utils.OK(c, gin.H{"message": "Tipe properti diperbarui."})
}

// ==================== FACILITIES (public + admin) ====================

func ListFacilities(c *gin.Context) {
	var facilities []models.Facility
	database.DB.Where("is_active = ?", true).Order("name ASC").Find(&facilities)
	utils.OK(c, facilities)
}

func AdminListFacilities(c *gin.Context) {
	var facilities []models.Facility
	database.DB.Order("name ASC").Find(&facilities)
	utils.OK(c, facilities)
}

func AdminCreateFacility(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	f := models.Facility{Name: req.Name}
	if req.Icon != "" {
		f.Icon = &req.Icon
	}
	if err := database.DB.Create(&f).Error; err != nil {
		utils.InternalError(c, "Gagal membuat fasilitas.")
		return
	}
	utils.Created(c, f)
}

func AdminUpdateFacility(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name     *string `json:"name"`
		Icon     *string `json:"icon"`
		IsActive *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	database.DB.Model(&models.Facility{}).Where("id = ?", id).Updates(updates)
	utils.OK(c, gin.H{"message": "Fasilitas diperbarui."})
}

// ==================== LOCATIONS (public + admin) ====================

func ListLocations(c *gin.Context) {
	var locations []models.Location
	query := database.DB.Where("is_active = ?", true)

	if province := c.Query("province"); province != "" {
		query = query.Where("province ILIKE ?", "%"+province+"%")
	}
	if q := c.Query("q"); q != "" {
		query = query.Where("city ILIKE ? OR province ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	query.Order("province ASC, city ASC").Limit(50).Find(&locations)
	utils.OK(c, locations)
}

func AdminListLocations(c *gin.Context) {
	var locations []models.Location
	database.DB.Order("province ASC, city ASC").Find(&locations)
	utils.OK(c, locations)
}

func AdminCreateLocation(c *gin.Context) {
	var req struct {
		City      string  `json:"city" binding:"required"`
		Province  string  `json:"province" binding:"required"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	loc := models.Location{City: req.City, Province: req.Province}
	if req.Country != "" {
		loc.Country = req.Country
	}
	if req.Latitude != 0 {
		loc.Latitude = &req.Latitude
	}
	if req.Longitude != 0 {
		loc.Longitude = &req.Longitude
	}

	if err := database.DB.Create(&loc).Error; err != nil {
		utils.InternalError(c, "Gagal membuat lokasi.")
		return
	}
	utils.Created(c, loc)
}

func AdminUpdateLocation(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		City      *string  `json:"city"`
		Province  *string  `json:"province"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		IsActive  *bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.Province != nil {
		updates["province"] = *req.Province
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	database.DB.Model(&models.Location{}).Where("id = ?", id).Updates(updates)
	utils.OK(c, gin.H{"message": "Lokasi diperbarui."})
}

// ==================== SUBSCRIPTION PLANS (admin) ====================

func AdminListSubscriptionPlans(c *gin.Context) {
	var plans []models.SubscriptionPlan
	database.DB.Order("created_at ASC").Find(&plans)
	utils.OK(c, plans)
}

func AdminCreateSubscriptionPlan(c *gin.Context) {
	var req struct {
		Name                   string `json:"name" binding:"required"`
		Slug                   string `json:"slug" binding:"required"`
		MaxSalesmen            int    `json:"max_salesmen"`
		MaxListingsPerSalesman int    `json:"max_listings_per_salesman"`
		Description            string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	plan := models.SubscriptionPlan{
		Name:                   req.Name,
		Slug:                   req.Slug,
		MaxSalesmen:            req.MaxSalesmen,
		MaxListingsPerSalesman: req.MaxListingsPerSalesman,
	}
	if req.Description != "" {
		plan.Description = &req.Description
	}
	if plan.MaxSalesmen == 0 {
		plan.MaxSalesmen = 5
	}
	if plan.MaxListingsPerSalesman == 0 {
		plan.MaxListingsPerSalesman = 5
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		utils.InternalError(c, "Gagal membuat paket.")
		return
	}
	utils.Created(c, plan)
}

func AdminUpdateSubscriptionPlan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name                   *string `json:"name"`
		Slug                   *string `json:"slug"`
		MaxSalesmen            *int    `json:"max_salesmen"`
		MaxListingsPerSalesman *int    `json:"max_listings_per_salesman"`
		Description            *string `json:"description"`
		IsActive               *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Field tidak valid.", nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.MaxSalesmen != nil {
		updates["max_salesmen"] = *req.MaxSalesmen
	}
	if req.MaxListingsPerSalesman != nil {
		updates["max_listings_per_salesman"] = *req.MaxListingsPerSalesman
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	database.DB.Model(&models.SubscriptionPlan{}).Where("id = ?", id).Updates(updates)
	utils.OK(c, gin.H{"message": "Paket diperbarui."})
}
