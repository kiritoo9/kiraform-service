package models

import (
	"time"

	"github.com/google/uuid"
)

type CampaignFormEntries struct {
	ID                      uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID                  uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	User                    Users         `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	CampaignID              uuid.UUID     `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign                Campaigns     `gorm:"foreignKey:CampaignID;references:ID;constraint:OnDelete:CASCADE" json:"campaign"`
	CampaignFormID          uuid.UUID     `gorm:"type:uuid;not null" json:"campaign_form_id"`
	CampaignForm            CampaignForms `gorm:"foreignKey:CampaignFormID;references:ID;constraint:OnDelete:CASCADE" json:"campaign_form"`
	CampaignFormAttributeID uuid.UUID     `gorm:"type:uuid;null;comment:ID for form attribute, it can be null for form type that have no attributes" json:"campaign_form_attribute_id"`
	Value                   string        `gorm:"type:text" json:"value"`
	Remark                  string        `gorm:"type:text" json:"remark"`
	Deleted                 bool          `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt               time.Time     `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt               *time.Time    `gorm:"type:timestamp" json:"updated_at"`
}
