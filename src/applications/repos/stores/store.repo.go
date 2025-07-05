package storerepo

import (
	"kiraform/src/applications/models"

	"gorm.io/gorm"
)

type StoreRepository interface {
	FindStore() (models.Stores, error)
}

type StoreQuery struct {
	DB *gorm.DB
}

func NewStoreRepository(DB *gorm.DB) *StoreQuery {
	return &StoreQuery{DB: DB}
}

func (q *StoreQuery) FindStore() (models.Stores, error) {
	return models.Stores{}, nil
}
