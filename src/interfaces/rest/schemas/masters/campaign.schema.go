package masterschema

import (
	"time"

	"github.com/google/uuid"
)

type CampaignPayload struct {
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description"`
	IsPublish   bool                  `json:"is_publish" default:"false"`
	Forms       []CampaignFormPayload `json:"forms" validate:"required"`
}

type CampaignSchema struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	Title       string     `json:"title"`
	Key         string     `json:"key"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	IsPublish   bool       `json:"is_publish"`
	CreatedAt   *time.Time `json:"created_at"`
}

type DetailCampaignSchema struct {
	ID          uuid.UUID                  `json:"id"`
	WorkspaceID string                     `json:"workspace_id"`
	Title       string                     `json:"title"`
	Key         string                     `json:"key"`
	Slug        string                     `json:"slug"`
	Description string                     `json:"description"`
	IsPublish   bool                       `json:"is_publish"`
	CreatedAt   *time.Time                 `json:"created_at"`
	Forms       []DetailCampaignFormSchema `json:"forms"`
}
