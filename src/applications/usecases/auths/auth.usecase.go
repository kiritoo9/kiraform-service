package authusecase

import (
	"errors"
	userrepo "kiraform/src/applications/repos/masters"
	authschema "kiraform/src/interfaces/rest/schemas/auths"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	Login(body authschema.LoginPayload) (*string, error)
}

type Service struct {
	UserRepo userrepo.Repository
}

func NewUsecase(userRepo userrepo.Repository) *Service {
	return &Service{
		UserRepo: userRepo,
	}
}

func (s *Service) Login(body authschema.LoginPayload) (*string, error) {
	// get data
	data, err := s.UserRepo.FindByEmail(body.Email)
	if err != nil {
		return nil, err
	}

	// validate matching password
	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(body.Password)); err != nil {
		return nil, errors.New("password does not match")
	}

	// get user role
	roleName := "user" // default
	role, err := s.UserRepo.GetRoleByUser(data.ID)
	if err == nil && role.Role.Name != "" {
		roleName = role.Role.Name
	}

	// prepare data as response
	response := map[string]any{
		"id":        data.ID,
		"role_name": roleName,
	}

	// convert into jwt token
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	for k, v := range response {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte("mykey")
	signedToken, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}
