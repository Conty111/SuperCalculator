package initializers

import "github.com/Conty111/SuperCalculator/back-end/agent/internal/services"

func InitializeMonitor() *services.Monitor {
	return services.NewMonitor()
}
