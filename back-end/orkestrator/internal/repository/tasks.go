package repository

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TasksRepository struct {
	Database *gorm.DB
}

func NewTasksRepository(db *gorm.DB) *TasksRepository {
	return &TasksRepository{
		Database: db,
	}
}

func (tr *TasksRepository) GetNotExecutedTasks() ([]*models.TasksModel, error) {
	var tasks []*models.TasksModel
	r := tr.Database.Model(models.TasksModel{}).Where("is_executed = ?", false).Find(&tasks)
	return tasks, r.Error
}

func (tr *TasksRepository) GetAllTasks() ([]*models.TasksModel, error) {
	var tasks []*models.TasksModel
	r := tr.Database.
		Model(models.TasksModel{}).
		Preload("User").
		Find(&tasks)
	log.Print(tasks[0].User)
	return tasks, r.Error
}

func (tr *TasksRepository) GetByExpression(expression string) (*models.TasksModel, error) {
	task := models.TasksModel{
		Expression: expression,
	}
	r := tr.Database.
		Model(models.TasksModel{}).
		Preload(clause.Associations).
		Find(&task)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected == 0 {
		return nil, clierrs.ErrTaskNotFound
	}
	return &task, nil
}

func (tr *TasksRepository) GetByID(taskID uint) (*models.TasksModel, error) {
	var task models.TasksModel
	r := tr.Database.
		Model(models.TasksModel{}).
		Where("id = ?", taskID).
		Preload(clause.Associations).
		Find(&task)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected == 0 {
		return nil, clierrs.ErrTaskNotFound
	}
	return &task, nil
}

func (tr *TasksRepository) Create(task *models.TasksModel) error {
	return tr.Database.Model(models.TasksModel{}).Create(task).Error
}

func (tr *TasksRepository) Update(task *models.TasksModel, fields map[string]interface{}) error {
	return tr.Database.Model(task).Updates(fields).Error
}

func (tr *TasksRepository) Delete(task *models.TasksModel) error {
	return tr.Database.Model(task).Delete(task).Error
}
