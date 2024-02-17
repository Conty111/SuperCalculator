package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/kafka-broker"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

func InitializeProducer(container *dependencies.Container) *kafka_broker.AppProducer {
	producer, err := sarama.NewAsyncProducer(container.Config.BrokerCfg.Brokers, container.Config.BrokerCfg.SaramaCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize producer")
	}
	return kafka_broker.NewAppProducer(
		producer,
		container.Config.BrokerCfg.ProduceTopic,
	)
}
