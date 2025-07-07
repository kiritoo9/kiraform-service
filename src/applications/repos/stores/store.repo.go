package storerepo

import (
	"kiraform/src/applications/models"

	"gorm.io/gorm"
)

type StoreRepository interface {
	FindStoreByUser(userID string) (*models.Stores, error)
	CreateStore(data models.Stores) error
	CreateStoreUser(data models.StoreUsers) error
	UpdateStore(ID string, data models.Stores) error
}

type StoreQuery struct {
	DB *gorm.DB
}

func NewStoreRepository(DB *gorm.DB) *StoreQuery {
	return &StoreQuery{DB: DB}
}

func (q *StoreQuery) FindStoreByUser(userID string) (*models.Stores, error) {
	var store models.Stores
	if err := q.DB.Model(&models.Stores{}).
		Where("stores.deleted = ? AND store_users.deleted = ? AND store_users.user_id = ?", false, false, userID).
		Joins("LEFT JOIN store_users ON stores.id = store_users.store_id AND store_users.deleted = ?", false).
		First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (q *StoreQuery) CreateStore(data models.Stores) error {
	if err := q.DB.Model(&models.Stores{}).Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *StoreQuery) CreateStoreUser(data models.StoreUsers) error {
	if err := q.DB.Model(&models.StoreUsers{}).Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *StoreQuery) UpdateStore(ID string, data models.Stores) error {
	if err := q.DB.Model(&models.Stores{}).Where("id = ? AND deleted = ?", ID, false).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}
