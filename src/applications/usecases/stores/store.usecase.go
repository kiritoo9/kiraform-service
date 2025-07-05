package storeusecase

import (
	"kiraform/src/applications/models"
	storerepo "kiraform/src/applications/repos/stores"
)

type StoreUsecase interface {
	FindStore() (*models.Stores, error)
}

type StoreService struct {
	storeRepo storerepo.StoreRepository
}

func NewStoreUsecase(storeRepo storerepo.StoreRepository) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
	}
}

func (s *StoreService) FindStore() (*models.Stores, error) {
	data, err := s.storeRepo.FindStore()
	if err != nil {
		return nil, err
	}
	return &data, nil
}
