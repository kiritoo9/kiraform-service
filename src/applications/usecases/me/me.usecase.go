package meusecase

import (
	"errors"
	"kiraform/src/applications/models"
	masterrepo "kiraform/src/applications/repos/masters"
	meschema "kiraform/src/interfaces/rest/schemas/me"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MeUsecase interface {
	GetProfile(userID string) (*meschema.MeResponse, error)
	UpdateProfile(userID string, body meschema.UserProfilePayload) error
	ChangePassword(userID string, body meschema.ChangePasswordPayload) error
}

type MeService struct {
	userrepo masterrepo.UserRepository
}

func NewMeUsecase(userrepo masterrepo.UserRepository) *MeService {
	return &MeService{userrepo: userrepo}
}

func (s *MeService) GetProfile(userID string) (*meschema.MeResponse, error) {
	// get user account
	user, err := s.userrepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	// get user profile
	var userProfile *meschema.UserProfile
	up, err := s.userrepo.FindUserProfile(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if up != nil {
		userProfile = &meschema.UserProfile{
			FirstName:   up.FirstName,
			MiddleName:  up.MiddleName,
			LastName:    up.LastName,
			Address:     up.Address,
			Phone:       up.Phone,
			Province:    up.Province,
			City:        up.City,
			District:    up.District,
			SubDistrict: up.SubDistrict,
			Avatar:      up.Avatar,
			UpdatedAt:   up.UpdatedAt,
		}
	}

	// get user roles
	var userRoles []meschema.UserRole
	ur, err := s.userrepo.GetRoleByUser(user.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if ur != nil {
		userRoles = append(userRoles, meschema.UserRole{
			RoleID:    ur.RoleID.String(),
			RoleName:  ur.Role.Name,
			CreatedAt: ur.CreatedAt,
		})
	}

	// prepare response data
	data := meschema.MeResponse{
		UserAccount: meschema.UserAccount{
			ID:           user.ID.String(),
			UserIdentity: user.UserIdentity,
			Email:        user.Email,
			Fullname:     user.Fullname,
			IsActive:     user.IsActive,
			CreatedAt:    user.CreatedAt,
		},
		UserProfile: userProfile,
		UserRoles:   userRoles,
	}
	return &data, nil
}

func (s *MeService) UpdateProfile(userID string, body meschema.UserProfilePayload) error {
	var avatar string
	t := time.Now()
	UUIDuserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// generate avatars name
	// by get first letter of user name {fist, middle, last}
	if len(body.FirstName) > 0 {
		avatar += body.FirstName[0:1]
	}

	if len(body.MiddleName) > 0 {
		avatar += body.MiddleName[0:1]
	}

	if len(body.LastName) > 0 {
		avatar += body.LastName[0:1]
	}

	// preparing data to update
	data := models.UserProfiles{
		FirstName:   body.FirstName,
		MiddleName:  body.MiddleName,
		LastName:    body.LastName,
		Address:     body.Address,
		Phone:       body.Phone,
		Province:    body.Province,
		City:        body.City,
		District:    body.District,
		SubDistrict: body.SubDistrict,
		Avatar:      avatar,
		UpdatedAt:   &t,
	}

	// check existing profile to decided action insert or update
	exists, err := s.userrepo.FindUserProfile(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if exists == nil {
		// complete the data
		data.ID = uuid.New()
		data.UserID = UUIDuserID

		// perform to create profile
		if err := s.userrepo.CreateUserProfile(data); err != nil {
			return err
		}
	} else {
		// perform to update profile
		if err := s.userrepo.UpdateUserProfile(userID, data); err != nil {
			return err
		}
	}
	return nil
}

func (s *MeService) ChangePassword(userID string, body meschema.ChangePasswordPayload) error {
	// check confirm password
	if body.NewPassword != body.ConfirmPassword {
		return errors.New("confirm password does not match")
	}

	// get user first
	user, err := s.userrepo.FindUserByID(userID)
	if err != nil {
		return err
	}

	// confirm hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		return errors.New("old password does not match, please input right password")
	}

	// prepare update data
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// perform to update user
	t := time.Now()
	newUser := models.Users{
		Password:  string(hashedPassword),
		UpdatedAt: &t,
	}
	err = s.userrepo.UpdateUser(userID, newUser)
	if err != nil {
		return err
	}
	return nil
}
