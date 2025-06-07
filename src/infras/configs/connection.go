package configs

import (
	"fmt"
	"kiraform/src/infras/migrations"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connection(config Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME)

	gormConfig := &gorm.Config{}
	if config.APP_ENV != "" && strings.ToLower(config.APP_ENV) == "dev" {
		// only show query log on dev mode
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}

	DB, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// start migrating table
	if config.MIGRATION {
		migrations.Migrate(DB)
	}

	// start seeding data
	if config.SEEDER {
		migrations.Seeder(DB)
	}

	// send db to global connection
	return DB
}
