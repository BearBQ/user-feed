package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitDB() error {
	dbHost, dbUser, dbPass, dbName, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect postgres")
	}
	log.Printf("подключение к базе postgres %s успешно", dbName)

	_ = db
	return nil
}
