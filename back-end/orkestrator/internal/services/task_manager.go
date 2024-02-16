package services

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
)

type TaskManager struct {
	Repo *repository.TasksRepository
}

func NewTaskManager(rep *repository.TasksRepository) *TaskManager {
	return &TaskManager{
		Repo: rep,
	}
}

func (tm *TaskManager) GetTaskByID(taskID uint) (*models.TasksModel, error) {
	var task models.TasksModel
	task.ID = taskID
	return tm.GetTaskByID(taskID)
}

func (tm *TaskManager) GetTaskByExpression(taskID uint) (*models.TasksModel, error) {
	var task models.TasksModel
	task.ID = taskID
	return tm.GetTaskByID(taskID)
}
