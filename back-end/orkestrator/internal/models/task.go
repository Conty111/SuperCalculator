package models

type Task struct {
	ID         uint   `json:"id"`
	Expression string `json:"expression"`
}
