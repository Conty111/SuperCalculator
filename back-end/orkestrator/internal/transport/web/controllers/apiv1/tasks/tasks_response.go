package tasks

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
)

// Response is a declaration for a status response
type Response struct {
	Status  string `jsonapi:"attr,status"`
	Message string `json:"message"`
}

// TasksListResponse is a declaration of response to GetTasks endpoint
type TasksListResponse struct {
	Status string                  `jsonapi:"attr,status"`
	Tasks  []*helpers.TaskResponse `json:"tasks"`
}
