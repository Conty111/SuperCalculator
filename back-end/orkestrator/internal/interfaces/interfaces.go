package interfaces

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
)

type Service interface {
	GetAllTasks() ([]*models.TasksModel, error)
	SetCalculationSettings(settings *models.CalculationSettings) ([]map[string]interface{}, []int)
	CreateTask(expression string) (*models.TasksModel, error)
	GetWorkersInfo() ([]map[string]interface{}, []int)
	DeleteTaskByID(taskID uint) error
	DeleteTaskByExpression(expression string) error
	SaveResult(res *models.Result) error
}
