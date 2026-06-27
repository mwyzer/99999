package dto

import "time"

// ==================== AUTH ====================

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Phone    string `json:"phone" binding:"required,min=8,max=20"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  UserBrief `json:"user"`
}

type UserBrief struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Phone      *string `json:"phone"`
	PhotoURL   *string `json:"photo_url"`
	Role       string  `json:"role"`
	TenantID   *string `json:"tenant_id"`
	TenantName *string `json:"tenant_name"`
}

// ==================== PROFILE ====================

type UpdateProfileRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=2,max=100"`
	Phone *string `json:"phone" binding:"omitempty,min=8,max=20"`
}

// ==================== LISTING ====================

type CreateListingRequest struct {
	Title           string                 `json:"title" binding:"required,min=5,max=300"`
	Description     *string                `json:"description" binding:"omitempty,max=5000"`
	Price           float64                `json:"price" binding:"required,gt=0"`
	ListingType     string                 `json:"listing_type" binding:"required,oneof=sale rent"`
	PropertyType    string                 `json:"property_type" binding:"required,oneof=house land apartment shophouse warehouse office villa"`
	SourceType      string                 `json:"source_type" binding:"omitempty,oneof=regular bank_auction company_asset"`
	Address         *string                `json:"address" binding:"omitempty,max=500"`
	City            *string                `json:"city" binding:"omitempty,min=2,max=100"`
	Province        *string                `json:"province" binding:"omitempty,min=2,max=100"`
	Latitude        *float64               `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude       *float64               `json:"longitude" binding:"omitempty,min=-180,max=180"`
	LandArea        *float64               `json:"land_area" binding:"omitempty,gt=0"`
	BuildingArea    *float64               `json:"building_area" binding:"omitempty,gt=0"`
	Bedrooms        *int                   `json:"bedrooms" binding:"omitempty,min=0,max=99"`
	Bathrooms       *int                   `json:"bathrooms" binding:"omitempty,min=0,max=99"`
	Floors          *int                   `json:"floors" binding:"omitempty,min=0,max=200"`
	CertificateType *string                `json:"certificate_type" binding:"omitempty,oneof=SHM SHGB Girik Lainnya"`
	Facilities      map[string]interface{} `json:"facilities"`
}

type UpdateListingRequest struct {
	Title           *string                 `json:"title" binding:"omitempty,min=5,max=300"`
	Description     *string                 `json:"description" binding:"omitempty,max=5000"`
	Price           *float64                `json:"price" binding:"omitempty,gt=0"`
	ListingType     *string                 `json:"listing_type" binding:"omitempty,oneof=sale rent"`
	PropertyType    *string                 `json:"property_type" binding:"omitempty,oneof=house land apartment shophouse warehouse office villa"`
	SourceType      *string                 `json:"source_type" binding:"omitempty,oneof=regular bank_auction company_asset"`
	Address         *string                 `json:"address"`
	City            *string                 `json:"city"`
	Province        *string                 `json:"province"`
	Latitude        *float64                `json:"latitude"`
	Longitude       *float64                `json:"longitude"`
	LandArea        *float64                `json:"land_area"`
	BuildingArea    *float64                `json:"building_area"`
	Bedrooms        *int                    `json:"bedrooms"`
	Bathrooms       *int                    `json:"bathrooms"`
	Floors          *int                    `json:"floors"`
	CertificateType *string                 `json:"certificate_type"`
	Facilities      map[string]interface{}  `json:"facilities"`
}

type ReorderPhotosRequest struct {
	PhotoIDs []string `json:"photo_ids" binding:"required,min=1"`
}

// ==================== TENANT ADMIN ====================

type UpdateTenantRequest struct {
	OrganizationName *string `json:"organization_name" binding:"omitempty,min=2,max=200"`
	Description      *string `json:"description" binding:"omitempty,max=2000"`
	Phone            *string `json:"phone" binding:"omitempty,min=8,max=20"`
	Address          *string `json:"address" binding:"omitempty,max=500"`
}

type AddSalesmanRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Phone    string `json:"phone" binding:"required,min=8,max=20"`
	Password string `json:"password" binding:"required,min=8"`
}

type RequestUpgradeRequest struct {
	PlanType string `json:"plan_type" binding:"required,oneof=premium"`
}

// ==================== PLATFORM ADMIN ====================

type CreateTenantRequest struct {
	OrganizationName string `json:"organization_name" binding:"required,min=2,max=200"`
	SubdomainSlug    string `json:"subdomain_slug" binding:"required,min=3,max=100"`
	AdminName        string `json:"admin_name" binding:"required,min=2,max=100"`
	AdminEmail       string `json:"admin_email" binding:"required,email,max=255"`
	AdminPhone       string `json:"admin_phone" binding:"required,min=8,max=20"`
	AdminPassword    string `json:"admin_password" binding:"required,min=8"`
	PlanType         string `json:"plan_type" binding:"omitempty,oneof=free premium"`
}

type ChangePlanRequest struct {
	PlanType string `json:"plan_type" binding:"required,oneof=free premium"`
}

type RejectListingRequest struct {
	Reason string `json:"reason" binding:"required,min=10,max=500"`
}

// ==================== PUBLIC RESPONSES ====================

type PropertyCardResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Price         string  `json:"price"`
	ListingType   string  `json:"listing_type"`
	PropertyType  string  `json:"property_type"`
	SourceType    string  `json:"source_type"`
	City          *string `json:"city"`
	Province      *string `json:"province"`
	LandArea      *string `json:"land_area"`
	BuildingArea  *string `json:"building_area"`
	Bedrooms      *int    `json:"bedrooms"`
	Bathrooms     *int    `json:"bathrooms"`
	Status        string  `json:"status"`
	MainPhotoURL  *string `json:"main_photo_url"`
	Salesman      AgentBrief `json:"salesman"`
	Tenant        TenantBrief `json:"tenant"`
	CreatedAt     time.Time `json:"created_at"`
}

type AgentBrief struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	PhotoURL *string `json:"photo_url"`
	Phone    *string `json:"phone"`
}

type TenantBrief struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo_url"`
}

type NearbyPropertyResponse struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Price        string  `json:"price"`
	ListingType  string  `json:"listing_type"`
	PropertyType string  `json:"property_type"`
	City         *string `json:"city"`
	DistanceKm   float64 `json:"distance_km"`
	MainPhotoURL *string `json:"main_photo_url"`
	Salesman     AgentBrief `json:"salesman"`
	Tenant       TenantBrief `json:"tenant"`
}
