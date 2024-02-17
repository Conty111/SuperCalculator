package kafka_broker

import (
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type AppProducer struct {
	Producer  sarama.SyncProducer
	Topic     string
	Partition int32
	InChan    chan models.Task
	Done      chan interface{}
}

func NewAppProducer(
	producer sarama.SyncProducer,
	topic string) *AppProducer {

	return &AppProducer{
		InChan:   make(chan models.Task),
		Producer: producer,
		Done:     make(chan interface{}),
		Topic:    topic,
	}
}

func (ap *AppProducer) Start() {
	go func() {
		defer func(Producer sarama.SyncProducer) {
			err := Producer.Close()
			if err != nil {
				log.Error().Err(err).Msg("Error while closing producer")
			}
		}(ap.Producer)
		log.Info().Msg("producer are ready to send messages")
		for {
			select {
			case <-ap.Done:
				return
			case msg := <-ap.InChan:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Error().Err(err).Msg("Error while trying to marshal result")
					continue
				}
				prodMsg := &sarama.ProducerMessage{
					Topic: ap.Topic,
					Value: sarama.ByteEncoder(data),
				}
				log.Info().Str("task", string(data)).Msg("sending task")
				p, _, err := ap.Producer.SendMessage(prodMsg)
				if err != nil {
					log.Error().Int32("Partition", p).Str("message", string(data)).Err(err).Msg("Error while sending message")
					continue
				}
			}
		}
	}()
}

func (ap *AppProducer) Stop() {
	ap.Done <- "stop"
}
