package routes

import (
	"property-hub-backend/config"
	"property-hub-backend/handlers"
	"property-hub-backend/middleware"
	"property-hub-backend/models"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORSMiddleware(cfg))
	r.Use(middleware.RateLimitGlobal(cfg))

	// Static file server for uploads (photos)
	r.Static("/uploads", cfg.UploadDir)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "app": cfg.AppName})
	})

	api := r.Group("/api/v1")

	// ==================== PUBLIC ====================
	auth := api.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimitLogin(cfg), handlers.Register)
		auth.POST("/login", middleware.RateLimitLogin(cfg), handlers.Login)
	}

	properties := api.Group("/properties")
	{
		properties.GET("", handlers.ListProperties)
		properties.GET("/featured", handlers.FeaturedProperties)
		properties.GET("/nearby", handlers.NearbyProperties)
		properties.GET("/:id", handlers.GetPropertyDetail)
		properties.POST("/:id/view", handlers.LogPropertyView)
	}

	api.GET("/property-types", handlers.ListPropertyTypes)
	api.GET("/facilities", handlers.ListFacilities)

	locations := api.Group("/locations")
	{
		locations.GET("", handlers.ListLocations)
	}

	// ==================== AUTHENTICATED COMMON ====================
	me := api.Group("/me")
	me.Use(middleware.AuthRequired(cfg))
	{
		me.GET("/profile", handlers.GetMyProfile)
		me.PUT("/profile", handlers.UpdateMyProfile)
	}

	// ==================== BUYER ====================
	buyer := api.Group("/buyer")
	buyer.Use(middleware.AuthRequired(cfg))
	buyer.Use(middleware.RequireRole(models.RoleBuyer))
	{
		buyer.GET("/favorites", handlers.BuyerListSaved)
		buyer.POST("/favorites/:propertyId", handlers.BuyerSaveProperty)
		buyer.DELETE("/favorites/:propertyId", handlers.BuyerRemoveSaved)
		buyer.POST("/inquiries", handlers.BuyerCreateInquiry)
		buyer.GET("/inquiries", handlers.BuyerListInquiries)
	}

	// Keep old /me/saved for backward compatibility
	meSaved := api.Group("/me")
	meSaved.Use(middleware.AuthRequired(cfg))
	meSaved.Use(middleware.RequireRole(models.RoleBuyer))
	{
		meSaved.GET("/saved", handlers.BuyerListSaved)
		meSaved.POST("/saved/:propertyId", handlers.BuyerSaveProperty)
		meSaved.DELETE("/saved/:propertyId", handlers.BuyerRemoveSaved)
	}

	// ==================== SALESMAN ====================
	salesman := api.Group("/salesman")
	salesman.Use(middleware.AuthRequired(cfg))
	salesman.Use(middleware.RequireRole(models.RoleSalesman, models.RoleTenantAdmin))
	salesman.Use(middleware.TenantScope())
	{
		salesman.GET("/dashboard", handlers.SalesmanDashboard)
		salesman.GET("/listings", handlers.SalesmanListListings)
		salesman.POST("/listings", handlers.SalesmanCreateListing)
		salesman.GET("/listings/:id", handlers.SalesmanGetListing)
		salesman.PUT("/listings/:id", handlers.SalesmanUpdateListing)
		salesman.DELETE("/listings/:id", handlers.SalesmanDeleteListing)
		salesman.POST("/listings/:id/submit", handlers.SalesmanSubmitListing)
		salesman.POST("/listings/:id/deactivate", handlers.SalesmanDeactivateListing)
		salesman.POST("/listings/:id/mark-sold", handlers.SalesmanMarkSold)
		salesman.POST("/listings/:id/mark-rented", handlers.SalesmanMarkRented)
		salesman.POST("/listings/:id/photos", handlers.SalesmanUploadPhotos)
		salesman.DELETE("/listings/:id/photos/:photoId", handlers.SalesmanDeletePhoto)
		salesman.PUT("/listings/:id/photos/reorder", handlers.SalesmanReorderPhotos)
		salesman.GET("/quota", handlers.SalesmanGetQuota)
		salesman.GET("/inquiries", handlers.SalesmanListInquiries)
		salesman.PUT("/inquiries/:id", handlers.SalesmanUpdateInquiry)
	}

	// ==================== TENANT ADMIN ====================
	tenant := api.Group("/tenant")
	tenant.Use(middleware.AuthRequired(cfg))
	tenant.Use(middleware.RequireRole(models.RoleTenantAdmin))
	tenant.Use(middleware.TenantScope())
	{
		tenant.GET("/dashboard", handlers.TenantDashboard)
		tenant.GET("/profile", handlers.TenantGetProfile)
		tenant.PUT("/profile", handlers.TenantUpdateProfile)
		tenant.GET("/salesmen", handlers.TenantListSalesmen)
		tenant.POST("/salesmen", handlers.TenantAddSalesman)
		tenant.DELETE("/salesmen/:id", handlers.TenantRemoveSalesman)
		tenant.GET("/listings", handlers.TenantListListings)
		tenant.GET("/subscription", handlers.TenantGetSubscription)
		tenant.POST("/subscription/upgrade", handlers.TenantRequestUpgrade)
		tenant.GET("/inquiries", handlers.TenantListInquiries)
	}

	// ==================== PLATFORM ADMIN ====================
	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg))
	admin.Use(middleware.RequireRole(models.RolePlatformAdmin))
	{
		admin.GET("/dashboard", handlers.AdminDashboard)
		admin.GET("/tenants", handlers.AdminListTenants)
		admin.POST("/tenants", handlers.AdminCreateTenant)
		admin.GET("/tenants/:id", handlers.AdminGetTenant)
		admin.POST("/tenants/:id/suspend", handlers.AdminSuspendTenant)
		admin.POST("/tenants/:id/activate", handlers.AdminActivateTenant)
		admin.PUT("/tenants/:id/plan", handlers.AdminChangePlan)
		admin.GET("/listings/pending", handlers.AdminListPending)
		admin.POST("/listings/:id/approve", handlers.AdminApproveListing)
		admin.POST("/listings/:id/reject", handlers.AdminRejectListing)
		admin.GET("/audit-logs", handlers.AdminAuditLogs)
		admin.GET("/listings", handlers.AdminListAllListings)
		// Master data CRUD
		admin.GET("/subscription-plans", handlers.AdminListSubscriptionPlans)
		admin.POST("/subscription-plans", handlers.AdminCreateSubscriptionPlan)
		admin.PUT("/subscription-plans/:id", handlers.AdminUpdateSubscriptionPlan)
		admin.GET("/property-types", handlers.AdminListPropertyTypes)
		admin.POST("/property-types", handlers.AdminCreatePropertyType)
		admin.PUT("/property-types/:id", handlers.AdminUpdatePropertyType)
		admin.GET("/facilities", handlers.AdminListFacilities)
		admin.POST("/facilities", handlers.AdminCreateFacility)
		admin.PUT("/facilities/:id", handlers.AdminUpdateFacility)
		admin.GET("/locations", handlers.AdminListLocations)
		admin.POST("/locations", handlers.AdminCreateLocation)
		admin.PUT("/locations/:id", handlers.AdminUpdateLocation)
		admin.PUT("/tenants/:id/subscription", handlers.AdminChangePlanByID)
	}

	return r
}
