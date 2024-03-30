package dependencies

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/build"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/kafka-broker"
	"gorm.io/gorm"
)

// Container is a DI container for application
type Container struct {
	BuildInfo    *build.Info
	Database     *gorm.DB
	Config       *config.Configuration
	Consumer     *kafka_broker.AppConsumer
	Producer     *kafka_broker.AppProducer
	TaskManager  interfaces.TaskManager
	AgentManager interfaces.AgentManager
}
