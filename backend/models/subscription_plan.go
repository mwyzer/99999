package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SubStatusActive         = "active"
	SubStatusExpired        = "expired"
	SubStatusCancelled      = "cancelled"
	SubStatusPendingUpgrade = "pending_upgrade"
)

type SubscriptionPlan struct {
	ID                     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name                   string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug                   string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	MaxSalesmen            int       `gorm:"not null;default:5" json:"max_salesmen"`
	MaxListingsPerSalesman int       `gorm:"not null;default:5" json:"max_listings_per_salesman"`
	Description            *string   `gorm:"type:text" json:"description"`
	IsActive               bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt              time.Time `gorm:"autoCreateTime" json:"created_at"`

	Subscriptions []TenantSubscription `gorm:"foreignKey:PlanID" json:"subscriptions,omitempty"`
}

func (SubscriptionPlan) TableName() string { return "subscription_plans" }

type TenantSubscription struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"tenant_id"`
	PlanID    uuid.UUID  `gorm:"type:uuid;not null" json:"plan_id"`
	StartDate time.Time  `gorm:"not null;default:now()" json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Status    string     `gorm:"type:varchar(20);not null;default:active" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Tenant *Tenant           `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Plan   *SubscriptionPlan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

func (TenantSubscription) TableName() string { return "tenant_subscriptions" }
