package db

import (
	"fmt"
	"log"
	"os"
	"user-feed/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBasePostgres struct {
	DB *gorm.DB
}

// Close() закрывает соединение с постгрес
func (db *DataBasePostgres) Close() error {
	sqlBase, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlBase.Close()
}

// Automigrate - обновляет таблицы Users и POST в Postgres
func (db *DataBasePostgres) AutoMigrate() error {
	err := db.DB.AutoMigrate(&models.Users{}, &models.Post{})
	if err != nil {
		return fmt.Errorf("failed to migrate data: %w", err)
	}
	log.Println("Automigrate done")
	return nil
}

//CloseConnection

// NewDataBase() - создает подключение к базе данных Postgres
func NewDataBase() (*DataBasePostgres, error) {
	dbHost, dbUser, dbPass, dbName, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	log.Printf("подключение к базе postgres %s успешно", dbName)

	return &DataBasePostgres{DB: db}, nil
}
