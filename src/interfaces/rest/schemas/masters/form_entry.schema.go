package masterschema

import "time"

type FormEntryPayload struct {
	CampaignFormID          string  `json:"campaign_form_id" validate:"required"`
	CampaignFormAttributeID *string `json:"campaign_form_attribute_id" default:"null"`
	Value                   string  `json:"value"`
}

type FormDetailEntrySchema struct {
	ID                      string    `json:"id"`
	CampaignFormID          string    `json:"campaign_form_id"`
	CampaignFormAttributeID *string   `json:"campaign_form_attribute_id"`
	CampaignFormTitle       string    `json:"campaign_form_title"`
	CampaignFormDescription string    `json:"campaign_form_description"`
	CampaignFormName        string    `json:"campaign_form_name"`
	Value                   string    `json:"value"`
	CreatedAt               time.Time `json:"created_at"`
}

type FormEntrySchema struct {
	ID            string                  `json:"id"`
	UserID        string                  `json:"user_id"`
	UserName      string                  `json:"user_name"`
	UserEmail     string                  `json:"user_email"`
	CampaignID    string                  `json:"campaign_id"`
	CampaignTitle string                  `json:"campaign_title"`
	Status        string                  `json:"status"`
	Remark        string                  `json:"remark"`
	Detail        []FormDetailEntrySchema `json:"detail"`
	CreatedAt     time.Time               `json:"created_at"`
}
