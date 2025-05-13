package models

import (
	"time"

	"github.com/google/uuid"
)

type Campaigns struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID  `gorm:"type:uuid;not null"`
	Workspace   Workspaces `gorm:"foreignKey:WorkspaceID;references:ID;constraint:OnDelete:CASCADE" json:"workspace"`
	Key         string     `gorm:"type:varchar(100);not null;unique;comment:Generate by system" json:"key"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string     `gorm:"type:varchar(255);not null" json:"slug"`
	Description string     `gorm:"type:text" json:"description"`
	Thumbnail   string     `gorm:"type:varchar(100)" json:"thumbnail"`
	IsPublish   bool       `gorm:"type:bool;default:false" json:"is_publish"`
	Deleted     bool       `gorm:"type:bool;default:false" json:"deleted"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"type:timestamp" json:"updated_at"`
}
