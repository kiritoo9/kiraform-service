package masterrepo

import (
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"strings"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	FindWorkspaces(userID *string, params *commonschema.QueryParams) ([]models.Workspaces, error)
	FindCountWorkspace(userID *string, params *commonschema.QueryParams) (int64, error)
	FindWorkspaceByID(ID string) (*models.Workspaces, error)
	CreateWorkspace(data models.Workspaces) error
	UpdateWorkspace(ID string, data models.Workspaces) error
	FindWorkspaceUsers(workspaceID string, params *commonschema.QueryParams) ([]masterschema.WorkspaceUserSchema, error)
	FindCountWorkspaceUser(workspaceID string, params *commonschema.QueryParams) (int64, error)
	FindWorkspaceUserByID(workspaceID string, ID string) (*masterschema.WorkspaceUserSchema, error)
	FindWorkspaceUserByUser(workspaceID string, userID string) (*masterschema.WorkspaceUserSchema, error)
	FindWorkspaceUserByUserApproved(workspaceID string, userID string) (*masterschema.WorkspaceUserSchema, error)
	CreateWorkspaceUser(data models.WorkspaceUsers) error
	UpdateWorkspaceUser(workspaceID string, ID string, data models.WorkspaceUsers) error
	FindCountCampaignByWorkspace(workspaceID string) (int64, error)
	FindCountFormSubmissionByWorkspace(workspaceID string) (int64, error)
}

type WorkspaceQuery struct {
	DB *gorm.DB
}

func NewWorkspaceRepository(DB *gorm.DB) *WorkspaceQuery {
	return &WorkspaceQuery{DB: DB}
}

