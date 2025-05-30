package masterschema

import (
	"time"

	"github.com/google/uuid"
)

type CampaignFormAttributePayload struct {
	ID        *string `json:"id"`
	Label     string  `json:"label" validate:"required"`
	Value     string  `json:"value" validate:"required"`
	IsDefault bool    `json:"is_default" default:"false"`
}

type CampaignFormAttributeSchemas struct {
	ID        uuid.UUID  `json:"id"`
	Label     string     `json:"label"`
	Value     string     `json:"value"`
	IsDefault bool       `json:"is_default"`
	CreatedAt *time.Time `json:"created_at"`
}
