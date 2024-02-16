package kafka_broker

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"time"
)

type AppConsumer struct {
	Service  *services.ExpressionService
	Consumer sarama.PartitionConsumer
	Monitor  *services.Monitor
	Done     chan interface{}
}

func NewAppConsumer(svc *services.ExpressionService,
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
				res, err := ac.Proccess(message)
				if err != nil {
					log.Debug().Err(err).Msg("invalid message")
				} else {
					out <- *res
				}
				log.Print(ac.Consumer.HighWaterMarkOffset())
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
	res, err := ac.Service.Proccess(msg)
	log.Info().Str("time of calculation", time.Since(t1).String()).Msg("calculated")
	// TODO handle errors here
	return res, err
}
