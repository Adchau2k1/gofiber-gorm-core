package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// Config func to get env value
func LoadConfig(config *Config) {
	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading .env file!")
		}

		config.DBPort = os.Getenv("DB_PORT")
		config.DBHost = os.Getenv("DB_HOST")
		config.DBUser = os.Getenv("DB_USER")
		config.DBPass = os.Getenv("DB_PASS")
		config.DBName = os.Getenv("DB_NAME")

	} else {
		config.DBPort = os.Getenv("DB_PORT")
		config.DBUser = os.Getenv("DB_USER")
		config.DBPass = os.Getenv("DB_PASS")
		config.DBName = os.Getenv("DB_NAME")
		config.DBHost = os.Getenv("DB_HOST")
	}
}

func GetConfigByKey(key string) string {
	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading .env file!")
		}

	}

	return os.Getenv(key)
}
