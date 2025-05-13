package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPackages struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User       Users      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	PackageID  uuid.UUID  `gorm:"type:uuid;not null" json:"package_id"`
	Package    Packages   `gorm:"foreignKey:PackageID;references:ID;constraint:OnDelete:CASCADE" json:"package"`
	ActiveDate time.Time  `gorm:"type:timestamp" json:"active_date"`
	ExpireDate time.Time  `gorm:"type:timestamp" json:"expire_date"`
	Remark     string     `gorm:"type:text" json:"remark"`
	IsActive   bool       `gorm:"type:boolean;default:false" json:"is_active"`
	Deleted    bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt  time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
