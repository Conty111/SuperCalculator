package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/services"
)

func InitializeTaskManager(container *dependencies.Container) interfaces.TaskManager {
	rep := repository.NewTasksRepository(container.Database)

	return services.NewTaskManager(
		rep,
		container.Producer.InChan,
		container.Config.App.Agents,
		container.Config.App.TimeoutResponse,
		container.Config.App.TimeToRetry,
	)
}

func InitializeAgentManager(container *dependencies.Container) interfaces.AgentManager {
	return services.NewAgentManager(
		container.Config.App.ApiToUse,
		container.Config.App.Agents,
		container.Config.App.TimeoutResponse,
	)
}
