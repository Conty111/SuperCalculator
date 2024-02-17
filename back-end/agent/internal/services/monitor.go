package services

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"sync"
)

type Monitor struct {
	Lock            *sync.RWMutex
	AgentID         int32
	EmployedWorkers uint
	FreeWorkers     uint
	CompletedTasks  uint
	LastTaskID      uint
}

func NewMonitor(agentID int32) *Monitor {
	return &Monitor{
		AgentID: agentID,
		Lock:    &sync.RWMutex{},
	}
}

func (m *Monitor) GetStats() *models.Stats {
	return &models.Stats{
		//EmployedWorkers: m.EmployedWorkers,
		//FreeWorkers:     m.FreeWorkers,
		CompletedTasks: m.CompletedTasks,
		LastTaskID:     m.LastTaskID,
	}
}

func (m *Monitor) CompleteWork(taskID uint) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	m.CompletedTasks++
	m.FreeWorkers++
	m.EmployedWorkers--
	m.LastTaskID = taskID
}

func (m *Monitor) AddWork() {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	m.FreeWorkers--
	m.EmployedWorkers++
}
