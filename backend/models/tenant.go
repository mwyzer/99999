package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant status
const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
)

type Tenant struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationName string         `gorm:"type:varchar(200);not null" json:"organization_name"`
	SubdomainSlug    string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"subdomain_slug"`
	LogoURL          *string        `gorm:"type:varchar(500)" json:"logo_url"`
	Description      *string        `gorm:"type:text" json:"description"`
	Phone            *string        `gorm:"type:varchar(20)" json:"phone"`
	Address          *string        `gorm:"type:text" json:"address"`
	Status           string         `gorm:"type:varchar(20);not null;default:active" json:"status"`
	WhatsappNumber   *string        `gorm:"type:varchar(20)" json:"whatsapp_number"`
	ShowWhatsapp     *bool          `gorm:"type:boolean;default:true" json:"show_whatsapp"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Users         []User         `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	Subscription  *Subscription  `gorm:"foreignKey:TenantID" json:"subscription,omitempty"`
	Listings      []PropertyListing `gorm:"foreignKey:TenantID" json:"listings,omitempty"`
}

func (Tenant) TableName() string {
	return "tenants"
}
