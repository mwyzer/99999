package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AuctionStatusUpcoming  = "upcoming"
	AuctionStatusOpen      = "open"
	AuctionStatusClosed    = "closed"
	AuctionStatusCancelled = "cancelled"
	AuctionStatusSold      = "sold"
)

type BankAuctionDetail struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PropertyID         uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"property_id"`
	BankName           string     `gorm:"type:varchar(200);not null" json:"bank_name"`
	AuctionNumber      *string    `gorm:"type:varchar(100)" json:"auction_number"`
	AuctionLimitPrice  *float64   `gorm:"type:decimal(16,2)" json:"auction_limit_price"`
	AuctionDeposit     *float64   `gorm:"type:decimal(16,2)" json:"auction_deposit"`
	AuctionDate        *time.Time `json:"auction_date"`
	AuctionLocation    *string    `gorm:"type:text" json:"auction_location"`
	AuctionDocumentURL *string    `gorm:"type:varchar(500)" json:"auction_document_url"`
	AuctionStatus      string     `gorm:"type:varchar(20);not null;default:upcoming" json:"auction_status"`
	Notes              *string    `gorm:"type:text" json:"notes"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Property *PropertyListing `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (BankAuctionDetail) TableName() string { return "bank_auction_details" }
