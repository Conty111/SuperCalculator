package app

import (
	"context"
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"net/http"

	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/initializers"
	"github.com/rs/zerolog/log"
)

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

	return BuildApplication(ctx)
}

func BuildApplication(ctx context.Context) (*Application, error) {
	cfg := config.GetConfig(ctx)
	info := initializers.InitializeBuildInfo()
	db := initializers.InitializeDatabase(cfg.DB.DSN)

	container := &dependencies.Container{
		BuildInfo: info,
		Config:    cfg,
		Database:  db,
	}
	producer := initializers.InitializeProducer(container)
	container.Producer = producer
	svc := initializers.InitializeService(container)
	container.Service = svc
	consumer := initializers.InitializeConsumer(container)
	container.Consumer = consumer
	router := initializers.InitializeRouter(container)
	server := initializers.InitializeHTTPServer(router, cfg.HTTPConfig)

	return &Application{
		httpServer: server,
		Container:  container,
	}, nil
}

// Start starts application services
func (a *Application) Start(ctx context.Context, cli bool) {
	if cli {
		return
	}
	a.Container.Consumer.Start()
	a.Container.Producer.Start()
	a.Container.Service.Start(ctx)
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
