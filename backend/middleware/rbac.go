package middleware

import (
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole creates middleware that checks for specific roles
func RequireRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool)
	for _, r := range roles {
		roleSet[r] = true
	}

	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "AUTHZ_FORBIDDEN", "Anda tidak memiliki izin untuk mengakses resource ini.")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || !roleSet[roleStr] {
			utils.Forbidden(c, "AUTHZ_FORBIDDEN", "Anda tidak memiliki izin untuk mengakses resource ini.")
			c.Abort()
			return
		}

		c.Next()
	}
}
