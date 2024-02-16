package kafka_broker

import (
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type AppProducer struct {
	Producer  sarama.AsyncProducer
	Monitor   *services.Monitor
	Topic     string
	Partition int32
	Done      chan interface{}
}

func NewAppProducer(
	producer sarama.AsyncProducer,
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
				log.Info().Msg("sending result")
				ap.Producer.Input() <- &prodMsg
			}
		}
	}()
}

func (ap *AppProducer) Stop() {
	ap.Done <- "stop"
}
