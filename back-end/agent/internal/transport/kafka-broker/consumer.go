package kafka_broker

import (
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/agent_errors"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"time"
)

type AppConsumer struct {
	Service  *services.CalculatorService
	Consumer sarama.PartitionConsumer
	Monitor  *services.Monitor
	Done     chan interface{}
}

func NewAppConsumer(svc *services.CalculatorService,
	consumer sarama.PartitionConsumer,
	mon *services.Monitor) *AppConsumer {
	return &AppConsumer{
		Service:  svc,
		Consumer: consumer,
		Monitor:  mon,
		Done:     make(chan interface{}, 5),
	}
}

func (ac *AppConsumer) Start() <-chan models.Result {
	out := make(chan models.Result)
	go func() {
		defer func(Consumer sarama.PartitionConsumer) {
			err := Consumer.Close()
			if err != nil {
				log.Error().Err(err).Msg("error while closing consumer")
			}
		}(ac.Consumer)
		for {
			select {
			case <-ac.Done:
				return
			case message := <-ac.Consumer.Messages():
				// Обработка полученного сообщения
				ac.Monitor.AddWork()
				res, err := ac.Proccess(message)
				if err != nil {
					log.Debug().Err(err).Msg("invalid message")
				} else {
					out <- *res
				}
			}
		}
	}()
	return out
}

func (ac *AppConsumer) Stop() {
	ac.Done <- "stop"
}

func (ac *AppConsumer) Proccess(msg *sarama.ConsumerMessage) (*models.Result, error) {
	t1 := time.Now()
	log.Info().
		Time("start_time", msg.Timestamp).
		Str("message", string(msg.Value)).
		Msg("started processing a message")
	res := ac.Service.Execute(ac.parseMessageToTask(msg))
	if res == nil {
		return nil, agent_errors.ErrInvalidMessage
	}
	log.Info().Str("time of calculation", time.Since(t1).String()).Msg("calculated")
	return res, nil
}

func (ac *AppConsumer) parseMessageToTask(msg *sarama.ConsumerMessage) *models.Task {
	var t models.Task

	if err := json.Unmarshal(msg.Value, &t); err != nil {
		log.Error().Msg("Error while parsing json")
		return nil
	}
	log.Debug().Any("task", t).Str("key", string(msg.Key)).Msg("parsed to json")
	return &t
}
