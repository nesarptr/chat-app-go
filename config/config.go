package config

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	d, err := gorm.Open(postgres.Open(GetEnv("PGURI")), &gorm.Config{})
	if err != nil {
		return err
	}
	db = d
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func GetEnv(key string) string {
	godotenv.Load("./.env")
	return os.Getenv(key)
}
