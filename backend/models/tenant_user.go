package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TenantRoleAdmin    = "tenant_admin"
	TenantRoleSalesman = "salesman"
)

type TenantUser struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_tu_tenant_user" json:"tenant_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_tu_tenant_user" json:"user_id"`
	TenantRole string    `gorm:"type:varchar(20);not null" json:"tenant_role"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (TenantUser) TableName() string { return "tenant_users" }
