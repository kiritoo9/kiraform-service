package seeders

import (
	"kiraform/src/applications/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Forms(DB *gorm.DB) error {
	forms := []models.Forms{
		{ID: uuid.New(), Code: "INPT_TEXT", Name: "Text", Description: "Free text input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "INPT_NUMBER", Name: "Number", Description: "Number input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "INPT_EMAIL", Name: "Email", Description: "Email input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "INPT_PASSWORD", Name: "Password", Description: "Password input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "INPT_FILE", Name: "Upload", Description: "File upload", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "SELC_OPTION", Name: "Select", Description: "Select option input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "SELC_RADIO", Name: "Radio", Description: "Radio button input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "CHCK_BOX", Name: "Checkbox", Description: "Checkbox for multiple input", CreatedAt: time.Now()},
		{ID: uuid.New(), Code: "TXT_AREA", Name: "Text Area", Description: "Textarea input", CreatedAt: time.Now()},
	}

	for _, data := range forms {
		if err := DB.FirstOrCreate(&data, models.Forms{Code: data.Code}).Error; err != nil {
			return err
		}
	}

	return nil
}
