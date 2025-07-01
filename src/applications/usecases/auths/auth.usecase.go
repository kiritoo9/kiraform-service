package authusecase

import (
	"errors"
	"kiraform/src/applications/models"
	repomasters "kiraform/src/applications/repos/masters"
	"kiraform/src/infras/configs"
	authschema "kiraform/src/interfaces/rest/schemas/auths"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Login(body authschema.LoginPayload) (*string, error)
	Register(body authschema.RegisterPayload) (*string, error)
}

type AuthService struct {
	UserRepo repomasters.UserRepository
	RoleRepo repomasters.RoleRepository
}

func NewAuthUsecase(userRepo repomasters.UserRepository, roleRepo repomasters.RoleRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
}

func (s *AuthService) Login(body authschema.LoginPayload) (*string, error) {
	// get data
	data, err := s.UserRepo.FindUserByEmail(body.Email)
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
	key := []byte(configs.Environment().SECRET_KEY)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func (s *AuthService) Register(body authschema.RegisterPayload) (*string, error) {
	// check existing email
	// validate to return error query, not error not found
	user, err := s.UserRepo.FindUserByEmail(body.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if user != nil {
		return nil, errors.New("email is already taken, try another one")
	}

	// load data role[user]
	role, err := s.RoleRepo.FindRoleByName("user")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role for this registartion is not found, please contact admin")
		}
		return nil, err
	}

	// hasing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// generating user identity
	userID := uuid.New()
	uuidStr := strings.Split(userID.String(), "-")
	userIdentity := ""
	if len(uuidStr) > 1 {
		userIdentity = strings.ToUpper(uuidStr[0] + uuidStr[1])
	}

	// preparing data user and user profile
	dataUser := models.Users{
		ID:           userID,
		UserIdentity: userIdentity,
		Email:        body.Email,
		Password:     string(hashedPassword),
		Fullname:     body.Fullname,
		IsActive:     true, // default true for now, next version needs to active it manually using OTP
		CreatedAt:    time.Now(),
	}

	firstName, middleName, lastName := "", "", ""
	nameParts := strings.Fields(body.Fullname)

	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}

	if len(nameParts) > 1 {
		middleName = nameParts[1]
	}

	if len(lastName) > 2 {
		lastName = nameParts[2]
	}

	dataUserProfile := models.UserProfiles{
		ID:         uuid.New(),
		UserID:     dataUser.ID,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		CreatedAt:  time.Now(),
	}

	dataUserRole := models.UserRoles{
		ID:        uuid.New(),
		UserID:    dataUser.ID,
		RoleID:    role.ID,
		CreatedAt: time.Now(),
	}

	// perform to insert data
	err = s.UserRepo.CreateUser(dataUser, dataUserProfile, dataUserRole)
	if err != nil {
		return nil, err
	}

	// response
	responseMsg := "Your account is successfully registered"
	return &responseMsg, nil
}
