package masterdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	masterusecase "kiraform/src/applications/usecases/masters"

	"gorm.io/gorm"
)

type FormDependencies struct {
	DB *gorm.DB
	UC masterusecase.FormUsecase
}

func NewFormDependencies(DB *gorm.DB) *FormDependencies {
	formRepo := masterrepo.NewFormRepository(DB)
	UC := masterusecase.NewFormUsecase(formRepo)
	return &FormDependencies{
		DB: DB,
		UC: UC,
	}
}
