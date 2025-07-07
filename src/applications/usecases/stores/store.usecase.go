package storeusecase

import (
	"errors"
	"fmt"
	"kiraform/src/applications/models"
	storerepo "kiraform/src/applications/repos/stores"
	storeschema "kiraform/src/interfaces/rest/schemas/stores"
	"kiraform/src/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type StoreUsecase interface {
	FindStore(userID string) (*storeschema.StoreResponse, error)
	UpdateStore(userID string, body storeschema.StorePayload) error
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
