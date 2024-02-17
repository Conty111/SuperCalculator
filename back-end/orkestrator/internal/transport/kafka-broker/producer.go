package kafka_broker

import (
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type AppProducer struct {
	Producer  sarama.AsyncProducer
	Topic     string
	Partition int32
	InChan    chan models.Task
	Done      chan interface{}
}

func NewAppProducer(
	producer sarama.AsyncProducer,
	topic string,
	partition int32) *AppProducer {

	return &AppProducer{
		InChan:    make(chan models.Task),
		Producer:  producer,
		Done:      make(chan interface{}),
		Topic:     topic,
		Partition: partition,
	}
}

func (ap *AppProducer) Start() {
	go func() {
		var prodMsg sarama.ProducerMessage
		prodMsg.Topic = ap.Topic
		prodMsg.Partition = ap.Partition
		log.Info().Msg("producer are ready to send tasks")
		for {
			select {
			case <-ap.Done:
				return
			case msg := <-ap.InChan:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Error().Err(err).Msg("Error while trying to marshal result")
				}
				prodMsg.Value = sarama.ByteEncoder(data)
				log.Info().Str("result", string(data)).Msg("sending result")
				ap.Producer.Input() <- &prodMsg
			}
		}
	}()
}

func (ap *AppProducer) Stop() {
	ap.Done <- "stop"
}
