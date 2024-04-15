package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"constraint:unique_index" binding:"required,email"`
	Password string `json:"password" gorm:"column:password" binding:"required"`
	Name     string `json:"name" gorm:"column:name"`
}

type Token struct {
	Expires time.Time `json:"expires"` // Time when token expires
	UserID  string    `json:"userID"`
}