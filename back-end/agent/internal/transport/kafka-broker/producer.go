package kafka_broker

import (
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type AppProducer struct {
	Producer  sarama.SyncProducer
	Monitor   *services.Monitor
	Topic     string
	Partition int32
	Done      chan interface{}
}

func NewAppProducer(
	producer sarama.SyncProducer,
	mon *services.Monitor,
	topic string,
	partition int32) *AppProducer {

	return &AppProducer{
		Producer:  producer,
		Monitor:   mon,
		Done:      make(chan interface{}, 5),
		Topic:     topic,
		Partition: partition,
	}
}

func (ap *AppProducer) Start(messages <-chan models.Result) {
	go func() {
		var prodMsg sarama.ProducerMessage
		prodMsg.Topic = ap.Topic
		prodMsg.Partition = ap.Partition
		for {
			select {
			case <-ap.Done:
				return
			case msg := <-messages:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Error().Err(err).Msg("Error while trying to marshal result")
				}

				prodMsg.Value = sarama.ByteEncoder(data)

				log.Info().Str("task", string(data)).Msg("sending result")

				p, _, err := ap.Producer.SendMessage(&prodMsg)
				if err != nil {
					log.Error().Int32("Partition", p).Str("message", string(data)).Err(err).Msg("Error while sending message")
				}

				err = ap.Monitor.CompleteTask(msg.ID)
				if err != nil {
					log.Error().Err(err).Msg("error completing task in monitor")
				}
			}
		}
	}()
}

func (ap *AppProducer) Stop() {
	ap.Done <- "stop"
}
