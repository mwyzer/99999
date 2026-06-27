package models

import (
	"time"

	"github.com/google/uuid"
)

type SavedProperty struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BuyerID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_saved_buyer_listing;index" json:"buyer_id"`
	ListingID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_saved_buyer_listing;index" json:"listing_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Buyer   *User            `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Listing *PropertyListing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}

func (SavedProperty) TableName() string {
	return "saved_properties"
}
