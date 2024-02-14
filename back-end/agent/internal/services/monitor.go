package services

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/models"
	"sync"
)

type Monitor struct {
	Lock            *sync.RWMutex
	EmployedWorkers uint
	FreeWorkers     uint
	CompletedTasks  uint
}

func NewMonitor() *Monitor {
	return &Monitor{
		Lock: &sync.RWMutex{},
	}
}

func (m *Monitor) GetStats() *models.Stats {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return &models.Stats{
		EmployedWorkers: m.EmployedWorkers,
		FreeWorkers:     m.FreeWorkers,
		CompletedTasks:  m.CompletedTasks,
	}
}

func (m *Monitor) CompleteWork() {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	m.CompletedTasks++
	m.FreeWorkers++
	m.EmployedWorkers--
}

func (m *Monitor) AddWork() {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	m.FreeWorkers--
	m.EmployedWorkers++
}
