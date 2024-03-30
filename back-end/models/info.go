package models

type AgentInfo struct {
	Name    string `json:"name"`
	AgentID int32  `json:"agent_id"`
	//EmployedWorkers uint `json:"employed_workers"`
	//FreeWorkers     uint `json:"free_workers"`
	CompletedTasks uint `json:"completed_tasks"`
	LastTaskID     uint `json:"last_task_id"`
}
