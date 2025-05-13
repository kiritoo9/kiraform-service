package seeders

import (
	"kiraform/src/applications/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Users(DB *gorm.DB, roleID uuid.UUID) error {
	users := []models.Users{
		{ID: uuid.New(), UserIdentity: "KYYHSAG1999ADM", Email: "admin@admin.com", Password: "admin123", Fullname: "Administrator", IsActive: true, CreatedAt: time.Now()},
	}

	for _, data := range users {
		hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(hashed)

		if err := DB.FirstOrCreate(&data, models.Users{UserIdentity: data.UserIdentity}).Error; err != nil {
			return err
		} else {
			userRole := models.UserRoles{
				ID:        uuid.New(),
				UserID:    data.ID,
				RoleID:    roleID,
				CreatedAt: time.Now(),
			}
			DB.Create(&userRole)
		}
	}

	return nil
}
