package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/build"
)

// InitializeBuildInfo creates new build.Info
func InitializeBuildInfo() *build.Info {
	return build.NewInfo()
}
