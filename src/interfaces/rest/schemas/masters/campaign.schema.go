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

type CampaignSchemaWithSummary struct {
	ID           uuid.UUID  `json:"id"`
	WorkspaceID  string     `json:"workspace_id"`
	Title        string     `json:"title"`
	Key          string     `json:"key"`
	Slug         string     `json:"slug"`
	Description  string     `json:"description"`
	IsPublish    bool       `json:"is_publish"`
	TotalVisitor int64      `json:"total_visitor"`
	TotalSubmit  int64      `json:"total_submit"`
	CreatedAt    *time.Time `json:"created_at"`
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

type CampaignDashboard struct {
	TotalVisitor int64 `json:"total_visitor"`
	TotalSubmit  int64 `json:"total_submit"`
}

type CampaignFormEntryChart struct {
	Total int64  `json:"total"`
	Date  string `json:"date"`
}

type FormEntryList struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	CampaignID    string `json:"campaign_id"`
	WorkspaceID   string `json:"workspace_id"`
	Status        string `json:"status"`
	Remark        string `json:"remark"`
	CampaignTitle string `json:"campaign_title"`
	CampaignKey   string `json:"campaign_key"`
	CampaignSlug  string `json:"campaign_slug"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	CreatedAt     string `json:"created_at"`
}
