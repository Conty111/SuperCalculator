package repository

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TasksRepository struct {
	Database *gorm.DB
}

func NewTasksRepository(db *gorm.DB) *TasksRepository {
	return &TasksRepository{
		Database: db,
	}
}

func (tr *TasksRepository) GetTaskByExpression(expression string) (*models.TasksModel, error) {
	task := models.TasksModel{
		Expression: expression,
	}
	r := tr.Database.Model(models.TasksModel{}).Find(&task)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected == 0 {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (tr *TasksRepository) GetTaskByID(taskID string) (*models.TasksModel, error) {
	var task models.TasksModel
	r := tr.Database.Model(models.TasksModel{}).Where("id = ?", taskID).Find(&task)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected == 0 {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (tr *TasksRepository) CreateTask(task *models.TasksModel) error {
	return tr.Database.Model(models.TasksModel{}).Create(task).Error
}

func (tr *TasksRepository) UpdateTask(task *models.TasksModel, param string, value any) error {
	return tr.Database.Model(task).Update(param, value).Error
}
