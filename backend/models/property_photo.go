package models

import (
	"time"

	"github.com/google/uuid"
)

type PropertyPhoto struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ListingID       uuid.UUID `gorm:"type:uuid;not null;index" json:"listing_id"`
	FileName        *string   `gorm:"type:varchar(255)" json:"file_name"`
	OriginalURL     string    `gorm:"type:varchar(500);not null" json:"original_url"`
	ThumbnailURL    *string   `gorm:"type:varchar(500)" json:"thumbnail_url"`
	MediumURL       *string   `gorm:"type:varchar(500)" json:"medium_url"`
	WatermarkedURL  string    `gorm:"type:varchar(500);not null" json:"watermarked_url"`
	WatermarkStatus string    `gorm:"type:varchar(20);not null;default:pending" json:"watermark_status"`
	IsPrimary       bool      `gorm:"not null;default:false" json:"is_primary"`
	SortOrder       int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Listing *PropertyListing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}

func (PropertyPhoto) TableName() string {
	return "property_photos"
}
