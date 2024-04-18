package services

import (
	"context"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type TaskManager struct {
	TasksRepository *repository.TasksRepository
	UserRepo        interfaces.UserManager
	ProduceChan     chan<- models.Task
	Agents          []models.AgentConfig
	timeRetry       time.Duration
	cachedTasks     map[uint]interface{}
	lock            sync.RWMutex
}

func NewTaskManager(
	tasksRep *repository.TasksRepository,
	userRep interfaces.UserManager,
	produceCg chan<- models.Task,
	agents []models.AgentConfig,
	timeRetry time.Duration) *TaskManager {

	return &TaskManager{
		TasksRepository: tasksRep,
		UserRepo:        userRep,
		ProduceChan:     produceCg,
		Agents:          agents,
		timeRetry:       timeRetry,
		cachedTasks:     make(map[uint]interface{}),
		lock:            sync.RWMutex{},
	}
}

// Start starts task manager
func (tm *TaskManager) Start(ctx context.Context) {
	tm.getNotExecutedTasks()
	tm.EnableRetrying(ctx)
}

// getNotExecutedTasks gets not executed tasks from database and saves it in manager cache
func (tm *TaskManager) getNotExecutedTasks() {
	tasks, err := tm.TasksRepository.GetNotExecutedTasks()
	if err != nil {
		log.Error().Err(err).Msg("error while trying to get not executed tasks from database")
	}
	for _, t := range tasks {
		tm.cachedTasks[t.ID] = t
	}
}

// GetAllTasks returns all tasks in database
func (tm *TaskManager) GetAllTasks() ([]*models.TasksModel, error) {
	return tm.TasksRepository.GetAllTasks()
}

// CreateTask creates task in database
func (tm *TaskManager) CreateTask(expression string, callerID uint) (*models.TasksModel, error) {
	user, err := tm.UserRepo.GetUserByID(callerID)
	if err != nil {
		return nil, err
	}
	task := models.TasksModel{Expression: expression, User: user}
	err = tm.TasksRepository.Create(&task)
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

// DeleteTaskByID deletes task by id
func (tm *TaskManager) DeleteTaskByID(taskID uint) error {
	t := &models.TasksModel{}
	t.ID = taskID
	return tm.TasksRepository.Delete(t)
}

// DeleteTaskByExpression deletes task by expression
func (tm *TaskManager) DeleteTaskByExpression(expression string) error {
	t := &models.TasksModel{}
	t.Expression = expression
	return tm.TasksRepository.Delete(t)
}

// SaveResult saves (updates) task in database if it didn't be executed
func (tm *TaskManager) SaveResult(res *models.Result) error {
	tm.lock.RLock()
	delete(tm.cachedTasks, res.ID)
	tm.lock.RUnlock()
	t, err := tm.TasksRepository.GetByID(res.ID)
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
	return tm.TasksRepository.Update(t, fields)
}

// EnableRetrying starts goroutine that periodically retries cached tasks
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
			t, err := tm.TasksRepository.GetByID(id)
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
