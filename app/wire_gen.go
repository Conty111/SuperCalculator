// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package app

import (
	"github.com/ildarusmanov/gobase/app/dependencies"
	"github.com/ildarusmanov/gobase/app/initializers"
)

// Injectors from wire.go:

func BuildApplication() (*Application, error) {
	info := initializers.InitializeBuildInfo()
	container := &dependencies.Container{
		BuildInfo: info,
	}
	engine := initializers.InitializeRouter(container)
	httpServerConfig := initializers.InitializeHTTPServerConfig(engine)
	server, err := initializers.InitializeHTTPServer(httpServerConfig)
	if err != nil {
		return nil, err
	}
	application := &Application{
		httpServer: server,
		Container:  container,
	}
	return application, nil
}
