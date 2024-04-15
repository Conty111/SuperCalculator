package interfaces

import (
	"context"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
)

type AgentManager interface {
	SetSettings(settings *models.Settings) []*helpers.AgentResponse
	GetWorkersInfo() []*helpers.AgentResponse
}

type AgentAPIClient interface {
	GetAgentInfo()
	SetSettings(settings *models.Settings)
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
