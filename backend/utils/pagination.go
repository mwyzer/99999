package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams extracts and validates pagination from query
type PaginationParams struct {
	Page    int
	PerPage int
	Offset  int
}

func GetPagination(c *gin.Context) PaginationParams {
	page := getQueryInt(c, "page", 1)
	perPage := getQueryInt(c, "per_page", 20)

	// Clamp
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 1
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}

func getQueryInt(c *gin.Context, key string, fallback int) int {
	val := c.Query(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}

// RoundTo2Dec rounds a float to 2 decimal places
func RoundTo2Dec(val float64) float64 {
	return math.Round(val*100) / 100
}
