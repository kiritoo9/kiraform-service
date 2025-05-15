package repomasters

import (
	"kiraform/src/applications/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(email string) (*models.Users, error)
	GetRoleByUser(userID uuid.UUID) (*models.UserRoles, error)
	CreateUser(user models.Users, userProfile models.UserProfiles, userRole models.UserRoles) error
}

type UserQuery struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &UserQuery{DB: DB}
}

func (q *UserQuery) FindUserByEmail(email string) (*models.Users, error) {
	var user models.Users

	if err := q.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
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
