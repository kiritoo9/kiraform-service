package masterschema

import (
	"time"

	"github.com/google/uuid"
)

type CampaignFormPayload struct {
	ID           *string                         `json:"id"`
	FormID       string                          `json:"form_id" validate:"required"`
	Title        string                          `json:"title" validate:"required"`
	Description  string                          `json:"description"`
	Placeholder  string                          `json:"placeholder"`
	DefaultValue string                          `json:"default_value"`
	IsRequired   bool                            `json:"is_required" default:"false"`
	IsMultiple   bool                            `json:"is_multiple" default:"false"`
	Attributes   *[]CampaignFormAttributePayload `json:"attributes"`
}

type CampaignFormSchema struct {
	ID           uuid.UUID  `json:"id"`
	FormID       uuid.UUID  `json:"form_id"`
	FormCode     string     `json:"form_code"`
	FormName     string     `json:"form_name"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Placeholder  string     `json:"placeholder"`
	DefaultValue string     `json:"default_value"`
	IsRequired   bool       `json:"is_required"`
	IsMultiple   bool       `json:"is_multiple"`
	CreatedAt    *time.Time `json:"created_at"`
}

type DetailCampaignFormSchema struct {
	ID           uuid.UUID                      `json:"id"`
	FormID       uuid.UUID                      `json:"form_id"`
	FormCode     string                         `json:"form_code"`
	FormName     string                         `json:"form_name"`
	Title        string                         `json:"title"`
	Description  string                         `json:"description"`
	Placeholder  string                         `json:"placeholder"`
	DefaultValue string                         `json:"default_value"`
	IsRequired   bool                           `json:"is_required"`
	IsMultiple   bool                           `json:"is_multiple"`
	CreatedAt    *time.Time                     `json:"created_at"`
	Attributes   []CampaignFormAttributeSchemas `json:"attributes"`
}
