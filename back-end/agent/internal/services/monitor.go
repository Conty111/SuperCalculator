package services

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"sync"
)

type Monitor struct {
	Lock            *sync.RWMutex
	AgentID         int32
	Name            string
	EmployedWorkers uint
	FreeWorkers     uint
	CompletedTasks  uint
	CurrentTasks    map[uint]interface{}
	LastTaskID      uint
}

func NewMonitor(agentID int32, name string) *Monitor {
	return &Monitor{
		AgentID:      agentID,
		Name:         name,
		Lock:         &sync.RWMutex{},
		CurrentTasks: make(map[uint]interface{}),
	}
}

func (m *Monitor) GetInfo() *models.AgentInfo {
	return &models.AgentInfo{
		Name:           m.Name,
		AgentID:        m.AgentID,
		CompletedTasks: m.CompletedTasks,
		LastTaskID:     m.LastTaskID,
	}
}

func (m *Monitor) CompleteTask(taskID uint) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.CurrentTasks[taskID]; !ok {
		return agent_errors.ErrTaskNotFound
	}
	delete(m.CurrentTasks, taskID)
	m.CompletedTasks++
	m.FreeWorkers++
	m.EmployedWorkers--
	m.LastTaskID = taskID
	return nil
}

func (m *Monitor) AddTask(taskID uint) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.CurrentTasks[taskID] = "task"
	m.FreeWorkers--
	m.EmployedWorkers++
}
