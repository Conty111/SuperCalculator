package services

import (
	"bytes"
	"encoding/json"
	"errors"
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
	timeout      time.Duration
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
		timeout:      time.Second * 5,
	}
}

// GetAllTasks returns all tasks in database
func (tm *TaskManager) GetAllTasks() ([]*models.TasksModel, error) {
	log.Info().Msg("going to db to get all tasks")
	return tm.Repo.GetAllTasks()
}

// SetCalculationSettings sends settings to agents in parallel and returns their responses
func (tm *TaskManager) SetCalculationSettings(settings *models.CalculationSettings) ([]map[string]interface{}, []int) {
	if settings.Timeout != 0 {
		tm.timeout = time.Duration(settings.Timeout) * time.Second
	}

	client := http.Client{Timeout: tm.timeout}
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
			body, status, err := sendRequestToAgent(&client, bytes.NewReader(reqBody), agentAddr, http.MethodPut)
			if err != nil {
				statuses[i] = http.StatusInternalServerError
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
	client := http.Client{Timeout: tm.timeout}
	wg := sync.WaitGroup{}
	wg.Add(int(tm.AgentCount))

	responseBodys := make([]map[string]interface{}, int(tm.AgentCount))
	statuses := make([]int, int(tm.AgentCount))

	for i, agentAddr := range tm.AgentAddress {
		agentAddr := agentAddr
		i := i
		go func() {
			defer wg.Done()
			body, status, err := sendRequestToAgent(&client, nil, agentAddr, http.MethodGet)
			if err != nil {
				statuses[i] = http.StatusInternalServerError
			}
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

// sendRequestToAgent sends request and returns body and status
func sendRequestToAgent(
	client *http.Client,
	reqBody io.Reader,
	agentAddr string,
	method string) (map[string]interface{}, int, error) {
	body := make(map[string]interface{})
	req, err := http.NewRequest(method, agentAddr, reqBody)
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
