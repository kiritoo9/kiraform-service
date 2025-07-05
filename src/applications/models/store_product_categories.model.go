package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreProductCategories struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID     uuid.UUID  `gorm:"type:uuid;not null" json:"store_id"`
	Store       Stores     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE" json:"store"`
	Name        string     `gorm:"type:varchar" json:"name"`
	Description string     `gorm:"type:varchar" json:"description"`
	Deleted     bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
