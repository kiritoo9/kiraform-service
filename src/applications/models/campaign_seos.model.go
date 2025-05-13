package models

import (
	"time"

	"github.com/google/uuid"
)

type CampaignSeos struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignID uuid.UUID  `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign   Campaigns  `gorm:"foreignKey:CampaignID;references:ID;constraint:OnDelete:CASCADE" json:"campaign"`
	Platform   string     `gorm:"type:varchar(50);not null" json:"platform"`
	Event      string     `gorm:"type:varchar(50);not null" json:"event"`
	AccessKey  string     `gorm:"type:text;not null" json:"access_key"`
	Deleted    bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt  time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
