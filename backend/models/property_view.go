package models

import (
	"time"

	"github.com/google/uuid"
)

type PropertyView struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PropertyID uuid.UUID `gorm:"type:uuid;not null;index" json:"property_id"`
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	IPAddress  *string   `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  *string   `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	Property *PropertyListing `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	User     *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PropertyView) TableName() string { return "property_views" }
