package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app/dependencies"
	kafka_broker "github.com/Conty111/SuperCalculator/back-end/agent/internal/transport/kafka-broker"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"time"
)

func InitializeConsumer(container *dependencies.Container) *kafka_broker.AppConsumer {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = time.Duration(container.Config.ConsumerCfg.CommitInterval) * time.Second

	consumer, err := sarama.NewConsumer(container.Config.ConsumerCfg.Brokers, cfg)
	if err != nil {
		log.Error().Msg("Error creating Kafka receiver")
		log.Fatal().Err(err)
	}
	return kafka_broker.NewAppConsumer(container.ExpressionSvc, consumer, &container.Config.ConsumerCfg, container.Monitor)
}
