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
	// load repositories
	workspaceRepo := masterrepo.NewWorkspaceRepository(DB)
	userRepo := masterrepo.NewUserRepository(DB)

	// init dependencies
	UC := masterusecase.NewWorkspaceUsecase(workspaceRepo, userRepo)
	return &WorkspaceDependencies{
		DB: DB,
		UC: UC,
	}
}
