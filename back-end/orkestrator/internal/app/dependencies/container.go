package dependencies

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/build"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/config"
)

// Container is a DI container for application
type Container struct {
	BuildInfo *build.Info
	Config    *config.Configuration
}
