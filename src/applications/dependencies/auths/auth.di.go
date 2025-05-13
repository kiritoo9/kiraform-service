package authdi

import (
	userrepo "kiraform/src/applications/repos/masters"
	authusecase "kiraform/src/applications/usecases/auths"

	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
	UC authusecase.Usecase
}

func NewDependencies(DB *gorm.DB) *Dependencies {
	// load necessary repositories
	// it possible to more than one
	userRepo := userrepo.NewRepository(DB)

	// load the usecase and inject into Dependency
	authUC := authusecase.NewUsecase(userRepo)
	return &Dependencies{
		DB: DB,
		UC: authUC,
	}
}
