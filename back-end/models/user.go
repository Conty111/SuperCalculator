package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"constraint:unique_index" binding:"required,email"`
	Password string `json:"password" gorm:"column:password" binding:"required"`
	Name     string `json:"name" gorm:"column:name"`
}
