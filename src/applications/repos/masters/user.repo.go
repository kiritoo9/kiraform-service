package userrepo

import (
	"kiraform/src/applications/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*models.Users, error)
	GetRoleByUser(userID uuid.UUID) (*models.UserRoles, error)
}

type Query struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &Query{DB: db}
}

func (q *Query) FindByEmail(email string) (*models.Users, error) {
	var user models.Users

	if err := q.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *Query) GetRoleByUser(userID uuid.UUID) (*models.UserRoles, error) {
	var userRole models.UserRoles
	if err := q.DB.
		Where("deleted = ? AND user_id = ?", false, userID).
		Preload("Role"). // join table
		First(&userRole).Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}
