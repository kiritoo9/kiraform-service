package migrations

import (
	"fmt"
	"kiraform/src/applications/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	err := DB.AutoMigrate(
		&models.Users{}, &models.UserProfiles{},
		&models.Roles{}, &models.Packages{},
		&models.UserRoles{}, &models.UserPackages{},
		&models.Workspaces{}, &models.Campaigns{}, &models.Forms{},
		&models.CampaignSeos{}, &models.CampaignForms{}, &models.CampaignFormAttributes{},
		&models.CampaignFormEntries{},
		&models.Billings{}, &models.BillingDetails{},
	)
	if err != nil {
		log.Fatal(fmt.Printf("Error while migrating database: %v", err))
	}
	fmt.Println("Database successfully migrated!")
}
