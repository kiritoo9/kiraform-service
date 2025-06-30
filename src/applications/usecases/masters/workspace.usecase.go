package masterusecase

import (
	"errors"
	"fmt"
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

type WorkspaceUsecase interface {
	FindWorkspaces(userID *string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindWorkspaceByID(ID string) (*models.Workspaces, error)
	CreateWorkspace(userID string, body masterschema.WorkspacePayload) error
	UpdateWorkspace(ID string, body masterschema.WorkspacePayload) error
	DeleteWorkspace(ID string) error
	FindWorkspaceUsers(workspaceID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindWorkspaceUserByID(workspaceID string, ID string) (*masterschema.WorkspaceUserSchema, error)
	CreateWorkspaceUser(workspaceID string, body masterschema.WorkspaceUserPayload) error
	UpdateWorkspaceUser(workspaceID string, ID string, body masterschema.WorkspaceUserUpdatePayload) error
	DeleteWorkspaceUser(workspaceID string, ID string) error
}

type WorkspaceService struct {
	workspaceRepo masterrepo.WorkspaceRepository
	userRepo      masterrepo.UserRepository
}

func NewWorkspaceUsecase(workspaceRepo masterrepo.WorkspaceRepository, userRepo masterrepo.UserRepository) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

func (s *WorkspaceService) FindWorkspaces(userID *string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.workspaceRepo.FindWorkspaces(userID, params)
	if err != nil {
		return nil, err
	}

	// converting format data from []models.Workspaces -> []masterschema.WorkspaceList
	var list []masterschema.WorkspaceList
	for _, v := range rows {
		// get total form from this workspace
		var totalForm int64
		countForm, err := s.workspaceRepo.FindCountCampaignByWorkspace(v.ID.String())
		if err != nil {
			return nil, err
		}
		totalForm = countForm

		// get total form created from entire campaign in this workspace
		var totalSubmit int64
		countSubmit, err := s.workspaceRepo.FindCountFormSubmissionByWorkspace(v.ID.String())
		if err != nil {
			return nil, err
		}
		totalSubmit = countSubmit

		// appending data
		list = append(list, masterschema.WorkspaceList{
			ID:          v.ID.String(),
			Title:       v.Title,
			Key:         v.Key,
			Slug:        v.Slug,
			Description: v.Description,
			Thumbnail:   v.Thumbnail,
			TotalForm:   totalForm,
			TotalSubmit: totalSubmit,
			CreatedAt:   v.CreatedAt,
		})
	}

	// get count data
	count, err := s.workspaceRepo.FindCountWorkspace(userID, params)
	if err != nil {
		return nil, err
	}
	totalPage := 1
	if count > 0 {
		totalPage = int(math.Ceil(float64(int(count)) / float64(params.Limit)))
	}

	// send response
	response.TotalPage = totalPage
	response.Rows = list
	return &response, nil
}

func (s *WorkspaceService) FindWorkspaceByID(ID string) (*models.Workspaces, error) {
	data, err := s.workspaceRepo.FindWorkspaceByID(ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *WorkspaceService) CreateWorkspace(userID string, body masterschema.WorkspacePayload) error {
	// prepare data to insert
	ID := uuid.New()
	arrOfID := strings.Split(ID.String(), "-")
	key := ""
	if len(arrOfID) > 0 {
		key = arrOfID[0]
	}

	data := models.Workspaces{
		ID:          ID,
		Title:       body.Title,
		Key:         key,
		Slug:        slug.Make(body.Title),
		Description: body.Description,
		IsPublish:   body.IsPublish,
		Thumbnail:   body.Thumbnail,
	}

	// perform to insert data
	err := s.workspaceRepo.CreateWorkspace(data)
	if err != nil {
		return err
	}

	// prepare and insert workspace user
	UUIDuserID, err := uuid.Parse(userID)
	if err == nil {
		wu := models.WorkspaceUsers{
			ID:          uuid.New(),
			UserID:      UUIDuserID,
			WorkspaceID: ID,
			Status:      "S5", // as an owner
		}
		_ = s.workspaceRepo.CreateWorkspaceUser(wu)
	}

	// set as success
	return nil
}

func (s *WorkspaceService) UpdateWorkspace(ID string, body masterschema.WorkspacePayload) error {
	// preparing data
	t := time.Now()
	data := models.Workspaces{
		Title:       body.Title,
		Slug:        slug.Make(body.Title),
		Description: body.Description,
		IsPublish:   body.IsPublish,
		Thumbnail:   body.Thumbnail,
		UpdatedAt:   &t,
	}

	// perform to update data
	err := s.workspaceRepo.UpdateWorkspace(ID, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) DeleteWorkspace(ID string) error {
	// check existing data
	_, err := s.FindWorkspaceByID(ID)
	if err != nil {
		return err
	}

	// start updating data
	t := time.Now()
	data := models.Workspaces{
		Deleted:   true,
		UpdatedAt: &t,
	}
	err = s.workspaceRepo.UpdateWorkspace(ID, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) FindWorkspaceUsers(workspaceID string, params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.workspaceRepo.FindWorkspaceUsers(workspaceID, params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.workspaceRepo.FindCountWorkspaceUser(workspaceID, params)
	if err != nil {
		return nil, err
	}
	totalPage := 1
	if count > 0 {
		totalPage = int(math.Ceil(float64(int(count)) / float64(params.Limit)))
	}

	// change value of status
	for i, v := range rows {
		status := strings.ToUpper(v.Status)
		if status == "S1" {
			rows[i].Status = "INVITED"
		} else if status == "S2" {
			rows[i].Status = "REQUESTED"
		} else if status == "S3" {
			rows[i].Status = "APPROVED"
		} else if status == "S4" {
			rows[i].Status = "REJECTED"
		}
	}

	// send response
	response.TotalPage = totalPage
	response.Rows = rows
	return &response, nil
}

func (s *WorkspaceService) FindWorkspaceUserByID(workspaceID string, ID string) (*masterschema.WorkspaceUserSchema, error) {
	data, err := s.workspaceRepo.FindWorkspaceUserByID(workspaceID, ID)
	if err != nil {
		return nil, err
	}

	// change value of status
	status := strings.ToUpper(data.Status)
	if status == "S1" {
		data.Status = "INVITED"
	} else if status == "S2" {
		data.Status = "REQUESTED"
	} else if status == "S3" {
		data.Status = "APPROVED"
	} else if status == "S4" {
		data.Status = "REJECTED"
	}

	// send response
	return data, nil
}

func (s *WorkspaceService) CreateWorkspaceUser(workspaceID string, body masterschema.WorkspaceUserPayload) error {
	// parse workspace id into uuid format
	UUIDworkspaceID, err := uuid.Parse(workspaceID)
	if err != nil {
		return err
	}

	// check user available by email or ID
	// priority check ID
	isExists := false
	var UUIDuserID uuid.UUID
	if body.UserID != nil {
		_, err := s.userRepo.FindUserByID(*body.UserID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		isExists = true

		uid, err := uuid.Parse(*body.UserID)
		if err != nil {
			return err
		}
		UUIDuserID = uid
	} else if body.UserEmail != nil {
		u, err := s.userRepo.FindUserByEmail(*body.UserEmail)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		isExists = true
		UUIDuserID = u.ID
	}

	if !isExists {
		return errors.New("user id or email is not found please try another user")
	}

	// check if user already registered in this workspace or not
	wu, err := s.workspaceRepo.FindWorkspaceUserByUser(workspaceID, UUIDuserID.String())
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	fmt.Println(wu)
	if wu != nil {
		return errors.New("this user already exists in this workspace")
	}

	// preparing data to insert
	data := models.WorkspaceUsers{
		ID:          uuid.New(),
		WorkspaceID: UUIDworkspaceID,
		UserID:      UUIDuserID,
		Status:      body.Status,
	}

	// perform to insert data
	err = s.workspaceRepo.CreateWorkspaceUser(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) UpdateWorkspaceUser(workspaceID string, ID string, body masterschema.WorkspaceUserUpdatePayload) error {
	// Only status that will updated in this section
	// User cannot update user_id
	t := time.Now()
	data := models.WorkspaceUsers{
		Status:    body.Status,
		UpdatedAt: t,
	}

	// perform to update data
	err := s.workspaceRepo.UpdateWorkspaceUser(workspaceID, ID, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) DeleteWorkspaceUser(workspaceID string, ID string) error {
	// check existing data
	_, err := s.FindWorkspaceUserByID(workspaceID, ID)
	if err != nil {
		return err
	}

	// start updating data
	t := time.Now()
	data := models.WorkspaceUsers{
		Deleted:   true,
		UpdatedAt: t,
	}
	err = s.workspaceRepo.UpdateWorkspaceUser(workspaceID, ID, data)
	if err != nil {
		return err
	}
	return nil
}
