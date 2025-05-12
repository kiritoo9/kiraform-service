package models

import (
	"time"

	"github.com/google/uuid"
)

type Roles struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);unique;not null" json:"name"`
	Description string    `gorm:"type:varchar(100)" json:"description"`
	Deleted     bool      `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt   time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp" json:"updated_at"`
}
