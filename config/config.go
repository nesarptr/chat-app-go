package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	HOST := GetEnv("HOST")
	USER := GetEnv("USER")
	PASSWORD := GetEnv("PASSWORD")
	DATABASE := GetEnv("DATABASE")
	PORT := GetEnv("PGPORT")
	SSLMODE := GetEnv("SSLMODE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", HOST, USER, PASSWORD, DATABASE, PORT, SSLMODE)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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