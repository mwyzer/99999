package models

import "github.com/google/uuid"

type Location struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	City      string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_loc_city_prov" json:"city"`
	Province  string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_loc_city_prov" json:"province"`
	Country   string     `gorm:"type:varchar(50);not null;default:Indonesia" json:"country"`
	Latitude  *float64   `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude *float64   `gorm:"type:decimal(10,7)" json:"longitude"`
	IsActive  bool       `gorm:"not null;default:true" json:"is_active"`

	Listings []PropertyListing `gorm:"foreignKey:LocationID" json:"listings,omitempty"`
}

func (Location) TableName() string { return "locations" }
