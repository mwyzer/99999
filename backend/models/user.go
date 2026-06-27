package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User roles
const (
	RoleBuyer         = "buyer"
	RoleSalesman      = "salesman"
	RoleTenantAdmin   = "tenant_admin"
	RolePlatformAdmin = "platform_admin"
)

// User status
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID     *uuid.UUID     `gorm:"type:uuid;index" json:"tenant_id"` // null for buyer & platform_admin
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"` // never exposed
	Name         string         `gorm:"type:varchar(200);not null" json:"name"`
	Phone        *string        `gorm:"type:varchar(20)" json:"phone"`
	PhotoURL     *string        `gorm:"type:varchar(500)" json:"photo_url"`
	Role         string         `gorm:"type:varchar(20);not null" json:"role"`
	Status       string         `gorm:"type:varchar(20);not null;default:active" json:"status"`
	WhatsappNumber *string      `gorm:"type:varchar(20)" json:"whatsapp_number"`
	ShowWhatsapp   *bool        `gorm:"type:boolean;default:true" json:"show_whatsapp"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Tenant          *Tenant            `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Listings        []PropertyListing  `gorm:"foreignKey:SalesmanID" json:"listings,omitempty"`
	ApprovedListings []PropertyListing `gorm:"foreignKey:ApprovedBy" json:"-"`
	SavedProperties []SavedProperty    `gorm:"foreignKey:BuyerID" json:"-"`
	AuditLogs       []AuditLog         `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}
