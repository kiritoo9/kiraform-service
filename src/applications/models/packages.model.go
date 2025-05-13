package models

import (
	"time"

	"github.com/google/uuid"
)

type Packages struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Deleted     bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
