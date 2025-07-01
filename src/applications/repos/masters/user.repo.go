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
	FindCountFormByUser(userID string) (int64, error)
	FindCountFormSubmitByUser(userID string) (int64, error)
	FindCountFormSubmittedByUser(userID string) (int64, error)
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

func (q *UserQuery) FindCountFormByUser(userID string) (int64, error) {
	query := `
		SELECT 
			COUNT(1)
		FROM campaigns
		JOIN workspaces ON workspaces.id = campaigns.workspace_id
		WHERE
			campaigns.deleted = ?
			AND workspaces.deleted = ?
			AND workspaces.id IN (
				SELECT workspace_users.workspace_id
				FROM workspace_users
				WHERE
					workspace_users.workspace_id = workspaces.id
					AND workspace_users.deleted = ?
					AND workspace_users.user_id = ?
					AND workspace_users.status IN (?, ?)
			)
	`
	args := []any{false, false, false, userID, "S3", "S5"}

	var count int64
	if err := q.DB.Raw(query, args...).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *UserQuery) FindCountFormSubmitByUser(userID string) (int64, error) {
	query := `
		SELECT
			COUNT(1)
		FROM form_entries
		JOIN campaigns ON campaigns.id = form_entries.campaign_id
		JOIN workspaces ON workspaces.id = campaigns.workspace_id
		WHERE 
			form_entries.deleted = ?
			AND campaigns.deleted = ?
			AND workspaces.deleted = ?
			AND workspaces.id IN (
				SELECT workspace_users.workspace_id
				FROM workspace_users
				WHERE
					workspace_users.workspace_id = workspaces.id
					AND workspace_users.deleted = ?
					AND workspace_users.user_id = ?
					AND workspace_users.status IN (?, ?)
			)
	`
	args := []any{false, false, false, false, userID, "S3", "S5"}

	var count int64
	if err := q.DB.Raw(query, args...).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *UserQuery) FindCountFormSubmittedByUser(userID string) (int64, error) {
	query := `
		SELECT
			COUNT(1)
		FROM form_entries
		WHERE
			form_entries.deleted = ?
			AND form_entries.user_id = ?
	`
	args := []any{false, userID}

	var count int64
	if err := q.DB.Raw(query, args...).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
