package masterrepo

import (
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"strings"

	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindCampaigns(workspaceID string, params *commonschema.QueryParams) ([]models.Campaigns, error)
	FindCountCampaign(workspaceID string, params *commonschema.QueryParams) (int64, error)
}

type CampaignQuery struct {
	DB *gorm.DB
}

func NewCampaignRepository(DB *gorm.DB) *CampaignQuery {
	return &CampaignQuery{
		DB: DB,
	}
}

func (q *CampaignQuery) FindCampaigns(workspaceID string, params *commonschema.QueryParams) ([]models.Campaigns, error) {
	var campaigns []models.Campaigns

	// define offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statemetns
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id::TEXT = ?", false, workspaceID)

	// add search condition
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// add orderby
	if params.OrderBy != "" {
		st = st.Order(params.OrderBy)
	}

	// add limit:offset
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get the data
	if err := st.Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (q *CampaignQuery) FindCountCampaign(workspaceID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// prepare condition
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id::TEXT = ?", false, workspaceID)
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
