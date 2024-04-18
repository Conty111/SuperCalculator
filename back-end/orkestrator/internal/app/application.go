package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/initializers"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/repository"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var attempt uint

const maxAttemptCount = 3

// Application is a main struct for the application that contains general information
type Application struct {
	httpServer *http.Server
	Container  *dependencies.Container
}

// InitializeApplication initializes new application
func InitializeApplication(ctx context.Context) (*Application, error) {
	initializers.InitializeEnvs()

	if err := initializers.InitializeLogs(); err != nil {
		return nil, err
	}

	return BuildApplication(), nil
}

func BuildApplication() *Application {
	defer func() {
		if r := recover(); r != nil {
			attempt++
			if attempt <= maxAttemptCount {
				log.Error().Uint("Attempt", attempt).Msg("Не удалось запустить приложение, пробуем еще раз")
				for i := 3; i > 0; i-- {
					log.Info().Msg(fmt.Sprintf("Перезапуск через %d...", i))
					time.Sleep(time.Second)
				}
				BuildApplication()
			}
			log.Fatal().Msg("Не удалось запустить приложение")
		}
	}()

	cfg := config.GetConfig()
	info := initializers.InitializeBuildInfo()
	db := initializers.InitializeDatabase(cfg.DB.DSN, cfg.DB.DBtype)

	container := &dependencies.Container{
		BuildInfo: info,
		Config:    cfg,
		Database:  db,
	}

	container.AgentManager = initializers.InitializeAgentManager(container)
	container.UserManager = repository.NewUserRepository(db)
	container.AuthManager = initializers.InitializeAuthManager(container)

	container.Producer = initializers.InitializeProducer(container)
	container.TaskManager = initializers.InitializeTaskManager(container)
	container.Consumer = initializers.InitializeConsumer(container)

	router := initializers.InitializeRouter(container)
	server := initializers.InitializeHTTPServer(router, cfg.HTTPConfig)

	return &Application{
		httpServer: server,
		Container:  container,
	}
}

// Start starts application services
func (a *Application) Start(ctx context.Context, cli bool) {
	if cli {
		return
	}
	a.Container.TaskManager.Start(ctx)
	a.Container.Consumer.Start()
	a.Container.Producer.Start()
	a.startHTTPServer()
}

// Stop stops application services
func (a *Application) Stop() (err error) {
	a.Container.Consumer.Stop()
	a.Container.Producer.Stop()
	return a.httpServer.Shutdown(context.TODO())
}

func (a *Application) startHTTPServer() {
	go func() {
		log.Info().Str("HTTPServerAddress", a.httpServer.Addr).Msg("started http server")

		// service connections
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic().Err(err).Msg("HTTP Server stopped")
		}
	}()
}
