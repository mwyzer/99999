package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Audit actions
const (
	AuditActionCreate   = "create"
	AuditActionUpdate   = "update"
	AuditActionDelete   = "delete"
	AuditActionApprove  = "approve"
	AuditActionReject   = "reject"
	AuditActionSuspend  = "suspend"
	AuditActionActivate = "activate"
)

// JSONMap for JSONB fields
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONMap: not []byte")
	}
	return json.Unmarshal(bytes, j)
}

type AuditLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	TenantID    *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id"`
	UserRole    string     `gorm:"type:varchar(20);not null" json:"user_role"`
	Action      string     `gorm:"type:varchar(50);not null" json:"action"`
	Module      string     `gorm:"type:varchar(50);not null;default:property" json:"module"`
	EntityType  string     `gorm:"type:varchar(50);not null;index" json:"entity_type"`
	EntityID    string     `gorm:"type:varchar(36);not null;index" json:"entity_id"`
	Description *string    `gorm:"type:text" json:"description"`
	OldValues   *JSONMap   `gorm:"type:jsonb" json:"old_values"`
	NewValues   *JSONMap   `gorm:"type:jsonb" json:"new_values"`
	IPAddress   *string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   *string    `gorm:"type:text" json:"user_agent"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;index" json:"created_at"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
