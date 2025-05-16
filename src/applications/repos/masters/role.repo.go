package masterrepo

import (
	"kiraform/src/applications/models"
	"strings"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindRoleByName(name string) (*models.Roles, error)
}

type RoleQuery struct {
	DB *gorm.DB
}

func NewRoleRepository(DB *gorm.DB) RoleRepository {
	return &RoleQuery{DB: DB}
}

func (q *RoleQuery) FindRoleByName(name string) (*models.Roles, error) {
	var role models.Roles
	if err := q.DB.Where("deleted = ? AND LOWER(name) = ?", false, strings.ToLower(name)).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
