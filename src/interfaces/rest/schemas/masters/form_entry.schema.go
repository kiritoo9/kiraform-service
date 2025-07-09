package masterschema

import (
	"time"
)

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
	FormName                string    `json:"form_name"`
	FormCode                string    `json:"form_code"`
	Value                   string    `json:"value"`
	CreatedAt               time.Time `json:"created_at"`
}

type FormEntrySchema struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	UserName            string    `json:"user_name"`
	UserEmail           string    `json:"user_email"`
	CampaignID          string    `json:"campaign_id"`
	CampaignTitle       string    `json:"campaign_title"`
	CampaignDescription string    `json:"campaign_description"`
	Status              string    `json:"status"`
	Remark              string    `json:"remark"`
	CreatedAt           time.Time `json:"created_at"`
	ProductID           *string   `json:"product_id"`
}

type ProductResponse struct {
	ID                  string    `json:"id"`
	StoreID             string    `json:"store_id"`
	CategoryID          string    `json:"category_id"`
	CampaignID          *string   `json:"campaign_id"`
	Key                 string    `json:"key"`
	Slug                string    `json:"slug"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Price               int64     `json:"price"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	Thumbnail           string    `json:"thumbnail"`
	CategoryName        string    `json:"category_name"`
	CategoryDescription string    `json:"category_description"`
	CampaignTitle       string    `json:"campaign_title"`
	CampaignDescription string    `json:"campaign_description"`
}

type FormEntryResponse struct {
	Header  FormEntrySchema         `json:"header"`
	Detail  []FormDetailEntrySchema `json:"detail"`
	Product *ProductResponse        `json:"product"`
}
