package masterusecase

import (
	"kiraform/src/applications/models"
	masterrepo "kiraform/src/applications/repos/masters"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"

	"github.com/google/uuid"
)

type FormEntryUsecase interface {
	EntryForm(campaignID string, userID *string, body []masterschema.FormEntryPayload) error
	GetHistory(userID string) ([]masterschema.FormEntrySchema, error)
	GetDetailHistory(userID string, ID string) (*masterschema.FormEntrySchema, error)
}

type FormEntryService struct {
	formEntryRepo masterrepo.FormEntryRepository
}

func NewFormEntryUsecase(formEntryRepo masterrepo.FormEntryRepository) *FormEntryService {
	return &FormEntryService{
		formEntryRepo: formEntryRepo,
	}
}

func (s *FormEntryService) EntryForm(campaignID string, userID *string, body []masterschema.FormEntryPayload) error {
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
	formEntry := models.FormEntries{
		ID:         uuid.New(),
		CampaignID: UUIDcampaignID,
		UserID:     UUIDuserID,
		Status:     "S1", // static as pending
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

		formDetailEntries = append(formDetailEntries, models.FormDetailEntries{
			ID:                      uuid.New(),
			FormEntryID:             formEntry.ID,
			CampaignFormID:          UUIDcampaignFormID,
			CampaignFormAttributeID: UUIDcampaignFormAttributeID,
			Value:                   v.Value,
		})
	}

	// perform to insert data
	if err := s.formEntryRepo.EntryForm(formEntry, formDetailEntries); err != nil {
		return err
	}
	return nil
}

func (s *FormEntryService) GetHistory(userID string) ([]masterschema.FormEntrySchema, error) {
	return nil, nil
}

func (s *FormEntryService) GetDetailHistory(userID string, ID string) (*masterschema.FormEntrySchema, error) {
	return nil, nil
}
