package seeders

import (
	"kiraform/src/applications/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Packages(DB *gorm.DB) error {
	packages := []models.Packages{
		{ID: uuid.New(), Code: "FREEMIUM", Name: "Freemium", Description: "Free feature unlocked", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "GOLD", Name: "Gold Member", Description: "Gold feature unlocked", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "DIAMOND", Name: "Diamond Member", Description: "Diamond feature unlocked", CreatedAt: time.Now()},
	}

	for _, data := range packages {
		if err := DB.FirstOrCreate(&data, models.Packages{Code: data.Code}).Error; err != nil {
			return err
		}
	}

	return nil
}
