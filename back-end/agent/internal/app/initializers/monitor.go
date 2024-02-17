package initializers

import "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"

func InitializeMonitor(agentID int32) *services.Monitor {
	return services.NewMonitor(agentID)
}
