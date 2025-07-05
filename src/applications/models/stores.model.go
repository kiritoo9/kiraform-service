package models

import (
	"time"

	"github.com/google/uuid"
)

type Stores struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Key             string     `gorm:"type:string;unique" json:"key"`
	Name            string     `gorm:"type:varchar;not null" json:"name"`
	Slug            string     `gorm:"type:varchar;not null" json:"slug"`
	Category        string     `gorm:"type:varchar;not null" json:"category"`
	Description     string     `gorm:"type:text" json:"description"`
	Thumbnail       string     `gorm:"type:text" json:"thumbnail"`
	OperationalHour string     `gorm:"type:text" json:"operational_hour"`
	Adddress        string     `gorm:"type:text" json:"address"`
	Phone           string     `gorm:"type:varchar(20)" json:"phone"`
	Email           string     `gorm:"type:varchar(70)" json:"email"`
	Status          string     `gorm:"type:char(2);not null;comment:S1=PENDING,S2=ACTIVE,S3=INACTIVE" json:"status"`
	Deleted         bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt       time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
