package helpers

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"time"
)

type TaskResponse struct {
	models.Result
	UserID     uint      `json:"user_id"`
	IsExecuted bool      `json:"is_executed"`
	CreatedAt  time.Time `json:"created_at"`
	ExecutedAt time.Time `json:"executed_at"`
}

func SerializeTasks(userID uint, userTasks []*models.TasksModel) []TaskResponse {
	newTasks := make([]TaskResponse, len(userTasks))
	for i, t := range userTasks {
		newTasks[i] = TaskResponse{
			UserID:     userID,
			IsExecuted: t.IsExecuted,
			CreatedAt:  t.CreatedAt,
			ExecutedAt: t.UpdatedAt,
			Result: models.Result{
				Task:  models.Task{ID: t.ID, Expression: t.Expression},
				Value: t.Value,
				Error: t.Error,
			},
		}
	}
	return newTasks
}
