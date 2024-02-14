package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/initializers"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/models"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/kafka-broker"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Application is a main struct for the application that contains general information
type Application struct {
	httpServer *http.Server
	consumer   *kafka_broker.AppConsumer
	Container  *dependencies.Container
}

// InitializeApplication initializes new application
func InitializeApplication() (*Application, error) {
	initializers.InitializeEnvs()

	if err := initializers.InitializeLogs(); err != nil {
		return nil, err
	}

	return BuildApplication()
}

func BuildApplication() (*Application, error) {
	cfg := config.GetConfig()
	info := initializers.InitializeBuildInfo()
	monitor := initializers.InitializeMonitor()
	svc := initializers.InitializeExpressionService()
	container := &dependencies.Container{
		BuildInfo:     info,
		Config:        cfg,
		Monitor:       monitor,
		ExpressionSvc: svc,
	}

	consumer := initializers.InitializeConsumer(container)
	router := initializers.InitializeRouter(container)
	server := initializers.InitializeHTTPServer(router, &cfg.HTTPConfig)

	return &Application{
		httpServer: server,
		consumer:   consumer,
		Container:  container,
	}, nil
}

// Start starts application services
func (a *Application) Start(ctx context.Context, cli bool) {
	if cli {
		return
	}

	a.startHTTPServer()
	a.startProducer(a.startConsumer())
}

// Stop stops application services
func (a *Application) Stop() (err error) {
	a.consumer.Stop()
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

func (a *Application) startConsumer() <-chan models.Result {
	return a.consumer.Start()
}

func (a *Application) startProducer(msgChannel <-chan models.Result) {
	go func() {
		for res := range msgChannel {
			data, err := json.Marshal(res)
			log.Error().Err(err).Msg("Error while trying to marshal to json")
			log.Debug().Str("task", string(data)).Msg("Executed")
		}
	}()
}
