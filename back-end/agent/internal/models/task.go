package models

type Task struct {
	ID         uint   `json:"id"`
	Expression string `json:"expression"`
}

type Result struct {
	Task
	Value float64 `json:"value"`
}
