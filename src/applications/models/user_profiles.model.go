package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfiles struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User        Users      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	FirstName   string     `gorm:"type:varchar(100);not null" json:"first_name"`
	MiddleName  string     `gorm:"type:varchar(100)" json:"middle_name"`
	LastName    string     `gorm:"type:varchar(100)" json:"last_name"`
	Address     string     `gorm:"type:varchar(255)" json:"address"`
	Phone       string     `gorm:"type:varchar(20)" json:"phone"`
	Province    string     `gorm:"type:varchar(50)" json:"province"`
	City        string     `gorm:"type:varchar(70)" json:"city"`
	District    string     `gorm:"type:varchar(70)" json:"district"`
	SubDistrict string     `gorm:"type:varchar(70)" json:"sub_district"`
	Avatar      string     `gorm:"type:varchar(100)" json:"avatar"`
	Remark      string     `gorm:"type:varchar(255)" json:"remark"`
	Deleted     bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt   time.Time  `gorm:"type:timestamp;" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
