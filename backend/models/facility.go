package models

import (
	"github.com/google/uuid"
)

type Facility struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Icon     *string   `gorm:"type:varchar(50)" json:"icon"`
	IsActive bool      `gorm:"not null;default:true" json:"is_active"`

	PropertyFacilities []PropertyFacility `gorm:"foreignKey:FacilityID" json:"-"`
}

func (Facility) TableName() string { return "facilities" }

type PropertyFacility struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PropertyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_pf_prop_fac" json:"property_id"`
	FacilityID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_pf_prop_fac" json:"facility_id"`
	Value      *string   `gorm:"type:text" json:"value"`

	Property *PropertyListing `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Facility *Facility        `gorm:"foreignKey:FacilityID" json:"facility,omitempty"`
}

func (PropertyFacility) TableName() string { return "property_facilities" }
