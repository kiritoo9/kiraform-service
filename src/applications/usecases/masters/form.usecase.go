package masterusecase

import (
	masterrepo "kiraform/src/applications/repos/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"math"
)

type FormUsecase interface {
	FindForms(params *commonschema.QueryParams) (*commonschema.ResponseList, error)
	FindFormByID(ID string) (*masterschema.FormSchema, error)
}

type FormService struct {
	formRepo masterrepo.FormRepository
}

func NewFormUsecase(formRepo masterrepo.FormRepository) *FormService {
	return &FormService{
		formRepo: formRepo,
	}
}

func (s *FormService) FindForms(params *commonschema.QueryParams) (*commonschema.ResponseList, error) {
	response := commonschema.ResponseList{
		Parameters: *params,
		TotalPage:  1,
		Rows:       nil,
	}

	// get list data
	rows, err := s.formRepo.FindForms(params)
	if err != nil {
		return nil, err
	}

	// get count data
	count, err := s.formRepo.FindCountForm(params)
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

func (s *FormService) FindFormByID(ID string) (*masterschema.FormSchema, error) {
	data, err := s.formRepo.FindFormByID(ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}
