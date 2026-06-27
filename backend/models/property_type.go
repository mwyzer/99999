package models

import "github.com/google/uuid"

type PropertyType struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	Description *string   `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`

	Listings []PropertyListing `gorm:"foreignKey:PropertyTypeID" json:"listings,omitempty"`
}

func (PropertyType) TableName() string { return "property_types" }
