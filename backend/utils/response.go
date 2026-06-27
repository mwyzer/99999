package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response helpers — consistent JSON envelope

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   APIError `json:"error"`
}

type APIError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []FieldError  `json:"details"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Success responses

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessPaginated(c *gin.Context, status int, data interface{}, meta Meta) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

func Created(c *gin.Context, data interface{}) {
	Success(c, http.StatusCreated, data)
}

func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error responses

func Error(c *gin.Context, status int, code string, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Success: false,
		Error: APIError{
			Code:    code,
			Message: message,
			Details: []FieldError{},
		},
	})
}

func ValidationError(c *gin.Context, message string, details []FieldError) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, ErrorResponse{
		Success: false,
		Error: APIError{
			Code:    "VAL_INPUT_INVALID",
			Message: message,
			Details: details,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "VAL_JSON_INVALID", message)
}

func Unauthorized(c *gin.Context, code string, message string) {
	Error(c, http.StatusUnauthorized, code, message)
}

func Forbidden(c *gin.Context, code string, message string) {
	Error(c, http.StatusForbidden, code, message)
}

func NotFound(c *gin.Context, code string, message string) {
	Error(c, http.StatusNotFound, code, message)
}

func Conflict(c *gin.Context, code string, message string) {
	Error(c, http.StatusConflict, code, message)
}

func Unprocessable(c *gin.Context, code string, message string) {
	Error(c, http.StatusUnprocessableEntity, code, message)
}

func TooManyRequests(c *gin.Context, code string, message string) {
	Error(c, http.StatusTooManyRequests, code, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "SRV_INTERNAL_ERROR", message)
}

// CalculateMeta computes pagination metadata
func CalculateMeta(page, perPage int, total int64) Meta {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
