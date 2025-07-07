package storeschema

import (
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"time"
)

type ProductImages struct {
	ID       *string `json:"id"`
	FileName string  `json:"file_name"`
	FileExt  string  `json:"file_ext"`
	FileSize string  `json:"file_size"`
	FilePath string  `json:"file_path"`
}

type ProductResponse struct {
	ID          string                       `json:"id"`
	StoreID     string                       `json:"store_id"`
	CategoryID  string                       `json:"category_id"`
	CampaignID  *string                      `json:"campaign_id"`
	Key         string                       `json:"key"`
	Slug        string                       `json:"slug"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Price       int64                        `json:"price"`
	Status      string                       `json:"status"`
	CreatedAt   time.Time                    `json:"created_at"`
	Category    ProductCategoryResponse      `json:"category"`
	Campaign    *masterschema.CampaignSchema `json:"campaign"`
	Images      []ProductImages              `json:"images"`
}

type ProductPayload struct {
	CategoryID  string   `json:"category_id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Status      string   `json:"status" validate:"required"`
	CampaignID  *string  `json:"campaign_id"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	Images      []string `json:"images"`
}
