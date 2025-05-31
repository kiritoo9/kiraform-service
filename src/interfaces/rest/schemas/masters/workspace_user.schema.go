package masterschema

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceUserPayload struct {
	UserID    *string `json:"user_id"`
	UserEmail *string `json:"user_email"`
	Status    string  `json:"status" validate:"required,max=2" default:"S1"`
}

type WorkspaceUserUpdatePayload struct {
	Status string `json:"status" validate:"required,max=2" default:"S1"`
}

type WorkspaceUserSchema struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	WorkspaceID    uuid.UUID `json:"workspace_id"`
	WorkspaceTitle string    `json:"workspace_title"`
	UserName       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
