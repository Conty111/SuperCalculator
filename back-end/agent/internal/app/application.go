package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/dependencies"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/initializers"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/config"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/kafka-broker"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

var attempt uint

const maxAttemptCount = 3

// Application is a main struct for the application that contains general information
type Application struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	consumer   *kafka_broker.AppConsumer
	producer   *kafka_broker.AppProducer
	Container  *dependencies.Container
}

// InitializeApplication initializes new application
func InitializeApplication(ctx context.Context) (*Application, error) {
	initializers.InitializeEnvs()

	if err := initializers.InitializeLogs(); err != nil {
		return nil, err
	}

	return BuildApplication(ctx), nil
}

func BuildApplication(ctx context.Context) *Application {
	defer func() {
		if r := recover(); r != nil {
			attempt++
			if attempt <= maxAttemptCount {
				log.Error().Uint("Attempt", attempt).Msg("Не удалось запустить приложение, пробуем еще раз")
				for i := 3; i > 0; i-- {
					log.Info().Msg(fmt.Sprintf("Перезапуск через %d...", i))
					time.Sleep(time.Second)
				}
				BuildApplication(ctx)
			}
			log.Fatal().Msg("Не удалось запустить приложение")
		}
	}()

	cfg := config.GetConfig(ctx)

	info := initializers.InitializeBuildInfo()
	monitor := initializers.InitializeMonitor(cfg.BrokerCfg.Partition, cfg.App.Name)
	svc := initializers.InitializeExpressionService()
	container := &dependencies.Container{
		BuildInfo:  info,
		Config:     cfg,
		Monitor:    monitor,
		Calculator: svc,
	}

	consumer := initializers.InitializeConsumer(container)
	producer := initializers.InitializeProducer(container)
	router := initializers.InitializeRouter(container)
	httpServer := initializers.InitializeHTTPServer(router, &cfg.HTTPConfig)
	grpcServer := initializers.InitializeGRPCServer(container)

	return &Application{
		httpServer: httpServer,
		grpcServer: grpcServer,
		consumer:   consumer,
		producer:   producer,
		Container:  container,
	}
}

// Start starts application services
func (a *Application) Start(ctx context.Context, cli bool) {
	if cli {
		return
	}
	a.startHTTPServer()
	a.startGRPCServer()
	a.startProducer(a.startConsumer())
}

// Stop stops application services
func (a *Application) Stop() (err error) {
	log.Info().Msg("gracefully stopping")
	a.consumer.Stop()
	a.producer.Stop()
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

func (a *Application) startGRPCServer() {
	go func() {
		addr := a.Container.Config.GRPCConfig.Host + ":" + a.Container.Config.GRPCConfig.Port

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Panic().Err(err).Msg("failed to start grpc server")
			return
		}

		log.Info().Str("GRPCServerAddress", addr).
			Msg("started grpc server")

		if err = a.grpcServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic().Err(err).Msg("HTTP Server stopped")
		}
	}()
}

func (a *Application) startConsumer() <-chan models.Result {
	return a.consumer.Start()
}

func (a *Application) startProducer(msgChannel <-chan models.Result) {
	a.producer.Start(msgChannel)
}
