package masterusecase

import (
	masterrepo "kiraform/src/applications/repos/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"math"
)

type WorkspaceUsecase interface {
	GetWorkspaces(params *commonschema.QueryParams) (*commonschema.ResponseList, error)
}

type WorkspaceService struct {
	workspaceRepo masterrepo.WorkspaceRepository
}

func NewWorkspaceUsecase(workspaceRepo masterrepo.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
	}
}

func (s *WorkspaceService) GetWorkspaces(params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
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
	count, err := s.workspaceRepo.FindCountWorkspaces(params)
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
