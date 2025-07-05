package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreProductImages struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	StoreProductID uuid.UUID     `gorm:"type:uuid;not null" json:"store_product_id"`
	StoreProduct   StoreProducts `gorm:"foreignKey:StoreProductID;references:ID;constarint:OnDelete:CASCADE" json:"store_product"`
	FileName       string        `gorm:"type:varchar" json:"file_name"`
	FileExt        string        `gorm:"type:varchar" json:"file_ext"`
	FileSize       string        `gorm:"type:varchar" json:"file_size"`
	FilePath       string        `gorm:"type:varchar" json:"file_path"`
	Deleted        bool          `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt      time.Time     `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt      *time.Time    `gorm:"type:timestamp" json:"updated_at"`
}
