package models

type Stats struct {
	EmployedWorkers uint `json:"employed_workers"`
	FreeWorkers     uint `json:"free_workers"`
	CompletedTasks  uint `json:"completed_tasks"`
	LastTaskID      uint `json:"last_task_id"`
}
