package models

import (
	"time"

	"github.com/google/uuid"
)

type CampaignFormAttributes struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignFormID uuid.UUID     `gorm:"type:uuid;not null" json:"campaign_form_id"`
	CampaignForm   CampaignForms `gorm:"foreignKey:CampaignFormID;references:ID;constraint:OnDelete:CASCADE" json:"campaign_form"`
	Label          string        `gorm:"type:varchar(50);not null" json:"label"`
	Value          string        `gorm:"type:varchar(50);not null" json:"value"`
	IsDefault      bool          `gorm:"type:boolean;default:false" json:"is_default"`
	Deleted        bool          `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt      time.Time     `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt      *time.Time    `gorm:"type:timestamp" json:"updated_at"`
}
