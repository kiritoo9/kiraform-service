package masterrepo

import (
	"kiraform/src/applications/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(email string) (*models.Users, error)
	FindUserByID(ID string) (*models.Users, error)
	FindUserProfile(userID string) (*models.UserProfiles, error)
	GetRoleByUser(userID uuid.UUID) (*models.UserRoles, error)
	CreateUser(user models.Users, userProfile models.UserProfiles, userRole models.UserRoles) error
	UpdateUser(ID string, user models.Users) error
	CreateUserProfile(userProfile models.UserProfiles) error
	UpdateUserProfile(userID string, userProfile models.UserProfiles) error
}

type UserQuery struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &UserQuery{DB: DB}
}

func (q *UserQuery) FindUserByID(id string) (*models.Users, error) {
	var user models.Users

	if err := q.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *UserQuery) FindUserByEmail(email string) (*models.Users, error) {
	var user models.Users

	if err := q.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *UserQuery) FindUserProfile(userID string) (*models.UserProfiles, error) {
	var userProfile *models.UserProfiles
	if err := q.DB.Model(&models.UserProfiles{}).Where("deleted = ? AND user_id = ?", false, userID).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (q *UserQuery) GetRoleByUser(userID uuid.UUID) (*models.UserRoles, error) {
	var userRole models.UserRoles
	if err := q.DB.
		Where("deleted = ? AND user_id = ?", false, userID).
		Preload("Role"). // join table
		First(&userRole).Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (q *UserQuery) CreateUser(user models.Users, userProfile models.UserProfiles, userRole models.UserRoles) error {
	// perform to insert using transaction:rollback
	err := q.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&userProfile).Error; err != nil {
			return err
		}
		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (q *UserQuery) UpdateUser(ID string, user models.Users) error {
	if err := q.DB.Model(&models.Users{}).Where("deleted = ? AND id = ?", false, ID).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (q *UserQuery) CreateUserProfile(userProfile models.UserProfiles) error {
	if err := q.DB.Model(&models.UserProfiles{}).Create(&userProfile).Error; err != nil {
		return err
	}
	return nil
}

func (q *UserQuery) UpdateUserProfile(userID string, userProfile models.UserProfiles) error {
	if err := q.DB.Model(&models.UserProfiles{}).Where("deleted = ? AND user_id = ?", false, userID).Updates(&userProfile).Error; err != nil {
		return err
	}
	return nil
}
