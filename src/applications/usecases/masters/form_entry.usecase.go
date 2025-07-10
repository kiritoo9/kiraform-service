package masterusecase

import (
	"kiraform/src/applications/models"
	masterrepo "kiraform/src/applications/repos/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"math"

	"github.com/google/uuid"
)

type FormEntryUsecase interface {
	EntryForm(campaignID string, userID *string, body []masterschema.FormEntryPayload, productID *string) error
	GetHistory(userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	GetDetailHistory(userID string, ID string) (*masterschema.FormEntryResponse, error)
}

type FormEntryService struct {
	formEntryRepo masterrepo.FormEntryRepository
}

func NewFormEntryUsecase(formEntryRepo masterrepo.FormEntryRepository) *FormEntryService {
	return &FormEntryService{
		formEntryRepo: formEntryRepo,
	}
}

func (s *FormEntryService) EntryForm(campaignID string, userID *string, body []masterschema.FormEntryPayload, productID *string) error {
	UUIDcampaignID, err := uuid.Parse(campaignID)
	if err != nil {
		return err
	}

	var UUIDuserID *uuid.UUID
	if userID != nil && *userID != "" {
		_userID, err := uuid.Parse(*userID)
		if err != nil {
			return err
		}
		UUIDuserID = &_userID
	}

	// preparing data
	formEntryID := uuid.New()
	formEntry := map[string]any{
		"id":          formEntryID,
		"user_id":     UUIDuserID,
		"campaign_id": UUIDcampaignID,
		"status":      "S1", // static as pending
	}

	if productID != nil && *productID != "" {
		UUIDproductID, err := uuid.Parse(*productID)
		if err != nil {
			return err
		}
		formEntry["product_id"] = &UUIDproductID
	}

	var formDetailEntries []models.FormDetailEntries
	for _, v := range body {
		UUIDcampaignFormID, err := uuid.Parse(v.CampaignFormID)
		if err != nil {
			return err
		}

		var UUIDcampaignFormAttributeID *uuid.UUID
		if v.CampaignFormAttributeID != nil && *v.CampaignFormAttributeID != "" {
			_attributeID, err := uuid.Parse(*v.CampaignFormAttributeID)
			if err != nil {
				return err
			}
			UUIDcampaignFormAttributeID = &_attributeID
		}

		fde := models.FormDetailEntries{
			ID:             uuid.New(),
			FormEntryID:    formEntryID,
			CampaignFormID: UUIDcampaignFormID,
			Value:          v.Value,
		}
		if UUIDcampaignFormAttributeID != nil {
			fde.CampaignFormAttributeID = UUIDcampaignFormAttributeID
		}

		formDetailEntries = append(formDetailEntries, fde)
	}

	// perform to insert data
	if err := s.formEntryRepo.EntryForm(formEntry, formDetailEntries); err != nil {
		return err
	}
	return nil
}

func (s *FormEntryService) GetHistory(userID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.formEntryRepo.FindFormEntries(userID, params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.formEntryRepo.FindCountFormEntry(userID, params)
	if err != nil {
		return nil, err
	}
	totalPage := 1
	if count > 0 {
		totalPage = int(math.Ceil(float64(int(count)) / float64(params.Limit)))
	}

	// send response
	response.TotalPage = totalPage
	response.Rows = rows
	return &response, nil
}

func (s *FormEntryService) GetDetailHistory(userID string, ID string) (*masterschema.FormEntryResponse, error) {
	// get form entry header
	formEntry, err := s.formEntryRepo.FindFormEntry(userID, ID)
	if err != nil {
		return nil, err
	}

	// get detail form entries
	formDetailEntry, err := s.formEntryRepo.FindDetailFormEntry(formEntry.ID)
	if err != nil {
		return nil, err
	}

	// prepare for response
	return &masterschema.FormEntryResponse{
		Header: *formEntry,
		Detail: formDetailEntry,
	}, nil
}
