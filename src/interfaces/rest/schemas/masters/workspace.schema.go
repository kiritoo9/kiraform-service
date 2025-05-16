package masterschema

type WorkspacePayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	IsPublish   bool   `json:"is_publish" default:"false"`
	Thumbnail   string `json:"thumbnail"`
}
