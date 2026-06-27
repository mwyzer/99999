package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	InquiryStatusUnread  = "unread"
	InquiryStatusRead    = "read"
	InquiryStatusReplied = "replied"
	InquiryStatusClosed  = "closed"
)

type Inquiry struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PropertyID uuid.UUID  `gorm:"type:uuid;not null;index" json:"property_id"`
	BuyerID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"buyer_id"`
	Message    *string    `gorm:"type:text" json:"message"`
	Status     string     `gorm:"type:varchar(20);not null;default:unread" json:"status"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Property *PropertyListing `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Buyer    *User            `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}

func (Inquiry) TableName() string { return "inquiries" }
