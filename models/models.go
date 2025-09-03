package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name  string `json:"name"`
	Posts []Post `json:"posts" gorm:"foreignKey:UserID;references:ID"`
}

type Post struct {
	gorm.Model
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
}

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type CreatePostRequest struct {
	Content string `json:"content" validate:"required"`
	USerID  uint   `json:"user_id" validate:"required"`
}
