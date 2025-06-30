package masterschema

import "time"

type WorkspacePayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	IsPublish   bool   `json:"is_publish" default:"false"`
	Thumbnail   string `json:"thumbnail"`
}

type WorkspaceList struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Key         string    `json:"key"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	IsPublish   bool      `json:"is_publish"`
	Thumbnail   string    `json:"thumbnail"`
	TotalForm   int64     `json:"total_form"`
	TotalSubmit int64     `json:"total_submit"`
	CreatedAt   time.Time `json:"created_at"`
}
