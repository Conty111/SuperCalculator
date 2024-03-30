package interfaces

import (
	"context"
	"github.com/Conty111/SuperCalculator/back-end/models"
)

type AgentManager interface {
	SetSettings(settings *models.Settings) ([]map[string]interface{}, []int)
	GetWorkersInfo() ([]map[string]interface{}, []int)
}

type TaskManager interface {
	GetAllTasks() ([]*models.TasksModel, error)
	CreateTask(expression string) (*models.TasksModel, error)
	DeleteTaskByID(taskID uint) error
	DeleteTaskByExpression(expression string) error
	SaveResult(res *models.Result) error
	Start(ctx context.Context)
}

type AuthService interface {
	Login(email, password string)
	Register()
}
