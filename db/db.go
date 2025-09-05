package db

import (
	"fmt"
	"log"
	"os"
	"user-feed/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBase interface {
	CreateUser(user *models.Users) error
	CreatePost(post *models.Post) error
	GetPosts(userID uint) (models.Users, error)
	CheckUserExist(email string) (bool, error)
	Close() error
	AutoMigrate() error
}

type DataBasePostgres struct {
	DB *gorm.DB
}

func (db *DataBasePostgres) CreateUser(user *models.Users) error {
	result := db.DB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (db *DataBasePostgres) CheckUserExist(email string) (bool, error) {
	var user models.Users
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil // Пользователь не найден
		}
		return false, fmt.Errorf("failed to check user existence: %w", result.Error)
	}

	return true, nil // Пользователь найден
}

func (db *DataBasePostgres) CreatePost(post *models.Post) error {
	result := db.DB.Create(post)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (db *DataBasePostgres) GetPosts(userID uint) (models.Users, error) {
	var user models.Users
	result := db.DB.Preload("Posts").First(&user, userID)
	if result.Error != nil {
		return models.Users{}, fmt.Errorf("failed to get user posts: %w", result.Error)
	}
	return user, nil
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
