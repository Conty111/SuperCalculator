package kafka_broker

import (
	"encoding/json"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/models"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"strconv"
)

type AppConsumer struct {
	Service  *services.ExpressionService
	Consumer sarama.Consumer
	Config   *config.ConsumerConfig
	Monitor  *services.Monitor
	Done     chan interface{}
}

func NewAppConsumer(svc *services.ExpressionService,
	consumer sarama.Consumer,
	cfg *config.ConsumerConfig,
	mon *services.Monitor) *AppConsumer {
	return &AppConsumer{
		Service:  svc,
		Consumer: consumer,
		Config:   cfg,
		Monitor:  mon,
		Done:     make(chan interface{}),
	}
}

func (ac *AppConsumer) Start() <-chan models.Result {
	out := make(chan models.Result)
	go func(out chan<- models.Result) {
		consumerGroup, err := ac.Consumer.ConsumePartition(
			ac.Config.Topic,
			ac.Config.Partition,
			sarama.OffsetNewest)
		if err != nil {
			log.Panic().Err(err).Msg("receiver stopped")
		}
		log.Info().Msg("started receiver")
		defer func() {
			if err = consumerGroup.Close(); err != nil {
				// Обработка ошибки при закрытии
				log.Panic().Err(err).Msg("receiver stopped")
			}
		}()
		select {
		case <-ac.Done:
			return
		default:
			for message := range consumerGroup.Messages() {
				// Обработка полученного сообщения
				res, err := ac.Proccess(message)
				if err != nil {
					log.Fatal().Err(err)
				}
				out <- *res
			}
		}
	}(out)
	return out
}

func (ac *AppConsumer) Stop() {
	log.Info().Msg("receiver graceful stopped")
	ac.Done <- "stop"
}

func (ac *AppConsumer) Proccess(msg *sarama.ConsumerMessage) (*models.Result, error) {
	log.Info().
		Time("start_time", msg.Timestamp).
		Msg(fmt.Sprintf("Starting processing a message: %s", string(msg.Value)))

	var t models.Task
	if err := json.Unmarshal(msg.Value, &t); err != nil {
		return nil, err
	}
	key, err := strconv.Atoi(string(msg.Key))
	if err != nil {
		return nil, err
	}
	t.ID = uint(key)
	expression, err := ac.Service.ValidateExpression(t.Expression)
	if err != nil {
		return nil, err
	}
	resNum, err := ac.Service.Calculate(expression)
	if err != nil {
		return nil, err
	}
	return &models.Result{
		Task:  t,
		Value: resNum,
	}, nil
}
