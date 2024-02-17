package manager

import "github.com/Conty111/SuperCalculator/back-end/models"

// Response is a declaration for a status response
type Response struct {
	Status  string `jsonapi:"attr,status"`
	Message string `json:"message"`
}

type Task struct {
	models.Result
	IsExecuted bool `json:"is_executed"`
}

// TasksListResponse is a declaration of response to GetTasks endpoint
type TasksListResponse struct {
	Status string  `jsonapi:"attr,status"`
	Tasks  []*Task `json:"tasks"`
}

type WorkerResponse struct {
	Status   string
	Response map[string]interface{}
}

type WorkersListResponse struct {
	Status    string           `json:"status"`
	Responses []WorkerResponse `json:"responses"`
}