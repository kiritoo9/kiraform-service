package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreUsers struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID   uuid.UUID  `gorm:"type:uuid;not null" json:"store_id"`
	Store     Stores     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE" json:"store"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User      Users      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	Remark    string     `gorm:"type:text;" json:"remark"`
	Deleted   bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
