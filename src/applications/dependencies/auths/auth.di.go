package authdi

import (
	repomasters "kiraform/src/applications/repos/masters"
	authusecase "kiraform/src/applications/usecases/auths"

	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
	UC authusecase.AuthUsecase
}

func NewDependencies(DB *gorm.DB) *Dependencies {
	// load necessary repositories
	// it possible to more than one
	userRepo := repomasters.NewUserRepository(DB)
	roleRepo := repomasters.NewRoleRepository(DB)

	// load the usecase and inject into Dependency
	authUC := authusecase.NewAuthUsecase(userRepo, roleRepo)
	return &Dependencies{
		DB: DB,
		UC: authUC,
	}
}
