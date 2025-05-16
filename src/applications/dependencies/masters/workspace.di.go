package masterdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	masterusecase "kiraform/src/applications/usecases/masters"

	"gorm.io/gorm"
)

type WorkspaceDependencies struct {
	DB *gorm.DB
	UC masterusecase.WorkspaceUsecase
}

func NewWorkspaceDependencies(DB *gorm.DB) *WorkspaceDependencies {
	workspaceRepo := masterrepo.NewWorkspaceRepository(DB)
	UC := masterusecase.NewWorkspaceUsecase(workspaceRepo)
	return &WorkspaceDependencies{
		DB: DB,
		UC: UC,
	}
}
