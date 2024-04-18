package kafka_broker

import (
	"context"
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"sync"
)

type AppConsumer struct {
	Service  interfaces.TaskManager
	Consumer sarama.Consumer
	Topic    string
	Done     chan interface{}
	lock     sync.Mutex
}

func NewAppConsumer(svc interfaces.TaskManager,
	consumer sarama.Consumer, topic string) *AppConsumer {
	return &AppConsumer{
		Topic:    topic,
		Service:  svc,
		Consumer: consumer,
		Done:     make(chan interface{}, 5),
		lock:     sync.Mutex{},
	}
}

func (ac *AppConsumer) Start() {
	go func() {
		partList, err := ac.Consumer.Partitions(ac.Topic)
		if err != nil {
			log.Panic().
				Err(err).
				Msg("error to get partitions from topic")
		}
		wg := sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.TODO())
		for _, part := range partList {
			pc, err := ac.Consumer.ConsumePartition(ac.Topic, part, sarama.OffsetNewest)
			if err != nil {
				log.Panic().
					Err(err).
					Int32("Partition", part).
					Msg("error creating consumer for partition")
			}
			wg.Add(1)
			go func(pc sarama.PartitionConsumer, ctx context.Context) {
				defer wg.Done()
				for {
					select {
					case message := <-pc.Messages():
						log.Info().
							Int32("Partition", message.Partition).
							Str("message", string(message.Value)).
							Msg("got message")
						res, err := parseMessage(message)
						if err == nil && res != nil {
							err = ac.Service.SaveResult(res)
							if err != nil {
								log.Error().Err(err).Msg("failed to save result")
							}
						}
					case <-ctx.Done():
						return
					}
				}
			}(pc, ctx)
		}
		log.Info().Msg("consumer are ready to receive messages")
		<-ac.Done
		cancel()
		wg.Wait()
	}()
}

func (ac *AppConsumer) Stop() {
	ac.Done <- "stop"
}

func parseMessage(msg *sarama.ConsumerMessage) (*models.Result, error) {
	var res *models.Result
	err := json.Unmarshal(msg.Value, &res)
	if err != nil {
		log.Error().Err(err).Str("message", string(msg.Value)).Msg("Failed to parse message to json")
		return nil, err
	}
	return res, nil
}
