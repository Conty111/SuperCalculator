package interfaces

import "github.com/Conty111/SuperCalculator/back-end/models"

type Service interface {
	GetTaskByID(taskID uint) (*models.TasksModel, error)
	GetTaskByExpression(expression string) (*models.TasksModel, error)
	CreateTask(expression string) (*models.TasksModel, error)
	DeleteTaskByID(taskID uint) error
	DeleteTaskByExpression(expression string) error
	UpdateTask(taskID uint, param string, value interface{}) error
}
