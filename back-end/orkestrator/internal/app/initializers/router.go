package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1"
	apiv1AgentsManager "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1/agents_manager"
	apiv1Status "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1/status"
	apiv1Swagger "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1/swagger"
	apiv1TasksManager "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1/tasks"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/middleware"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// InitializeRouter initializes new gin router
func InitializeRouter(container *dependencies.Container) *gin.Engine {
	r := router.NewRouter()

	initializeMiddlewares(r, container.Config.App)
	v1 := r.Group("/api/v1")

	ctrls := buildControllers(container)
	for i, ctrl := range ctrls {
		ctrlRouterGroup := v1.Group(ctrl.GetRelativePath())
		ctrls[i].DefineRoutes(ctrlRouterGroup)
	}

	return r
}

func initializeMiddlewares(r gin.IRouter, appConfig *config.App) {
	r.Use(middleware.LoggerWithConfig(appConfig.LoggerCfg, log.Logger))
	r.Use((middleware.TokenAuthMiddleware(appConfig.AuthPublicKeyPath)))
	r.Use(middleware.Recovery())
}

func buildControllers(container *dependencies.Container) []apiv1.Controller {
	return []apiv1.Controller{
		apiv1Status.NewController(container.BuildInfo),
		apiv1Swagger.NewController(),
		apiv1TasksManager.NewController(container.TaskManager),
		apiv1AgentsManager.NewController(container.AgentManager),
	}
}
