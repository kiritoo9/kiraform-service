package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceUsers struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID  `gorm:"type:uuid;not null" json:"workspace_id"`
	Workspace   Workspaces `gorm:"foreignKey:WorkspaceID;references:ID;constraint:OnDelete:CASCADE" json:"workspace"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User        Users      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
	Status      string     `gorm:"type:char(2);default:S1;comment:S1=INVITED,S2=REQUESTED,S3=APPROVED,S4=REJECTED,S5=OWNER" json:"status"`
	Deleted     bool       `gorm:"type:boolean;default:false" json:"deleted"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;default:null" json:"updated_at"`
}
