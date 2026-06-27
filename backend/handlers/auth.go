package handlers

import (
	"fmt"
	"net/http"

	"property-hub-backend/config"
	"property-hub-backend/database"
	"property-hub-backend/dto"
	"property-hub-backend/models"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var AppConfig *config.Config

// ==================== AUTH ====================

func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	// Check duplicate email
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Conflict(c, "AUTH_EMAIL_REGISTERED", "Email sudah terdaftar. Silakan gunakan email lain atau login.")
		return
	}

	hash, err := utils.HashPassword(req.Password, AppConfig.BcryptCost)
	if err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
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
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	utils.Created(c, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", "Email atau password salah.")
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", "Email atau password salah.")
		return
	}

	// Check user status
	if user.Status == models.UserStatusInactive {
		utils.Forbidden(c, "AUTH_ACCOUNT_INACTIVE", "Akun Anda tidak aktif. Hubungi administrator.")
		return
	}

	// Check tenant status if applicable
	if user.TenantID != nil {
		var tenant models.Tenant
		if err := database.DB.First(&tenant, "id = ?", user.TenantID).Error; err == nil {
			if tenant.Status == models.TenantStatusSuspended {
				utils.Forbidden(c, "AUTH_ACCOUNT_SUSPENDED", "Akun organisasi Anda sedang dinonaktifkan. Hubungi administrator.")
				return
			}
		}
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Role, user.TenantID, AppConfig)
	if err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	// Build response
	resp := dto.LoginResponse{
		Token: token,
		User: dto.UserBrief{
			ID:       user.ID.String(),
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			PhotoURL: user.PhotoURL,
			Role:     user.Role,
		},
	}

	if user.TenantID != nil {
		tid := user.TenantID.String()
		resp.User.TenantID = &tid

		var tenant models.Tenant
		if err := database.DB.First(&tenant, "id = ?", user.TenantID).Error; err == nil {
			resp.User.TenantName = &tenant.OrganizationName
		}
	}

	utils.OK(c, resp)
}

// ==================== PROFILE ====================

func GetMyProfile(c *gin.Context) {
	userID := c.MustGet("userID")

	var user models.User
	if err := database.DB.Preload("Tenant").First(&user, "id = ?", userID).Error; err != nil {
		utils.NotFound(c, "RES_USER_NOT_FOUND", "Pengguna tidak ditemukan.")
		return
	}

	data := gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"photo_url":  user.PhotoURL,
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	}

	if user.Tenant != nil {
		data["tenant"] = gin.H{
			"id":       user.Tenant.ID,
			"name":     user.Tenant.OrganizationName,
			"logo_url": user.Tenant.LogoURL,
		}
	}

	utils.OK(c, data)
}

func UpdateMyProfile(c *gin.Context) {
	userID := c.MustGet("userID")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Beberapa field tidak valid. Silakan periksa kembali.", parseValidationErrors(err))
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Terjadi kesalahan pada server. Silakan coba lagi nanti.")
		return
	}

	GetMyProfile(c)
}

// ==================== HELPERS ====================

func parseValidationErrors(err error) []utils.FieldError {
	if err == nil {
		return []utils.FieldError{}
	}

	// Handle go-playground/validator ValidationErrors
	if ve, ok := err.(validator.ValidationErrors); ok {
		fields := make([]utils.FieldError, 0, len(ve))
		for _, fe := range ve {
			field := toSnakeCase(fe.Field())
			msg := ""
			switch fe.Tag() {
			case "required":
				msg = "Field ini wajib diisi"
			case "email":
				msg = "Format email tidak valid"
			case "min":
				msg = fmt.Sprintf("Minimal %s karakter", fe.Param())
			case "max":
				msg = fmt.Sprintf("Maksimal %s karakter", fe.Param())
			case "gt":
				msg = fmt.Sprintf("Harus lebih besar dari %s", fe.Param())
			case "oneof":
				msg = fmt.Sprintf("Nilai harus salah satu dari: %s", fe.Param())
			case "gte":
				msg = fmt.Sprintf("Harus lebih besar atau sama dengan %s", fe.Param())
			case "lte":
				msg = fmt.Sprintf("Harus kurang dari atau sama dengan %s", fe.Param())
			default:
				msg = fmt.Sprintf("Tidak valid (%s)", fe.Tag())
			}
			fields = append(fields, utils.FieldError{
				Field:   field,
				Message: msg,
			})
		}
		return fields
	}

	// Handle generic error — return empty (the main message already covers it)
	return []utils.FieldError{}
}

// toSnakeCase converts CamelCase to snake_case for field names
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, r+32) // lowercase
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func getParamUUID(c *gin.Context, param string) string {
	return c.Param(param)
}

func jsonOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func jsonCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": data})
}
