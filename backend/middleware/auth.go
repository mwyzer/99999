package middleware

import (
	"strings"

	"property-hub-backend/config"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired validates JWT and injects claims into context
func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "AUTH_TOKEN_MISSING", "Token autentikasi diperlukan. Silakan login terlebih dahulu.")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Unauthorized(c, "AUTH_TOKEN_INVALID", "Token autentikasi tidak valid.")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], cfg)
		if err != nil {
			utils.Unauthorized(c, "AUTH_TOKEN_EXPIRED", "Sesi Anda telah berakhir. Silakan login kembali.")
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("tenantID", claims.TenantID)

		c.Next()
	}
}
