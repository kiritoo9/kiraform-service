package masterrepo

import (
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"strings"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	FindWorkspaces(params *commonschema.QueryParams) ([]models.Workspaces, error)
	FindCountWorkspaces(params *commonschema.QueryParams) (int64, error)
}

type WorkspaceQuery struct {
	DB *gorm.DB
}

func NewWorkspaceRepository(DB *gorm.DB) *WorkspaceQuery {
	return &WorkspaceQuery{DB: DB}
}

func (q *WorkspaceQuery) FindWorkspaces(params *commonschema.QueryParams) ([]models.Workspaces, error) {
	var workspaces []models.Workspaces

	// calculating offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statments
	st := q.DB.Model(&models.Workspaces{}).Where("deleted = ?", false)

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
	if err := st.Find(&workspaces).Error; err != nil {
		return nil, err
	}

	// send back
	return workspaces, nil
}

func (q *WorkspaceQuery) FindCountWorkspaces(params *commonschema.QueryParams) (int64, error) {
	var count int64

	// preparing conditions
	st := q.DB.Model(&models.Workspaces{}).Where("deleted = ?", false)
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}

	// send back
	return count, nil
}
