package migrations

import (
	"fmt"
	"kiraform/src/applications/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	err := DB.AutoMigrate(
		&models.Users{},
		&models.Roles{},
	)
	if err != nil {
		log.Fatal("Error while migrating database: %v", err)
	}
	fmt.Println("Database successfully migrated!")
}
