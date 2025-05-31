package masterschema

import "time"

type CampaignSeoPayload struct {
	Platform  string `json:"platform" validate:"required"`
	Event     string `json:"event" validate:"required"`
	AccessKey string `json:"access_key" validate:"required"`
}

type CampaignSeoSchema struct {
	ID         string     `json:"id"`
	CampaignID string     `json:"campaign_id"`
	Platform   string     `json:"platform"`
	Event      string     `json:"event"`
	AccessKey  string     `json:"access_key"`
	CreatedAt  *time.Time `json:"created_at"`
}
