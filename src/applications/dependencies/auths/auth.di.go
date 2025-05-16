package authdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	authusecase "kiraform/src/applications/usecases/auths"

	"gorm.io/gorm"
)

type AuthDependencies struct {
	DB *gorm.DB
	UC authusecase.AuthUsecase
}

func NewAuthDependencies(DB *gorm.DB) *AuthDependencies {
	// load necessary repositories
	// it possible to more than one
	userRepo := masterrepo.NewUserRepository(DB)
	roleRepo := masterrepo.NewRoleRepository(DB)

	// load the usecase and inject into Dependency
	authUC := authusecase.NewAuthUsecase(userRepo, roleRepo)
	return &AuthDependencies{
		DB: DB,
		UC: authUC,
	}
}
