package storeusecase

import (
	"errors"
	"fmt"
	"kiraform/src/applications/models"
	storerepo "kiraform/src/applications/repos/stores"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	storeschema "kiraform/src/interfaces/rest/schemas/stores"
	"kiraform/src/utils"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StoreUsecase interface {
	FindStore(c echo.Context, userID string) (*storeschema.StoreResponse, error)
	UpdateStore(userID string, body storeschema.StorePayload) error
	FindStoreProductCategories(userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindStoreProductCategory(userID string, ID string) (*storeschema.ProductCategoryResponse, error)
	CreateStoreProductCategory(userID string, body storeschema.ProductCategoryPayload) error
	UpdateStoreProductCategory(userID string, ID string, body storeschema.ProductCategoryPayload) error
	DeleteStoreProductCategory(userID string, ID string) error
	FindStoreProducts(c echo.Context, userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindStoreProductsByStoreKey(c echo.Context, key string, params *commonschema.QueryParams, category_id *string) (*commonschema.ResponseList, error)
	FindStoreProduct(c echo.Context, userID string, ID string) (*storeschema.ProductResponse, error)
	CreateStoreProduct(userID string, body storeschema.ProductPayload) error
	UpdateStoreProduct(userID string, ID string, body storeschema.ProductPayload) error
	DeleteStoreProduct(userID string, ID string) error
	FindStoreProductFormEntries(userID string, ID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindStoreByKey(c echo.Context, key string) (*storeschema.StoreResponse, error)
	FindStoreCategoriesByKey(key string) ([]storeschema.ProductCategoryResponse, error)
}

type StoreService struct {
	storeRepo storerepo.StoreRepository
}

func NewStoreUsecase(storeRepo storerepo.StoreRepository) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
	}
}

func (s *StoreService) findStore(userID string, returnOriginalError bool) (*storeschema.StoreResponse, error) {
	data, err := s.storeRepo.FindStoreByUser(userID)
	if err != nil {
		if returnOriginalError {
			return nil, err
		} else {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("store data not found, please create your store first")
			}
			return nil, err
		}
	}

	totalCategories, err := s.storeRepo.FindCountStoreProductCategories(data.ID.String(), nil)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.storeRepo.FindCountStoreProducts(data.ID.String(), nil, nil)
	if err != nil {
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
		TotalCategories: totalCategories,
		TotalProducts:   totalProducts,
	}, nil
}

func (s *StoreService) FindStore(c echo.Context, userID string) (*storeschema.StoreResponse, error) {
	data, err := s.findStore(userID, false)
	if err != nil {
		return nil, err
	}

	// generate image url
	data.Thumbnail = utils.ServeImage(c, data.Thumbnail)
	return data, nil
}

func (s *StoreService) FindStoreByKey(c echo.Context, key string) (*storeschema.StoreResponse, error) {
	data, err := s.storeRepo.FindStoreByKey(key)
	if err != nil {
		return nil, err
	}

	store := storeschema.StoreResponse{
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
		TotalCategories: 0,
		TotalProducts:   0,
	}

	// generate image url
	store.Thumbnail = utils.ServeImage(c, data.Thumbnail)
	return &store, nil
}

func (s *StoreService) UpdateStore(userID string, body storeschema.StorePayload) error {
	// check existing store by id
	// if exists then update it
	// otherwise insert new one with store_users

	exists, err := s.findStore(userID, true)
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

		if isExists && exists.Thumbnail != "" {
			// because user update the thumbnail
			// then remove last thumbnail in cdn/stores/{file_name} to make folder clean
			_ = utils.RemoveImage(exists.Thumbnail)
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
	store, err := s.findStore(userID, false)
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
		totalProducts, err := s.storeRepo.FindCountStoreProductsByCategory(store.ID, v.ID.String())
		if err != nil {
			return nil, err
		}

		data = append(data, storeschema.ProductCategoryResponse{
			ID:            v.ID.String(),
			Name:          v.Name,
			Description:   v.Description,
			CreatedAt:     v.CreatedAt,
			TotalProducts: totalProducts,
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
	store, err := s.findStore(userID, false)
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
	store, err := s.findStore(userID, false)
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
	_, err := s.findStore(userID, false)
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
	_, err := s.findStore(userID, false)
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

func (s *StoreService) findStoreProducts(c echo.Context, storeID string, params *commonschema.QueryParams, category_id *string) (*commonschema.ResponseList, error) {
	// perform to get product
	list, err := s.storeRepo.FindStoreProducts(storeID, params, category_id)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.storeRepo.FindCountStoreProducts(storeID, params, category_id)
	if err != nil {
		return nil, err
	}

	// convert to response schema
	data := []storeschema.ProductResponse{}
	for _, v := range list {
		var campaignID string
		if v.CampaignID != nil {
			campaignID = v.CampaignID.String()
		}

		// generate url for latest image of product
		thumbnail := ""
		images, err := s.storeRepo.FindImagesByProduct(v.ID.String())
		if err != nil {
			return nil, err
		} else if len(images) > 0 {
			thumbnail = images[0].FileName
		}

		if thumbnail != "" {
			thumbnail = utils.ServeImage(c, thumbnail)
		}

		productImages := []storeschema.ProductImages{}
		for _, j := range images {
			productImageID := j.ID.String()
			productImages = append(productImages, storeschema.ProductImages{
				ID:       &productImageID,
				FileName: utils.ServeImage(c, j.FileName),
			})
		}

		d := storeschema.ProductResponse{
			ID:          v.ID.String(),
			StoreID:     v.StoreID.String(),
			CategoryID:  v.CategoryID.String(),
			CampaignID:  &campaignID,
			Key:         v.Key,
			Slug:        v.Slug,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt,
			Thumbnail:   thumbnail,
			Images:      productImages,
			Category: storeschema.ProductCategoryResponse{
				ID:          v.CategoryID.String(),
				Name:        v.Category.Name,
				Description: v.Category.Description,
				CreatedAt:   v.Category.CreatedAt,
			},
		}

		if v.CampaignID != nil {
			d.Campaign = &masterschema.CampaignSchema{
				ID:          v.Campaign.ID,
				WorkspaceID: v.Campaign.WorkspaceID.String(),
				Title:       v.Campaign.Title,
				Key:         v.Campaign.Key,
				Slug:        v.Campaign.Slug,
				Description: v.Campaign.Description,
				IsPublish:   v.Campaign.IsPublish,
				CreatedAt:   &v.Campaign.CreatedAt,
			}
		}

		data = append(data, d)
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

func (s *StoreService) FindStoreProducts(c echo.Context, userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	// check valid store
	store, err := s.findStore(userID, false)
	if err != nil {
		return nil, err
	}

	list, err := s.findStoreProducts(c, store.ID, params, nil)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *StoreService) FindStoreProductsByStoreKey(c echo.Context, key string, params *commonschema.QueryParams, category_id *string) (*commonschema.ResponseList, error) {
	// check valid store
	store, err := s.storeRepo.FindStoreByKey(key)
	if err != nil {
		return nil, err
	}

	list, err := s.findStoreProducts(c, store.ID.String(), params, category_id)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *StoreService) FindStoreProduct(c echo.Context, userID string, ID string) (*storeschema.ProductResponse, error) {
	// check valid store
	store, err := s.findStore(userID, false)
	if err != nil {
		return nil, err
	}

	// check existing product
	data, err := s.storeRepo.FindStoreProduct(store.ID, ID)
	if err != nil {
		return nil, err
	}

	// generate url for thumbnail and entire images
	thumbnail := ""
	images, err := s.storeRepo.FindImagesByProduct(data.ID.String())
	if err != nil {
		return nil, err
	} else if len(images) > 0 {
		thumbnail = images[0].FileName
	}

	if thumbnail != "" {
		thumbnail = utils.ServeImage(c, thumbnail)
	}

	productImages := []storeschema.ProductImages{}
	for _, v := range images {
		imageID := v.ID.String()
		productImages = append(productImages, storeschema.ProductImages{
			ID:       &imageID,
			FileName: utils.ServeImage(c, v.FileName),
		})
	}

	// converting campaign ID
	var campaignID string
	if data.CampaignID != nil {
		campaignID = data.CampaignID.String()
	}

	// prepare for response
	product := storeschema.ProductResponse{
		ID:          data.ID.String(),
		StoreID:     data.StoreID.String(),
		CategoryID:  data.CategoryID.String(),
		CampaignID:  &campaignID,
		Key:         data.Key,
		Slug:        data.Slug,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		Thumbnail:   thumbnail,
		Category: storeschema.ProductCategoryResponse{
			ID:          data.CategoryID.String(),
			Name:        data.Category.Name,
			Description: data.Category.Description,
			CreatedAt:   data.Category.CreatedAt,
		},
		Images: productImages,
	}

	if data.CampaignID != nil {
		product.Campaign = &masterschema.CampaignSchema{
			ID:          data.Campaign.ID,
			WorkspaceID: data.Campaign.WorkspaceID.String(),
			Title:       data.Campaign.Title,
			Key:         data.Campaign.Key,
			Slug:        data.Campaign.Slug,
			Description: data.Campaign.Description,
			IsPublish:   data.Campaign.IsPublish,
			CreatedAt:   &data.Campaign.CreatedAt,
		}
	}

	// send response
	return &product, nil
}

func (s *StoreService) CreateStoreProduct(userID string, body storeschema.ProductPayload) error {
	// check valid store
	store, err := s.findStore(userID, false)
	if err != nil {
		return err
	}

	// converting uuid-string into uuid-type
	uuidStoreID, err := uuid.Parse(store.ID)
	if err != nil {
		return err
	}

	uuidCategoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		return err
	}

	// prepare data for insert
	ID := uuid.New()
	arrID := strings.Split(ID.String(), "-")
	key := ""
	if len(arrID) > 0 {
		key = arrID[0]
	}

	data := models.StoreProducts{
		ID:          ID,
		StoreID:     uuidStoreID,
		CategoryID:  uuidCategoryID,
		Name:        body.Name,
		Slug:        slug.Make(body.Name),
		Key:         key,
		Description: body.Description,
		Price:       body.Price,
		Status:      body.Status,
		CreatedAt:   time.Now(),
	}

	campaignID := body.CampaignID
	if campaignID != nil {
		uuidCampaignID, err := uuid.Parse(*campaignID)
		if err != nil {
			return err
		}
		data.CampaignID = &uuidCampaignID
	}

	// prepare to insert images
	dataImages := []models.StoreProductImages{}
	if body.Images != nil {
		for i, v := range body.Images {
			imageID := uuid.New()
			fileName, err := utils.UploadImage(v.FileName, "products", fmt.Sprintf("%s-%s", data.Key, strings.ReplaceAll(imageID.String(), "-", "")))
			if err != nil {
				return fmt.Errorf("failure to upload image number %d with detail error: %s", i, err.Error())
			}
			dataImages = append(dataImages, models.StoreProductImages{
				ID:             imageID,
				StoreProductID: data.ID,
				FileName:       *fileName,
				CreatedAt:      data.CreatedAt,
			})
		}
	}

	// perform to insert data
	err = s.storeRepo.CreateProduct(data)
	if err != nil {
		return err
	}

	err = s.storeRepo.CreateProductImages(dataImages)
	if err != nil {
		return err
	}

	// return success response
	// by set as no-error
	return nil
}

func (s *StoreService) UpdateStoreProduct(userID string, ID string, body storeschema.ProductPayload) error {
	// check valid store
	store, err := s.findStore(userID, false)
	if err != nil {
		return err
	}

	// check existing product
	product, err := s.storeRepo.FindStoreProduct(store.ID, ID)
	if err != nil {
		return err
	}

	// converting uuid-string data to uuid-type
	uuidCategoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		return err
	}

	// prepare data for update
	now := time.Now()
	data := models.StoreProducts{
		CategoryID:  uuidCategoryID,
		Name:        body.Name,
		Slug:        slug.Make(body.Name),
		Description: body.Description,
		Price:       body.Price,
		Status:      body.Status,
		UpdatedAt:   &now,
	}

	campaignID := body.CampaignID
	if campaignID != nil {
		uuidCampaignID, err := uuid.Parse(*campaignID)
		if err != nil {
			return err
		}
		data.CampaignID = &uuidCampaignID
	}

	// perform to update data
	err = s.storeRepo.UpdateStoreProduct(data, store.ID, ID)
	if err != nil {
		return err
	}

	// perform to delete images not exist from database
	images, err := s.storeRepo.FindImagesByProduct(ID)
	if err != nil {
		return err
	}
	for _, v := range images {
		isExists := false
		if body.Images != nil {
			for _, j := range body.Images {
				if j.ID != nil && v.ID.String() == *j.ID {
					isExists = true
				}
			}
		}
		if !isExists {
			// delete data from database
			// and remove image from directory
			err = s.storeRepo.DeleteProductImage(v.ID.String())
			if err != nil {
				return err
			}

			err := utils.RemoveImage(v.FileName)
			if err != nil {
				return err
			}
		}
	}

	// perform to create new images
	dataImages := []models.StoreProductImages{}
	if body.Images != nil {
		for i, v := range body.Images {
			if v.ID == nil {
				imageID := uuid.New()
				fileName, err := utils.UploadImage(v.FileName, "products", fmt.Sprintf("%s-%s", product.Key, strings.ReplaceAll(imageID.String(), "-", "")))
				if err != nil {
					return fmt.Errorf("failure to upload image number %d with detail error: %s", i, err.Error())
				}
				dataImages = append(dataImages, models.StoreProductImages{
					ID:             imageID,
					StoreProductID: product.ID,
					FileName:       *fileName,
					CreatedAt:      data.CreatedAt,
				})
			}
		}
	}

	err = s.storeRepo.CreateProductImages(dataImages)
	if err != nil {
		return err
	}

	// return success response
	// by set as no-error
	return nil
}

func (s *StoreService) DeleteStoreProduct(userID string, ID string) error {
	// check valid store
	store, err := s.findStore(userID, false)
	if err != nil {
		return err
	}

	// perform to set data as deleted
	now := time.Now()
	data := models.StoreProducts{
		Deleted:   true,
		UpdatedAt: &now,
	}
	err = s.storeRepo.UpdateStoreProduct(data, store.ID, ID)
	if err != nil {
		return err
	}

	// return success response
	// by set as no-error
	return nil
}

func (s *StoreService) FindStoreProductFormEntries(userID string, ID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	// check valid store
	_, err := s.findStore(userID, false)
	if err != nil {
		return nil, err
	}

	// perform to get form entries of prroduct
	list, err := s.storeRepo.FindStoreProductFormEntries(ID, params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.storeRepo.FindCountStoreProductFormEntries(ID, params)
	if err != nil {
		return nil, err
	}

	// prepare response list
	totalPage := 1
	if count > 0 && params.Limit > 0 {
		totalPage = int(math.Ceil(float64(count) / float64(params.Limit)))
	}

	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  totalPage,
		Rows:       list,
	}

	// return success response
	return &response, nil
}

func (s *StoreService) FindStoreCategoriesByKey(key string) ([]storeschema.ProductCategoryResponse, error) {
	store, err := s.storeRepo.FindStoreByKey(key)
	if err != nil {
		return nil, err
	}

	list, err := s.storeRepo.FindStoreProductCategories(store.ID.String(), nil)
	if err != nil {
		return nil, err
	}

	var data []storeschema.ProductCategoryResponse
	for _, v := range list {
		data = append(data, storeschema.ProductCategoryResponse{
			ID:            v.ID.String(),
			Name:          v.Name,
			Description:   v.Description,
			CreatedAt:     v.CreatedAt,
			TotalProducts: 0,
		})
	}

	return data, nil
}
