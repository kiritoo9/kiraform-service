package storerepo

import (
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"strings"

	"gorm.io/gorm"
)

type StoreRepository interface {
	FindStoreByUser(userID string) (*models.Stores, error)
	CreateStore(data models.Stores) error
	CreateStoreUser(data models.StoreUsers) error
	UpdateStore(ID string, data models.Stores) error
	FindStoreProductCategories(storeID string, paramns *commonschema.QueryParams) ([]models.StoreProductCategories, error)
	FindCountStoreProductCategories(storeID string, params *commonschema.QueryParams) (int64, error)
	FindStoreProductCategory(storeID string, ID string) (*models.StoreProductCategories, error)
	CreateStoreProductCategory(data models.StoreProductCategories) error
	UpdateStoreProductCategory(userID string, data models.StoreProductCategories) error
	FindStoreProducts(storeID string, params *commonschema.QueryParams) ([]models.StoreProducts, error)
	FindCountStoreProducts(storeID string, params *commonschema.QueryParams) (int64, error)
	CreateProduct(data models.StoreProducts) error
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

func (q *StoreQuery) FindStoreProductCategories(storeID string, params *commonschema.QueryParams) ([]models.StoreProductCategories, error) {
	var data []models.StoreProductCategories

	// init statement
	st := q.DB.Model(&models.StoreProductCategories{}).Where("deleted = ? AND store_id = ?", false, storeID)

	// handle search condition
	if params.Search != "" {
		st = st.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// handle pagination
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = (params.Limit * params.Page) - params.Limit
	}
	st = st.Order("created_at DESC")
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get data
	if err := st.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (q *StoreQuery) FindCountStoreProductCategories(storeID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// init statement
	st := q.DB.Model(&models.StoreProductCategories{}).Where("deleted = ? AND store_id = ?", false, storeID)

	// handle search condition
	if params.Search != "" {
		st = st.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *StoreQuery) FindStoreProductCategory(storeID string, ID string) (*models.StoreProductCategories, error) {
	var data models.StoreProductCategories
	if err := q.DB.Model(&models.StoreProductCategories{}).Where("store_id = ? AND id = ? AND deleted = ?", storeID, ID, false).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (q *StoreQuery) CreateStoreProductCategory(data models.StoreProductCategories) error {
	if err := q.DB.Model(&models.StoreProductCategories{}).Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *StoreQuery) UpdateStoreProductCategory(ID string, data models.StoreProductCategories) error {
	if err := q.DB.Model(&models.StoreProductCategories{}).Where("id = ?", ID).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *StoreQuery) FindStoreProducts(storeID string, params *commonschema.QueryParams) ([]models.StoreProducts, error) {
	var data []models.StoreProducts

	// init statement
	st := q.DB.Model(&models.StoreProducts{}).Where("store_products.deleted = ? AND store_products.store_id = ?", false, storeID).
		Preload("Store").
		Preload("Category").
		Preload("Campaign")

	// handle search condition
	if params.Search != "" {
		st = st.Where("LOWER(store_products.name) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// handle pagination
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = (params.Limit * params.Page) - params.Limit
	}
	st = st.Order("store_products.created_at DESC")
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get data
	if err := st.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (q *StoreQuery) FindCountStoreProducts(storeID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// init statement
	st := q.DB.Model(&models.StoreProducts{}).Where("deleted = ? AND store_id = ?", false, storeID)

	// handle search condition
	if params.Search != "" {
		st = st.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *StoreQuery) CreateProduct(data models.StoreProducts) error {
	if err := q.DB.Model(&models.StoreProducts{}).Create(&data).Error; err != nil {
		return err
	}
	return nil
}
