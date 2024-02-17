package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"sync"
	"time"
)

type TaskManager struct {
	Repo         *repository.TasksRepository
	ProduceChan  chan<- models.Task
	AgentCount   uint
	AgentAddress []string
	timeoutResp  time.Duration
	cachedTasks  map[uint]interface{}
	lock         sync.RWMutex
}

func NewTaskManager(rep *repository.TasksRepository,
	produceCg chan<- models.Task,
	agentCount uint,
	agentAddrs []string) *TaskManager {
	return &TaskManager{
		Repo:         rep,
		ProduceChan:  produceCg,
		AgentCount:   agentCount,
		AgentAddress: agentAddrs,
		timeoutResp:  time.Second * 5,
		cachedTasks:  make(map[uint]interface{}),
		lock:         sync.RWMutex{},
	}
}

// GetAllTasks returns all tasks in database
func (tm *TaskManager) GetAllTasks() ([]*models.TasksModel, error) {
	return tm.Repo.GetAllTasks()
}

// SetCalculationSettings sends settings to agents in parallel and returns their responses
func (tm *TaskManager) SetCalculationSettings(settings *models.CalculationSettings) ([]map[string]interface{}, []int) {
	if settings.Timeout != 0 {
		tm.timeoutResp = time.Duration(settings.Timeout) * time.Second
	}

	client := http.Client{Timeout: tm.timeoutResp}
	wg := sync.WaitGroup{}
	wg.Add(int(tm.AgentCount))

	responseBodys := make([]map[string]interface{}, int(tm.AgentCount))
	statuses := make([]int, int(tm.AgentCount))

	for i, agentAddr := range tm.AgentAddress {
		agentAddr := agentAddr
		i := i
		go func() {
			defer wg.Done()
			reqBody, err := json.Marshal(settings.DurationSettings)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal duration settings")
				statuses[i] = http.StatusInternalServerError
				return
			}
			body, status, err := sendRequestToAgent(
				&client,
				bytes.NewReader(reqBody),
				fmt.Sprintf("%s/calculator", agentAddr),
				http.MethodPut,
			)
			if err != nil {
				statuses[i] = http.StatusInternalServerError
				return
			}
			responseBodys[i] = body
			statuses[i] = status
		}()
	}
	wg.Wait()

	return responseBodys, statuses
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

// GetWorkersInfo gets info from workers in parallel and return their bodys and statuses
func (tm *TaskManager) GetWorkersInfo() ([]map[string]interface{}, []int) {
	client := http.Client{Timeout: tm.timeoutResp}
	wg := sync.WaitGroup{}
	wg.Add(int(tm.AgentCount))

	responseBodys := make([]map[string]interface{}, int(tm.AgentCount))
	statuses := make([]int, len(tm.AgentAddress))

	for i, agentAddr := range tm.AgentAddress {
		agentAddr := agentAddr
		i := i
		go func() {
			defer wg.Done()
			body, status, err := sendRequestToAgent(
				&client,
				nil,
				fmt.Sprintf("%s/status", agentAddr),
				http.MethodGet,
			)
			if err != nil {
				statuses[i] = http.StatusInternalServerError
				return
			}
			log.Print(body, statuses, err)
			responseBodys[i] = body
			statuses[i] = status
		}()
	}
	wg.Wait()

	return responseBodys, statuses
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
	fields := map[string]interface{}{
		"value":       res.Value,
		"is_executed": true,
		"error":       res.Error,
	}
	return tm.Repo.Update(t, fields)
}

func (tm *TaskManager) GetCachedTasksIDs() []uint {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	keys := make([]uint, len(tm.cachedTasks))
	var i int
	for k, _ := range tm.cachedTasks {
		keys[i] = k
	}
	return keys
}

// sendRequestToAgent sends request and returns body and status
func sendRequestToAgent(
	client *http.Client,
	reqBody io.Reader,
	agentAddr string,
	method string) (map[string]interface{}, int, error) {
	body := make(map[string]interface{})
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s", agentAddr), reqBody)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return nil, 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to send request")
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return nil, 0, err
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal body data")
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}