func (q *WorkspaceQuery) FindWorkspaces(userID *string, params *commonschema.QueryParams) ([]models.Workspaces, error) {
	var workspaces []models.Workspaces

	// calculating offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statments
	st := q.DB.Model(&models.Workspaces{}).Where("deleted = ?", false)
	if userID != nil {
		st = st.Where("workspaces.id IN (SELECT workspace_id FROM workspace_users WHERE workspace_users.workspace_id = workspaces.id AND deleted = ? AND user_id = ?)", false, userID)
	}

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

func (q *WorkspaceQuery) FindCountWorkspace(userID *string, params *commonschema.QueryParams) (int64, error) {
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

func (q *WorkspaceQuery) FindWorkspaceByID(ID string) (*models.Workspaces, error) {
	var workspace models.Workspaces

	// preparing query
	err := q.DB.Where("deleted = ? AND id = ?", false, ID).First(&workspace).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (q *WorkspaceQuery) CreateWorkspace(data models.Workspaces) error {
	if err := q.DB.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *WorkspaceQuery) UpdateWorkspace(ID string, data models.Workspaces) error {
	if err := q.DB.Model(&models.Workspaces{}).Where("id = ?", ID).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *WorkspaceQuery) FindWorkspaceUsers(workspaceID string, params *commonschema.QueryParams) ([]masterschema.WorkspaceUserSchema, error) {
	var workspaces []masterschema.WorkspaceUserSchema

	// calculating offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statments
	st := q.DB.Model(&models.WorkspaceUsers{}).Where("workspace_users.deleted = ? AND workspace_users.workspace_id = ?", false, workspaceID).
		Select("workspace_users.*", "users.email AS user_email", "users.fullname AS user_name", "workspaces.title AS workspace_title").
		Joins("JOIN users ON users.id = workspace_users.user_id").
		Joins("JOIN workspaces ON workspaces.id = workspace_users.workspace_id")

	// add search condition
	if params.Search != "" {
		st = st.Where("(LOWER(users.email) LIKE ? OR LOWER(users.fullname) LIKE ?)", "%"+strings.ToLower(params.Search)+"%", "%"+strings.ToLower(params.Search)+"%")
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

func (q *WorkspaceQuery) FindCountWorkspaceUser(workspaceID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// preparing conditions
	st := q.DB.Model(&models.WorkspaceUsers{}).Where("workspace_users.deleted = ? AND workspace_users.workspace_id = ?", false, workspaceID).
		Joins("JOIN users ON users.id = workspace_users.user_id")

	if params.Search != "" {
		st = st.Where("(LOWER(users.email) LIKE ? OR LOWER(users.fullname) LIKE ?)", "%"+strings.ToLower(params.Search)+"%", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}

	// send back
	return count, nil
}

func (q *WorkspaceQuery) FindWorkspaceUserByID(workspaceID string, ID string) (*masterschema.WorkspaceUserSchema, error) {
	var workspaceUser masterschema.WorkspaceUserSchema

	// preparing query
	err := q.DB.Model(&models.WorkspaceUsers{}).Where("workspace_users.deleted = ? AND workspace_users.workspace_id = ? AND workspace_users.id = ?", false, workspaceID, ID).
		Select("workspace_users.*", "users.email AS user_email", "users.fullname AS user_name", "workspaces.title AS workspace_title").
		Joins("JOIN users ON users.id = workspace_users.user_id").
		Joins("JOIN workspaces ON workspaces.id = workspace_users.workspace_id").
		First(&workspaceUser).Error
	if err != nil {
		return nil, err
	}
	return &workspaceUser, nil
}

func (q *WorkspaceQuery) FindWorkspaceUserByUser(workspaceID string, userID string) (*masterschema.WorkspaceUserSchema, error) {
	var workspaceUser masterschema.WorkspaceUserSchema

	// preparing query
	err := q.DB.Model(&models.WorkspaceUsers{}).Where("workspace_users.deleted = ? AND workspace_users.workspace_id = ? AND workspace_users.user_id = ?", false, workspaceID, userID).
		Select("workspace_users.*", "users.email AS user_email", "users.fullname AS user_name", "workspaces.title AS workspace_title").
		Joins("JOIN users ON users.id = workspace_users.user_id").
		Joins("JOIN workspaces ON workspaces.id = workspace_users.workspace_id").
		First(&workspaceUser).Error
	if err != nil {
		return nil, err
	}
	return &workspaceUser, nil
}

func (q *WorkspaceQuery) FindWorkspaceUserByUserApproved(workspaceID string, userID string) (*masterschema.WorkspaceUserSchema, error) {
	var workspaceUser masterschema.WorkspaceUserSchema

	// preparing query
	err := q.DB.Model(&models.WorkspaceUsers{}).Where("workspace_users.deleted = ? AND workspace_users.workspace_id = ? AND workspace_users.user_id = ? AND workspace_users.status IN (?, ?)", false, workspaceID, userID, "S3", "S5").
		Select("workspace_users.*", "users.email AS user_email", "users.fullname AS user_name", "workspaces.title AS workspace_title").
		Joins("JOIN users ON users.id = workspace_users.user_id").
		Joins("JOIN workspaces ON workspaces.id = workspace_users.workspace_id").
		First(&workspaceUser).Error
	if err != nil {
		return nil, err
	}
	return &workspaceUser, nil
}

func (q *WorkspaceQuery) CreateWorkspaceUser(data models.WorkspaceUsers) error {
	if err := q.DB.Model(&models.WorkspaceUsers{}).Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *WorkspaceQuery) UpdateWorkspaceUser(workspaceID string, ID string, data models.WorkspaceUsers) error {
	if err := q.DB.Model(&models.WorkspaceUsers{}).Where("workspace_id = ? AND id = ?", workspaceID, ID).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (q *WorkspaceQuery) FindCountCampaignByWorkspace(workspaceID string) (int64, error) {
	var count int64
	if err := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id = ?", false, workspaceID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *WorkspaceQuery) FindCountFormSubmissionByWorkspace(workspaceID string) (int64, error) {
	var count int64
	if err := q.DB.Raw("SELECT COUNT(1) FROM form_entries JOIN campaigns ON campaigns.id = form_entries.campaign_id WHERE campaigns.deleted = ? AND form_entries.deleted = ? AND campaigns.workspace_id = ?", false, false, workspaceID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
