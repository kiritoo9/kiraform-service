package storedi

import (
	storerepo "kiraform/src/applications/repos/stores"
	storeusecase "kiraform/src/applications/usecases/stores"

	"gorm.io/gorm"
)

type StoreDependencies struct {
	DB *gorm.DB
	UC storeusecase.StoreUsecase
}

func NewStoreDependencies(DB *gorm.DB) *StoreDependencies {
	storeRepo := storerepo.NewStoreRepository(DB)
	UC := storeusecase.NewStoreUsecase(storeRepo)
	return &StoreDependencies{
		DB: DB,
		UC: UC,
	}
}
