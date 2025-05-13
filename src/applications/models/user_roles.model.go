package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRoles struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User      Users      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	RoleID    uuid.UUID  `gorm:"type:uuid;not null" json:"role_id"`
	Role      Roles      `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE" json:"role"`
	Deleted   bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt time.Time  `gorm:"type:timestamp;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
