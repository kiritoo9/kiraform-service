package storeusecase

import (
	"errors"
	"fmt"
	"kiraform/src/applications/models"
	storerepo "kiraform/src/applications/repos/stores"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	storeschema "kiraform/src/interfaces/rest/schemas/stores"
	"kiraform/src/utils"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type StoreUsecase interface {
	FindStore(userID string) (*storeschema.StoreResponse, error)
	UpdateStore(userID string, body storeschema.StorePayload) error
	FindStoreProductCategories(userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindStoreProductCategory(userID string, ID string) (*storeschema.ProductCategoryResponse, error)
	CreateStoreProductCategory(userID string, body storeschema.ProductCategoryPayload) error
	UpdateStoreProductCategory(userID string, ID string, body storeschema.ProductCategoryPayload) error
	DeleteStoreProductCategory(userID string, ID string) error
}

type StoreService struct {
	storeRepo storerepo.StoreRepository
}

func NewStoreUsecase(storeRepo storerepo.StoreRepository) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
	}
}

func (s *StoreService) FindStore(userID string) (*storeschema.StoreResponse, error) {
	data, err := s.storeRepo.FindStoreByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store data not found, please create your store first")
		}
		return nil, err
	}
	return &storeschema.StoreResponse{
		ID:              data.ID.String(),
		Key:             data.Key,
		Name:            data.Name,
		Slug:            data.Slug,
		Category:        data.Category,
		Description:     data.Description,
		Phone:           data.Phone,
		Email:           data.Email,
		Address:         data.Address,
		OperationalHour: data.OperationalHour,
		Thumbnail:       data.Thumbnail,
		UpdatedAt:       data.UpdatedAt,
	}, nil
}

func (s *StoreService) UpdateStore(userID string, body storeschema.StorePayload) error {
	// check existing store by id
	// if exists then update it
	// otherwise insert new one with store_users

	exists, err := s.FindStore(userID)
	isExists := true
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isExists = false
		} else {
			return err
		}
	}

	// preparing data from payload
	store := models.Stores{
		Name:            body.Name,
		Slug:            fmt.Sprintf("ST-%s", slug.Make(body.Name)),
		Category:        body.Category,
		Description:     body.Description,
		Phone:           body.Phone,
		Email:           body.Email,
		Address:         body.Address,
		OperationalHour: body.OperationalHour,
	}

	// uploading image for store thumbnail
	if body.Thumbnail != nil {
		thumbnail, err := utils.UploadImage(*body.Thumbnail, "stores", store.Slug)
		if err != nil {
			return err
		}
		store.Thumbnail = *thumbnail

		// because user update the thumbnail
		// then remove last thumbnail in cdn/stores/{file_name} to make folder clean
		err = utils.RemoveImage(exists.Thumbnail)
		if err != nil {
			return err
		}
	}

	// perform to insert/update data
	if isExists {
		now := time.Now()

		// updating store data
		store.UpdatedAt = &now
		if err := s.storeRepo.UpdateStore(exists.ID, store); err != nil {
			return err
		}

		// updating store user data
	} else {
		store.ID = uuid.New()
		uuidArr := strings.Split(store.ID.String(), "-")
		if len(uuidArr) > 0 {
			store.Key = fmt.Sprintf("ST-%s", uuidArr[0])
		}
		store.CreatedAt = time.Now()
		store.Status = "S2" // force to Active for this version

		// preparing store user data
		uuidUserID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		storeUser := models.StoreUsers{
			ID:        uuid.New(),
			StoreID:   store.ID,
			UserID:    uuidUserID,
			CreatedAt: time.Now(),
		}

		// inserting store data
		if err := s.storeRepo.CreateStore(store); err != nil {
			return err
		}

		// inserting store user data
		if err := s.storeRepo.CreateStoreUser(storeUser); err != nil {
			return err
		}
	}

	// set as success response
	// by returning no error
	return nil
}

func (s *StoreService) FindStoreProductCategories(userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	// check valid store
	store, err := s.FindStore(userID)
	if err != nil {
		return nil, err
	}

	// perform to get product categories
	list, err := s.storeRepo.FindStoreProductCategories(store.ID, params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.storeRepo.FindCountStoreProductCategories(store.ID, params)
	if err != nil {
		return nil, err
	}

	// convert to response schema
	var data []storeschema.ProductCategoryResponse
	for _, v := range list {
		data = append(data, storeschema.ProductCategoryResponse{
			ID:          v.ID.String(),
			Name:        v.Name,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
		})
	}

	// prepare response list
	totalPage := 1
	if count > 0 && params.Limit > 0 {
		totalPage = int(math.Ceil(float64(count) / float64(params.Limit)))
	}

	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  totalPage,
		Rows:       data,
	}

	// return success response
	return &response, nil
}

func (s *StoreService) FindStoreProductCategory(userID string, ID string) (*storeschema.ProductCategoryResponse, error) {
	// check valid store
	store, err := s.FindStore(userID)
	if err != nil {
		return nil, err
	}

	// perform to get detail data
	data, err := s.storeRepo.FindStoreProductCategory(store.ID, ID)
	if err != nil {
		return nil, err
	}

	// return success response
	return &storeschema.ProductCategoryResponse{
		ID:          data.ID.String(),
		Name:        data.Name,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
	}, nil
}

func (s *StoreService) CreateStoreProductCategory(userID string, body storeschema.ProductCategoryPayload) error {
	// check valid store
	store, err := s.FindStore(userID)
	if err != nil {
		return err
	}

	// convert store_id to uuid format
	uuidStoreID, err := uuid.Parse(store.ID)
	if err != nil {
		return err
	}

	// perform to create data
	data := models.StoreProductCategories{
		ID:          uuid.New(),
		StoreID:     uuidStoreID,
		Name:        body.Name,
		Description: body.Description,
		CreatedAt:   time.Now(),
	}
	err = s.storeRepo.CreateStoreProductCategory(data)
	if err != nil {
		return err
	}

	// return success response
	// by flag as no-error
	return nil
}

func (s *StoreService) UpdateStoreProductCategory(userID string, ID string, body storeschema.ProductCategoryPayload) error {
	// check valid store
	_, err := s.FindStore(userID)
	if err != nil {
		return err
	}

	// perform to update data
	now := time.Now()
	data := models.StoreProductCategories{
		Name:        body.Name,
		Description: body.Description,
		UpdatedAt:   &now,
	}
	err = s.storeRepo.UpdateStoreProductCategory(ID, data)
	if err != nil {
		return err
	}

	// return success response
	// by flag as no-error
	return nil
}

func (s *StoreService) DeleteStoreProductCategory(userID string, ID string) error {
	_, err := s.FindStore(userID)
	if err != nil {
		return err
	}

	// perform to update data
	now := time.Now()
	data := models.StoreProductCategories{
		Deleted:   true,
		UpdatedAt: &now,
	}
	err = s.storeRepo.UpdateStoreProductCategory(ID, data)
	if err != nil {
		return err
	}

	// return success response
	// by flag as no-error
	return nil
}
