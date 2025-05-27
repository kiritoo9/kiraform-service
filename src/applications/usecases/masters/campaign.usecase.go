package masterusecase

import (
	masterrepo "kiraform/src/applications/repos/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"math"
)

type CampaignUsecase interface {
	FindCampaigns(workspaceID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
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
