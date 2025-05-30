package masterusecase

import (
	"errors"
	"kiraform/src/applications/models"
	masterrepo "kiraform/src/applications/repos/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type CampaignUsecase interface {
	FindCampaigns(workspaceID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindCampaign(workspaceID string, ID string) (*masterschema.DetailCampaignSchema, error)
	FindFormsByCampaign(campaignID string) ([]masterschema.DetailCampaignFormSchema, error)
	FindFormAttributes(campaignFormID string) ([]masterschema.CampaignFormAttributeSchemas, error)
	CreateCampaign(workspaceID string, body masterschema.CampaignPayload) error
	UpdateCampaign(workspaceID string, ID string, body masterschema.CampaignPayload) error
	DeleteCampaign(workspaceID string, ID string) error
	FindCampaignSeos(campaignID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindCampaignSeoByID(campaignID string, ID string) (*masterschema.CampaignSeoSchema, error)
	CreateCampaignSeo(campaignID string, body masterschema.CampaignSeoPayload) error
	UpdateCampaignSeo(campaignID string, ID string, body masterschema.CampaignSeoPayload) error
	DeleteCampaignSeo(campaignID string, ID string) error
}

type CampaignService struct {
	campaignRepo masterrepo.CampaignRepository
}

func NewCampaignUsecase(campaignRepo masterrepo.CampaignRepository) *CampaignService {
	return &CampaignService{
		campaignRepo: campaignRepo,
	}
}

func (s *CampaignService) FindCampaigns(workspaceID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.campaignRepo.FindCampaigns(workspaceID, params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.campaignRepo.FindCountCampaign(workspaceID, params)
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

func (s *CampaignService) FindCampaign(workspaceID string, ID string) (*masterschema.DetailCampaignSchema, error) {
	data, err := s.campaignRepo.FindCampaignByID(workspaceID, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("records not found")
		} else {
			return nil, err
		}
	}
	return &masterschema.DetailCampaignSchema{
		ID:          data.ID,
		WorkspaceID: data.WorkspaceID,
		Title:       data.Title,
		Description: data.Description,
		IsPublish:   data.IsPublish,
		CreatedAt:   data.CreatedAt,
	}, nil
}

func (s *CampaignService) FindFormsByCampaign(campaignID string) ([]masterschema.DetailCampaignFormSchema, error) {
	data, err := s.campaignRepo.FindFormsByCampaign(campaignID)
	if err != nil {
		return nil, err
	}

	// prepare response
	response := []masterschema.DetailCampaignFormSchema{}
	for _, v := range data {
		response = append(response, masterschema.DetailCampaignFormSchema{
			ID:           v.ID,
			FormID:       v.FormID,
			FormCode:     v.FormCode,
			FormName:     v.FormName,
			Title:        v.Title,
			Description:  v.Description,
			Placeholder:  v.Placeholder,
			DefaultValue: v.DefaultValue,
			IsRequired:   v.IsRequired,
			IsMultiple:   v.IsMultiple,
			CreatedAt:    v.CreatedAt,
		})
	}

	// send response
	return response, nil
}

func (s *CampaignService) FindFormAttributes(campaignFormID string) ([]masterschema.CampaignFormAttributeSchemas, error) {
	data, err := s.campaignRepo.FindFormAttributes(campaignFormID)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *CampaignService) CreateCampaign(workspaceID string, body masterschema.CampaignPayload) error {
	// prepare usable data
	campaignID := uuid.New()
	campaignIDarr := strings.Split(campaignID.String(), "-")
	key := ""
	if len(campaignIDarr) > 0 {
		key = campaignIDarr[0]
	}

	UUIDworkspaceID, err := uuid.Parse(workspaceID)
	if err != nil {
		return err
	}

	thumbnail := ""
	tArr := strings.Split(body.Title, " ")
	for _, t := range tArr {
		thumbnail += t[0:1] // get first character
	}

	// prepare data for campaign header
	campaign := models.Campaigns{
		ID:          campaignID,
		WorkspaceID: UUIDworkspaceID,
		Title:       body.Title,
		Key:         key,
		Slug:        slug.Make(body.Title),
		Description: body.Description,
		IsPublish:   body.IsPublish,
		Thumbnail:   thumbnail,
	}

	// prepare data for campaign forms and campaign form attributes
	var campaignForms []models.CampaignForms
	var campaignFormAttributes []models.CampaignFormAttributes

	for _, v := range body.Forms {
		formID, err := uuid.Parse(v.FormID)
		if err != nil {
			return err
		}

		cf := models.CampaignForms{
			ID:           uuid.New(),
			CampaignID:   campaignID,
			FormID:       formID,
			Title:        v.Title,
			Description:  v.Description,
			Placeholder:  v.Placeholder,
			DefaultValue: v.DefaultValue,
			IsRequired:   v.IsRequired,
			IsMultiple:   v.IsMultiple,
		}
		campaignForms = append(campaignForms, cf)

		// appending data attributes for this form
		for _, j := range *v.Attributes {
			fa := models.CampaignFormAttributes{
				ID:             uuid.New(),
				CampaignFormID: cf.ID,
				Label:          j.Label,
				Value:          j.Value,
				IsDefault:      j.IsDefault,
			}
			campaignFormAttributes = append(campaignFormAttributes, fa)
		}
	}

	// perform to insert entire data
	err = s.campaignRepo.CreateCampaign(campaign, campaignForms, campaignFormAttributes)
	if err != nil {
		return nil
	}

	return nil
}
func (s *CampaignService) UpdateCampaign(workspaceID string, ID string, body masterschema.CampaignPayload) error {
	// prepare usable data
	t := time.Now()
	thumbnail := ""
	tArr := strings.Split(body.Title, " ")
	for _, t := range tArr {
		thumbnail += t[0:1] // get first character
	}
	campaignID, err := uuid.Parse(ID)
	if err != nil {
		return err
	}

	// prepare data campaign
	campaign := models.Campaigns{
		Title:       body.Title,
		Slug:        slug.Make(body.Title),
		Description: body.Description,
		IsPublish:   body.IsPublish,
		Thumbnail:   thumbnail,
		UpdatedAt:   &t,
	}

	// prepare data campaign forms and
	campaignFormActions := map[string][]models.CampaignForms{
		"create": {},
		"update": {},
		"delete": {},
	}

	// check missing data from existing to flag as delete
	existingForms, err := s.campaignRepo.FindFormsByCampaign(ID)
	if err != nil {
		return err
	}
	dataNotExists := []uuid.UUID{}
	for _, v := range existingForms {
		isExists := false
		for _, j := range body.Forms {
			if j.ID != nil {
				campaignFormID, err := uuid.Parse(*j.ID)
				if err != nil {
					return err
				}
				if v.ID == campaignFormID {
					isExists = true
					break
				}
			}
		}
		if !isExists {
			dataNotExists = append(dataNotExists, v.ID)
		}
	}

	// perform to create and update data campaign form
	for _, v := range body.Forms {
		formID, err := uuid.Parse(v.FormID)
		if err != nil {
			return err
		}

		cf := models.CampaignForms{
			FormID:       formID,
			Title:        v.Title,
			Description:  v.Description,
			Placeholder:  v.Placeholder,
			DefaultValue: v.DefaultValue,
			IsRequired:   v.IsRequired,
			IsMultiple:   v.IsMultiple,
		}
		if v.ID != nil {
			cf.UpdatedAt = &t
			campaignFormID, err := uuid.Parse(*v.ID)
			if err != nil {
				return err
			}
			cf.ID = campaignFormID
			campaignFormActions["update"] = append(campaignFormActions["update"], cf)
		} else {
			cf.ID = uuid.New()
			cf.CampaignID = campaignID
			campaignFormActions["create"] = append(campaignFormActions["create"], cf)
		}
	}

	// means data from payload not exist in database data
	// so delete the existing data
	// assumed user delete it from the frontend
	for _, v := range dataNotExists {
		campaignFormActions["delete"] = append(campaignFormActions["delete"], models.CampaignForms{
			ID:        v,
			Deleted:   true,
			UpdatedAt: &t,
		})
	}

	// perform to query for entire data
	if err := s.campaignRepo.UpdateEntireCampaign(ID, campaign, campaignFormActions); err != nil {
		return err
	}
	return nil
}

func (s *CampaignService) DeleteCampaign(workspaceID string, ID string) error {
	// check existing data
	_, err := s.campaignRepo.FindCampaignByID(workspaceID, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}

	// start update data
	t := time.Now()
	campaign := models.Campaigns{
		Deleted:   true,
		UpdatedAt: &t,
	}
	err = s.campaignRepo.UpdateCampaign(ID, campaign)
	if err != nil {
		return err
	}
	return nil
}

func (s *CampaignService) FindCampaignSeos(campaignID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.campaignRepo.FindCampaignSeos(campaignID, params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.campaignRepo.FindCountCampaignSeo(campaignID, params)
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

func (s *CampaignService) FindCampaignSeoByID(campaignID string, ID string) (*masterschema.CampaignSeoSchema, error) {
	data, err := s.campaignRepo.FindCampaignSeoByID(campaignID, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CampaignService) CreateCampaignSeo(campaignID string, body masterschema.CampaignSeoPayload) error {
	UUIDcampaignID, err := uuid.Parse(campaignID)
	if err != nil {
		return err
	}

	campaignSeo := models.CampaignSeos{
		ID:         uuid.New(),
		CampaignID: UUIDcampaignID,
		Platform:   body.Platform,
		Event:      body.Event,
		AccessKey:  body.AccessKey,
	}
	if err := s.campaignRepo.CreateCampaignSeo(campaignSeo); err != nil {
		return err
	}
	return nil
}

func (s *CampaignService) UpdateCampaignSeo(campaignID string, ID string, body masterschema.CampaignSeoPayload) error {
	t := time.Now()
	campaignSeo := models.CampaignSeos{
		Platform:  body.Platform,
		Event:     body.Event,
		AccessKey: body.AccessKey,
		UpdatedAt: &t,
	}
	if err := s.campaignRepo.UpdateCampaignSeo(campaignID, ID, campaignSeo); err != nil {
		return err
	}
	return nil
}

func (s *CampaignService) DeleteCampaignSeo(campaignID string, ID string) error {
	t := time.Now()
	campaignSeo := models.CampaignSeos{
		Deleted:   true,
		UpdatedAt: &t,
	}
	if err := s.campaignRepo.UpdateCampaignSeo(campaignID, ID, campaignSeo); err != nil {
		return err
	}
	return nil
}
