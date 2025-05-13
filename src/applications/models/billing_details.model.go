package models

import (
	"time"

	"github.com/google/uuid"
)

type BillingDetails struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	BillingID uuid.UUID  `gorm:"type:uuid;not null" json:"billing_id"`
	Billing   Billings   `gorm:"foreignKey:BillingID;references:ID;constraint:OnDelete:CASCADE" json:"billing"`
	Item      string     `gorm:"type:varchar(120);not null" json:"item"`
	Qty       int        `gorm:"type:numeric;default:0" json:"qty"`
	Total     int        `gorm:"type:numeric;default:0" json:"total"`
	Remark    string     `gorm:"type:text;" json:"remark"`
	Deleted   bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
