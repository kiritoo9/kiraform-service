package storeschema

import "time"

type ProductCategoryPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ProductCategoryResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	TotalProducts int64     `json:"total_products"`
}
