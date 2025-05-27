package masterusecase

import (
	"kiraform/src/applications/models"
	masterrepo "kiraform/src/applications/repos/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type WorkspaceUsecase interface {
	FindWorkspaces(params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindWorkspaceByID(ID string) (*models.Workspaces, error)
	CreateWorkspace(body masterschema.WorkspacePayload) error
	UpdateWorkspace(ID string, body masterschema.WorkspacePayload) error
	DeleteWorkspace(ID string) error
}

type WorkspaceService struct {
	workspaceRepo masterrepo.WorkspaceRepository
}

func NewWorkspaceUsecase(workspaceRepo masterrepo.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
	}
}

func (s *WorkspaceService) FindWorkspaces(params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.workspaceRepo.FindWorkspaces(params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.workspaceRepo.FindCountWorkspace(params)
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

func (s *WorkspaceService) FindWorkspaceByID(ID string) (*models.Workspaces, error) {
	data, err := s.workspaceRepo.FindWorkspaceByID(ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *WorkspaceService) CreateWorkspace(body masterschema.WorkspacePayload) error {
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
