package middleware

import (
	"property-hub-backend/database"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantScope enforces tenant_id isolation on tenant-scoped queries.
// - platform_admin: bypass (access all tenants)
// - buyer: no tenant restriction (null tenant_id)
// - salesman / tenant_admin: enforce tenant_id from JWT
func TenantScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")

		// Platform admin bypasses tenant scope
		if role == models.RolePlatformAdmin {
			c.Next()
			return
		}

		// Buyer has no tenant — no scope to enforce
		if role == models.RoleBuyer {
			c.Next()
			return
		}

		// Salesman / Tenant Admin: verify tenant is active
		tenantIDVal, exists := c.Get("tenantID")
		if !exists || tenantIDVal == nil {
			utils.Forbidden(c, "AUTHZ_FORBIDDEN", "Akun Anda tidak terhubung dengan organisasi manapun.")
			c.Abort()
			return
		}

		// Check if tenant is suspended
		var tenant models.Tenant
		if err := database.DB.First(&tenant, "id = ?", tenantIDVal).Error; err != nil {
			utils.Forbidden(c, "AUTH_ACCOUNT_SUSPENDED", "Organisasi tidak ditemukan.")
			c.Abort()
			return
		}

		if tenant.Status == models.TenantStatusSuspended {
			utils.Forbidden(c, "AUTH_ACCOUNT_SUSPENDED", "Akun organisasi Anda sedang dinonaktifkan. Hubungi administrator.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetTenantID extracts tenant_id from context (returns uuid.Nil for buyer/platform_admin)
func GetTenantID(c *gin.Context) uuid.UUID {
	val, exists := c.Get("tenantID")
	if !exists || val == nil {
		return uuid.Nil
	}
	tid, ok := val.(*uuid.UUID)
	if !ok || tid == nil {
		return uuid.Nil
	}
	return *tid
}

// GetUserID extracts user_id from context
func GetUserID(c *gin.Context) uuid.UUID {
	val, _ := c.Get("userID")
	uid, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return uid
}

// GetRole extracts role from context
func GetRole(c *gin.Context) string {
	val, _ := c.Get("role")
	role, ok := val.(string)
	if !ok {
		return ""
	}
	return role
}
