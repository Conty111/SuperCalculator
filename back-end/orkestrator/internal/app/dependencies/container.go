package dependencies

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/build"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"gorm.io/gorm"
)

// Container is a DI container for application
type Container struct {
	BuildInfo *build.Info
	Database  *gorm.DB
	Config    *config.Configuration
}
