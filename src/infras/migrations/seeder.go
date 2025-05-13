package migrations

import (
	"fmt"
	"kiraform/src/infras/migrations/seeders"
	"log"

	"gorm.io/gorm"
)

func Seeder(DB *gorm.DB) {
	var err error

	roleID, err := seeders.Roles(DB)
	if err != nil {
		log.Fatal("Error while seeding data roles")
	}

	err = seeders.Forms(DB)
	if err != nil {
		log.Fatal("Error while seeding data forms")
	}

	err = seeders.Packages(DB)
	if err != nil {
		log.Fatal("Error while seeding data packages")
	}

	err = seeders.Users(DB, roleID)
	if err != nil {
		log.Fatal("Error while seeding data users")
	}

	fmt.Println("Seeding data is complete!")
}
