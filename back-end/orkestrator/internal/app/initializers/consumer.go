package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/dependencies"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/kafka-broker"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

func InitializeConsumer(container *dependencies.Container) *kafka_broker.AppConsumer {
	con, err := sarama.NewConsumer(
		container.Config.BrokerCfg.Brokers,
		container.Config.BrokerCfg.SaramaCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating Kafka consumer group")
	}
	log.Info().Msg("initialized consumer")
	return kafka_broker.NewAppConsumer(container.Service, con, container.Config.BrokerCfg.ConsumeTopic)
}
