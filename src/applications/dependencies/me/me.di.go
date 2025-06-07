package medi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	meusecase "kiraform/src/applications/usecases/me"

	"gorm.io/gorm"
)

type MeDependencies struct {
	DB *gorm.DB
	UC meusecase.MeUsecase
}

func NewMeDependencies(DB *gorm.DB) *MeDependencies {
	UC := meusecase.NewMeUsecase(masterrepo.NewUserRepository(DB))
	return &MeDependencies{
		DB: DB,
		UC: UC,
	}
}
