package models

import (
	"time"

	"github.com/google/uuid"
)

// Plan types
const (
	PlanFree    = "free"
	PlanPremium = "premium"
)

type Subscription struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID               uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"tenant_id"`
	PlanType               string     `gorm:"type:varchar(20);not null;default:free" json:"plan_type"`
	MaxSalesmen            int        `gorm:"not null;default:5" json:"max_salesmen"`
	MaxListingsPerSalesman int        `gorm:"not null;default:5" json:"max_listings_per_salesman"`
	StartDate              time.Time  `gorm:"autoCreateTime" json:"start_date"`
	EndDate                *time.Time `json:"end_date"`
	CreatedAt              time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
