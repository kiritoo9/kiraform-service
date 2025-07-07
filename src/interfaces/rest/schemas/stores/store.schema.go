package storeschema

import "time"

type StorePayload struct {
	Name            string  `json:"name" validate:"required"`
	Category        string  `json:"category" validate:"required"`
	Description     string  `json:"description"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	Address         string  `json:"address"`
	OperationalHour string  `json:"operational_hour"`
	Thumbnail       *string `json:"thumbnail"`
}

type StoreResponse struct {
	ID              string     `json:"id"`
	Key             string     `json:"key"`
	Name            string     `json:"name"`
	Slug            string     `json:"slug"`
	Category        string     `json:"category"`
	Description     string     `json:"description"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Address         string     `json:"address"`
	OperationalHour string     `json:"operational_hour"`
	Thumbnail       string     `json:"thumbnail"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
