package models

import (
	"time"

	"github.com/google/uuid"
)

type FormEntries struct {
	ID         uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     *uuid.UUID    `gorm:"type:uuid;null" json:"user_id"`
	User       Users         `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	CampaignID uuid.UUID     `gorm:"type:uuid;not null" json:"campaign_id"`
	Campaign   Campaigns     `gorm:"foreignKey:CampaignID;references:ID;constraint:OnDelete:CASCADE" json:"campaign"`
	ProductID  *uuid.UUID    `gorm:"type:uuid;not null" json:"product_id"`
	Product    StoreProducts `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE" json:"product"`
	Status     string        `gorm:"type:char(2);default:S1;comment:S1=PENDING,S2=APPROVED;S3=REJECTED" json:"status"`
	Remark     string        `gorm:"type:text" json:"remark"`
	Deleted    bool          `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt  time.Time     `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt  *time.Time    `gorm:"type:timestamp" json:"updated_at"`
}
