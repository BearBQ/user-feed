package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Posts []Post `json:"posts"`
}

type Post struct {
	gorm.Model
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
}
