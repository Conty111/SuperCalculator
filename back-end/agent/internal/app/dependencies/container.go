package dependencies

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/build"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
)

// Container is a DI container for application
type Container struct {
	BuildInfo  *build.Info
	Config     *config.Configuration
	Monitor    *services.Monitor
	Calculator *services.CalculatorService
}
