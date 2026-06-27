package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	DisposalTypeSale  = "sale"
	DisposalTypeRent  = "rent"
	DisposalTypeLease = "lease"

	AssetStatusAvailable   = "available"
	AssetStatusUnderReview = "under_review"
	AssetStatusSold        = "sold"
	AssetStatusRented      = "rented"
	AssetStatusInactive    = "inactive"
)

type CompanyAssetDetail struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PropertyID         uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"property_id"`
	CompanyName        string     `gorm:"type:varchar(200);not null" json:"company_name"`
	CompanyAssetCode   *string    `gorm:"type:varchar(100)" json:"company_asset_code"`
	DisposalType       string     `gorm:"type:varchar(20);not null" json:"disposal_type"`
	AssetStatus        string     `gorm:"type:varchar(20);not null;default:available" json:"asset_status"`
	PICName            *string    `gorm:"type:varchar(200)" json:"pic_name"`
	PICPhone           *string    `gorm:"type:varchar(20)" json:"pic_phone"`
	PICWhatsappNumber  *string    `gorm:"type:varchar(20)" json:"pic_whatsapp_number"`
	DocumentURL        *string    `gorm:"type:varchar(500)" json:"document_url"`
	InternalNote       *string    `gorm:"type:text" json:"internal_note"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Property *PropertyListing `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (CompanyAssetDetail) TableName() string { return "company_asset_details" }
