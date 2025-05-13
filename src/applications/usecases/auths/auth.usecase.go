package authusecase

import (
	userrepo "kiraform/src/applications/repos/masters"
	authschema "kiraform/src/interfaces/rest/schemas/auths"
)

type Usecase interface {
	Login(body authschema.LoginPayload) error
}

type Service struct {
	UserRepo userrepo.Query
}

func NewUsecase(userRepo userrepo.Query) *Service {
	return &Service{
		UserRepo: userRepo,
	}
}

func (s *Service) Login(body authschema.LoginPayload) error {
	_, err := s.UserRepo.FindByEmail(body.Email)
	if err != nil {
		return err
	}
	return nil
}
