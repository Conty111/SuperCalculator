package agents_manager

import "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"

type WorkersListResponse[T any] struct {
	Responses []*helpers.AgentResponse[T] `json:"responses"`
}
