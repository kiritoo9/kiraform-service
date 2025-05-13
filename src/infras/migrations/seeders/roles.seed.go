package seeders

import (
	"kiraform/src/applications/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Roles(DB *gorm.DB) (uuid.UUID, error) {
	roles := []models.Roles{
		{ID: uuid.New(), Name: "admin", Description: "Administrator", CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "user", Description: "User", CreatedAt: time.Now()},
	}

	for _, data := range roles {
		if err := DB.FirstOrCreate(&data, models.Roles{Name: data.Name}).Error; err != nil {
			return uuid.Nil, err
		}
	}

	return roles[0].ID, nil
}
