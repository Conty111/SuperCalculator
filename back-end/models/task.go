package models

import (
	"gorm.io/gorm"
)

type TasksModel struct {
	gorm.Model
	Expression string  `gorm:"primaryKey;column:expression;index:expression_hash_idx"`
	Value      float64 `gorm:"column:value"`
	IsExecuted bool    `gorm:"column:is_executed;default:false"`
	Error      string  `gorm:"column:error;type:text"`
}

type Task struct {
	ID         uint   `json:"id"`
	Expression string `json:"expression" binding:"required"`
}

type Result struct {
	Task
	Value float64 `json:"value"`
	Error string  `json:"error"`
}
