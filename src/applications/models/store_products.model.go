package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreProducts struct {
	ID          uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID     uuid.UUID              `gorm:"type:uuid;not null" json:"store_id"`
	Store       Stores                 `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE" json:"store"`
	CategoryID  uuid.UUID              `gorm:"type:uuid;not null" json:"category_id"`
	Category    StoreProductCategories `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE" json:"category"`
	CampaignID  *uuid.UUID             `gorm:"type:uuid;null" json:"form_id"`
	Campaign    Campaigns              `gorm:"foreignKey:CampaignID;references:ID;constraint:OnDelete:CASCADE" json:"campaign"`
	Key         string                 `gorm:"type:varchar;not null;unique" json:"key"`
	Name        string                 `gorm:"type:varchar;not null" json:"name"`
	Slug        string                 `gorm:"type:varchar;not null" json:"slug"`
	Description string                 `gorm:"type:text" json:"description"`
	Price       int64                  `gorm:"type:numeric;default:0" json:"price"`
	Status      string                 `gorm:"type:char(2);default:S1;comment:S1=DRAFT;S2=PUBLISH;S3=OUT_OF_STOCK" json:"status"`
	Deleted     bool                   `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt   time.Time              `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   *time.Time             `gorm:"type:timestamp" json:"updated_at"`
}
