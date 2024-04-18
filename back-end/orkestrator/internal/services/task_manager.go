package services

import (
	"context"
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type TaskManager struct {
	Repo        *repository.TasksRepository
	ProduceChan chan<- models.Task
	Agents      []models.AgentConfig
	timeRetry   time.Duration
	cachedTasks map[uint]interface{}
	lock        sync.RWMutex
}

func NewTaskManager(rep *repository.TasksRepository,
	produceCg chan<- models.Task,
	agents []models.AgentConfig,
	timeRetry time.Duration) *TaskManager {
	return &TaskManager{
		Repo:        rep,
		ProduceChan: produceCg,
		Agents:      agents,
		timeRetry:   timeRetry,
		cachedTasks: make(map[uint]interface{}),
		lock:        sync.RWMutex{},
	}
}

// Start initing service (sends not executed tasks from db and enbales retrying)
func (tm *TaskManager) Start(ctx context.Context) {
	tasks, err := tm.Repo.GetNotExecutedTasks()
	if err != nil {
		log.Error().Err(err).Msg("error while trying to get not executed tasks from database")
	}
	for _, t := range tasks {
		tm.cachedTasks[t.ID] = t
	}
	tm.EnableRetrying(ctx)
}

// GetAllTasks returns all tasks in database
func (tm *TaskManager) GetAllTasks() ([]*models.TasksModel, error) {
	return tm.Repo.GetAllTasks()
}

// CreateTask creates task in database
func (tm *TaskManager) CreateTask(expression string) (*models.TasksModel, error) {
	task := models.TasksModel{Expression: expression}
	t, err := tm.Repo.GetByExpression(expression)
	if errors.Is(err, clierrs.ErrTaskNotFound) {
		err = tm.Repo.Create(&task)
		if err != nil {
			return nil, err
		}
		tm.lock.RLock()
		tm.cachedTasks[task.ID] = expression
		tm.lock.RUnlock()
		msg := models.Task{
			ID:         task.ID,
			Expression: expression,
		}
		tm.ProduceChan <- msg
		return &task, err
	}
	if err == nil {
		return t, clierrs.ErrTaskAlreadyCreated
	}
	return nil, err
}

// DeleteTaskByID deletes task by id
func (tm *TaskManager) DeleteTaskByID(taskID uint) error {
	t := &models.TasksModel{}
	t.ID = taskID
	return tm.Repo.Delete(t)
}

// DeleteTaskByExpression deletes task by expression
func (tm *TaskManager) DeleteTaskByExpression(expression string) error {
	t := &models.TasksModel{}
	t.Expression = expression
	return tm.Repo.Delete(t)
}

// SaveResult saves (updates) result in database
func (tm *TaskManager) SaveResult(res *models.Result) error {
	tm.lock.RLock()
	delete(tm.cachedTasks, res.ID)
	tm.lock.RUnlock()
	t, err := tm.Repo.GetByID(res.ID)
	if err != nil {
		return err
	}
	if t.IsExecuted {
		return nil
	}
	fields := map[string]interface{}{
		"value":       res.Value,
		"is_executed": true,
	}
	if res.Error != "" {
		fields["error"] = res.Error
	}
	return tm.Repo.Update(t, fields)
}

func (tm *TaskManager) EnableRetrying(ctx context.Context) {
	go func() {
		for {
			tm.lock.Lock()
			t := tm.timeRetry
			tm.lock.Unlock()
			select {
			case <-ctx.Done():
				return
			default:
				tm.RetryTasks()
			}
			time.Sleep(t)
		}
	}()
}

// getCachedTasksIDs returns slice of current not completed tasks
func (tm *TaskManager) getCachedTasksIDs() []uint {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	keys := make([]uint, len(tm.cachedTasks))
	var i int
	for k, _ := range tm.cachedTasks {
		keys[i] = k
	}
	return keys
}

// RetryTasks sends tasks in queue again if they aren't updated yet
func (tm *TaskManager) RetryTasks() {
	ids := tm.getCachedTasksIDs()
	if len(ids) > 0 {
		log.Info().Any("taskIDs", ids).Msg("retrying tasks")
		for _, id := range ids {
			t, err := tm.Repo.GetByID(id)
			if err != nil {
				log.Error().Err(err).Msg("failed get task by id")
				continue
			}
			task := models.Task{
				ID:         id,
				Expression: t.Expression,
			}
			tm.ProduceChan <- task
		}
	}
}
