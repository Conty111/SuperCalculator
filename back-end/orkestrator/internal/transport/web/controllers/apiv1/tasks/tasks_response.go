package tasks

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"time"
)

// Response is a declaration for a status response
type Response struct {
	Status  string `jsonapi:"attr,status"`
	Message string `json:"message"`
}

type Task struct {
	models.Result
	IsExecuted bool      `json:"is_executed"`
	CreatedAt  time.Time `json:"created_at"`
	ExecutedAt time.Time `json:"executed_at"`
}

// TasksListResponse is a declaration of response to GetTasks endpoint
type TasksListResponse struct {
	Status string  `jsonapi:"attr,status"`
	Tasks  []*Task `json:"tasks"`
}
