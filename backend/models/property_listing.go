package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Listing type
const (
	ListingTypeSale = "sale"
	ListingTypeRent = "rent"
)

// Property types
const (
	PropertyTypeHouse     = "house"
	PropertyTypeLand      = "land"
	PropertyTypeApartment = "apartment"
	PropertyTypeShophouse = "shophouse"
	PropertyTypeWarehouse = "warehouse"
	PropertyTypeOffice    = "office"
	PropertyTypeVilla     = "villa"
)

// Source types
const (
	SourceTypeRegular      = "regular"
	SourceTypeBankAuction  = "bank_auction"
	SourceTypeCompanyAsset = "company_asset"
)

// Certificate types
const (
	CertSHM     = "SHM"
	CertSHGB    = "SHGB"
	CertGirik   = "Girik"
	CertLainnya = "Lainnya"
)

// Listing statuses
const (
	ListingStatusDraft    = "draft"
	ListingStatusPending  = "pending"
	ListingStatusApproved = "approved"
	ListingStatusRejected = "rejected"
	ListingStatusSold     = "sold"
	ListingStatusRented   = "rented"
	ListingStatusInactive = "inactive"
	ListingStatusDeleted  = "deleted"
)

// Facilities is a JSONB field
type Facilities map[string]interface{}

func (f Facilities) Value() (driver.Value, error) {
	if f == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(f)
}

func (f *Facilities) Scan(value interface{}) error {
	if value == nil {
		*f = make(Facilities)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Facilities: not []byte")
	}
	return json.Unmarshal(bytes, f)
}

type PropertyListing struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SalesmanID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"salesman_id"`
	Title           string         `gorm:"type:varchar(300);not null" json:"title"`
	Description     *string        `gorm:"type:text" json:"description"`
	Price           float64        `gorm:"type:decimal(16,2);not null" json:"price"`
	ListingType     string         `gorm:"type:varchar(10);not null" json:"listing_type"`
	PropertyType    string         `gorm:"type:varchar(20);not null" json:"property_type"`
	SourceType      string         `gorm:"type:varchar(20);not null;default:regular" json:"source_type"`
	RentPeriod      *string        `gorm:"type:varchar(20)" json:"rent_period"`
	PropertyTypeID  *uuid.UUID     `gorm:"type:uuid" json:"property_type_id"`
	LocationID      *uuid.UUID     `gorm:"type:uuid" json:"location_id"`
	Address         *string        `gorm:"type:text" json:"address"`
	City            *string        `gorm:"type:varchar(100);index" json:"city"`
	Province        *string        `gorm:"type:varchar(100)" json:"province"`
	Latitude        *float64       `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude       *float64       `gorm:"type:decimal(10,7)" json:"longitude"`
	LandArea        *float64       `gorm:"type:decimal(12,2)" json:"land_area"`
	BuildingArea    *float64       `gorm:"type:decimal(12,2)" json:"building_area"`
	Bedrooms        *int           `json:"bedrooms"`
	Bathrooms       *int           `json:"bathrooms"`
	Floors          *int           `json:"floors"`
	CertificateType *string        `gorm:"type:varchar(20)" json:"certificate_type"`
	Facilities      Facilities     `gorm:"type:jsonb;default:'{}'" json:"facilities"`
	Status          string         `gorm:"type:varchar(20);not null;default:draft;index" json:"status"`
	RejectReason    *string        `gorm:"type:text" json:"reject_reason"`
	ApprovedBy      *uuid.UUID     `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt      *time.Time     `json:"approved_at"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Tenant              *Tenant              `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Salesman            *User                `gorm:"foreignKey:SalesmanID" json:"salesman,omitempty"`
	Approver            *User                `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	PropertyTypeRel     *PropertyType        `gorm:"foreignKey:PropertyTypeID" json:"property_type_ref,omitempty"`
	Location            *Location            `gorm:"foreignKey:LocationID" json:"location_ref,omitempty"`
	Photos              []PropertyPhoto      `gorm:"foreignKey:ListingID" json:"photos,omitempty"`
	SavedBy             []SavedProperty      `gorm:"foreignKey:ListingID" json:"-"`
	BankAuctionDetail   *BankAuctionDetail   `gorm:"foreignKey:PropertyID" json:"bank_auction_detail,omitempty"`
	CompanyAssetDetail  *CompanyAssetDetail  `gorm:"foreignKey:PropertyID" json:"company_asset_detail,omitempty"`
	PropertyFacilities  []PropertyFacility   `gorm:"foreignKey:PropertyID" json:"property_facilities,omitempty"`
	Inquiries           []Inquiry            `gorm:"foreignKey:PropertyID" json:"-"`
	Views               []PropertyView       `gorm:"foreignKey:PropertyID" json:"-"`
}

func (PropertyListing) TableName() string {
	return "property_listings"
}

// IsQuotaCounted returns true if the status counts toward quota
func (l *PropertyListing) IsQuotaCounted() bool {
	return l.Status == ListingStatusDraft ||
		l.Status == ListingStatusPending ||
		l.Status == ListingStatusApproved
}

// CanBeEdited returns true if the listing can be edited by salesman
func (l *PropertyListing) CanBeEdited() bool {
	return l.Status == ListingStatusDraft || l.Status == ListingStatusRejected
}

// CanBeDeleted returns true if the listing can be soft-deleted by salesman
func (l *PropertyListing) CanBeDeleted() bool {
	return l.Status == ListingStatusDraft || l.Status == ListingStatusRejected
}

// CanBeSubmitted returns true if the listing can be submitted for review
func (l *PropertyListing) CanBeSubmitted() bool {
	return l.Status == ListingStatusDraft || l.Status == ListingStatusRejected
}

// ValidStatusTransitions maps allowed status transitions
var ValidStatusTransitions = map[string][]string{
	ListingStatusDraft:    {ListingStatusPending, ListingStatusDeleted},
	ListingStatusPending:  {ListingStatusApproved, ListingStatusRejected},
	ListingStatusRejected: {ListingStatusDraft, ListingStatusDeleted},
	ListingStatusApproved: {ListingStatusInactive, ListingStatusSold, ListingStatusRented},
	ListingStatusInactive: {ListingStatusApproved, ListingStatusDeleted},
	ListingStatusSold:     {},
	ListingStatusRented:   {},
	ListingStatusDeleted:  {},
}

// IsValidTransition checks if a status transition is allowed
func IsValidTransition(from, to string) bool {
	targets, ok := ValidStatusTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}
