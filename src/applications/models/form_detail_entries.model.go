package models

import (
	"time"

	"github.com/google/uuid"
)

type FormDetailEntries struct {
	ID                      uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	FormEntryID             uuid.UUID     `gorm:"type:uuid;not null" json:"form_entry_id"`
	FormEntry               FormEntries   `gorm:"foreignKey:FormEntryID;references:ID;constraint:OnDelete:Cascade" json:"form_entry"`
	CampaignFormID          uuid.UUID     `gorm:"type:uuid;not null" json:"campaign_form_id"`
	CampaignForm            CampaignForms `gorm:"foreignKey:CampaignFormID;references:ID;constraint:OnDelete:CASCADE" json:"campaign_form"`
	CampaignFormAttributeID *uuid.UUID    `gorm:"type:uuid;null;comment:ID for form attribute, it can be null for form type that have no attributes" json:"campaign_form_attribute_id"`
	Value                   string        `gorm:"type:text" json:"value"`
	Deleted                 bool          `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt               time.Time     `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt               *time.Time    `gorm:"type:timestamp" json:"updated_at"`
}
