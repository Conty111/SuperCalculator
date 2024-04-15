package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/auth"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/services"
	"github.com/rs/zerolog/log"
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
func InitializeUserManager(container *dependencies.Container) interfaces.UserManager {
	return repository.NewUserRepository(container.Database)
}

func InitializeAuthManager(container *dependencies.Container) interfaces.AuthManager {
	authManager, err := auth.NewAuth(container.Config.App)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize auth manager")
		return nil
	}
	return authManager
}
