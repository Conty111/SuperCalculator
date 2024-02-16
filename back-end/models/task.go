package models

import (
	"gorm.io/gorm"
)

type TasksModel struct {
	gorm.Model
	Expression string  `gorm:"primaryKey;type:text;column:expression;index:expression_hash_idx"`
	Value      float64 `gorm:"column:value"`
	IsExecuted bool    `gorm:"column:is_executed"`
	Error      error   `gorm:"column:error"`
}

type Task struct {
	ID         uint   `json:"id"`
	Expression string `json:"expression"`
}

type Result struct {
	Task
	Value float64 `json:"value"`
	Error error   `json:"error"`
}
