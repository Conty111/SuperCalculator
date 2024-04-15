package agents_manager

import "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"

type WorkersListResponse struct {
	Responses []*helpers.AgentResponse `json:"responses"`
}
