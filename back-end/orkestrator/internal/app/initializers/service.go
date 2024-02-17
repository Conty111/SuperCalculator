package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/services"
)

func InitializeService(container *dependencies.Container) interfaces.Service {
	rep := repository.NewTasksRepository(container.Database)
	return services.NewTaskManager(
		rep,
		container.Producer.InChan,
		container.Config.App.AgentCount,
		container.Config.HTTPConfig.AgentAddresses,
	)
}
