package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_NAME  string
	APP_PORT  string
	APP_ENV   string
	MIGRATION bool
	SEEDER    bool

	DB_HOST string
	DB_USER string
	DB_PASS string
	DB_NAME string
	DB_PORT string

	SECRET_KEY string
}

func Environment() Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(fmt.Sprintf(".env.%s", env))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// load environment variable
	MIGRATION := true
	envMigration := os.Getenv("MIGRATION")
	if envMigration != "" {
		m, err := strconv.ParseBool(envMigration)
		if err == nil {
			MIGRATION = m
		}
	}

	SEEDER := false
	envSeeder := os.Getenv("SEEDER")
	if envSeeder != "" {
		s, err := strconv.ParseBool(envSeeder)
		if err == nil {
			SEEDER = s
		}
	}

	APP_ENV := "dev" // default
	envApp := os.Getenv("APP_ENV")
	if envApp != "" {
		APP_ENV = envApp
	}

	return Config{
		APP_NAME:  os.Getenv("APP_NAME"),
		APP_PORT:  os.Getenv("APP_PORT"),
		APP_ENV:   APP_ENV,
		MIGRATION: MIGRATION,
		SEEDER:    SEEDER,

		DB_HOST: os.Getenv("DB_HOST"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_NAME: os.Getenv("DB_NAME"),
		DB_PORT: os.Getenv("DB_PORT"),

		SECRET_KEY: os.Getenv("SECRET_KEY"),
	}
}
