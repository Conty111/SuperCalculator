package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/dependencies"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/kafka-broker"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

func InitializeConsumer(container *dependencies.Container) *kafka_broker.AppConsumer {
	consumer, err := sarama.NewConsumer(container.Config.BrokerCfg.Brokers, container.Config.BrokerCfg.SaramaCfg)
	log.Print(container.Config.BrokerCfg.Brokers)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating Kafka consumer")
	}
	con, err := consumer.ConsumePartition(
		container.Config.BrokerCfg.ConsumeTopic,
		container.Config.BrokerCfg.Partition,
		sarama.OffsetNewest)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating Kafka consumer")
	}
	log.Info().Str("Partition", string(container.Config.BrokerCfg.Partition)).Msg("started consumer")
	return kafka_broker.NewAppConsumer(container.Calculator, con, container.Monitor)
}
