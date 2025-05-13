package models

import (
	"time"

	"github.com/google/uuid"
)

type CampaignForms struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignID   uuid.UUID  `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign     Campaigns  `gorm:"foreignKey:CampaignID;references:ID;constraint:OnDelete:CASCADE" json:"campaign"`
	FormID       uuid.UUID  `gorm:"type:uuid;not null" json:"form_id"`
	Form         Forms      `gorm:"foreignKey:FormID;references:ID;constraint:OnDelete:CASCADE" json:"form"`
	Title        string     `gorm:"type:varchar(100);not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Placeholder  string     `gorm:"type:varchar(150)" json:"placeholder"`
	DefaultValue string     `gorm:"type:varchar(150)" json:"default_value"`
	IsRequired   bool       `gorm:"type:boolean;default:false" json:"is_required"`
	IsMultiple   bool       `gorm:"type:boolean;default:false" json:"is_multiple"`
	Deleted      bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt    time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
