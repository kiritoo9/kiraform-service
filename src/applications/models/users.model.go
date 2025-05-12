package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	UserIdentity string    `gorm:"type:varchar(100);unique;" json:"user_identity"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password     string    `gorm:"type:varchar(100);not null" json:"password"`
	Fullname     string    `gorm:"type:varchar(255);not null" json:"fullname"`
	Deleted      bool      `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt    time.Time `gorm:"type:timestamp;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;" json:"updated_at"`
}
