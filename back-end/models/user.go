package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"constraint:unique_index" binding:"required,email"`
	Password string `json:"password" gorm:"column:password" binding:"required"`
	Username string `json:"username" gorm:"column:username" binding:"required"`
	Role     Role   `json:"role" gorm:"column:role"`
}

type Token struct {
	Expires time.Time `json:"expires"` // Time when token expires
	UserID  uint      `json:"userID"`
}
