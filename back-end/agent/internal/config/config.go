package config

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/enums"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type Configuration struct {
	App        App
	HTTPConfig HTTPConfig
	BrokerCfg  BrokerConfig
}

type App struct {
	TimeFormat    enums.TimeFormat
	LoggerCfg     gin.LoggerConfig
	TasksBuffSize uint
}

type HTTPConfig struct {
	Host string
	Port string
}

type BrokerConfig struct {
	SaramaCfg      *sarama.Config
	Brokers        []string
	ConsumerGroup  string
	CommitInterval uint
	ConsumeTopic   string
	ProduceTopic   string
	Partition      int32
}

var config *Configuration

func GetConfig() *Configuration {
	if config != nil {
		return config
	}

	cfg := getFromEnv()
	config = cfg

	return cfg
}

func getFromEnv() *Configuration {
	var cfg = &Configuration{}

	cfg.App = getAppConf()
	cfg.HTTPConfig = getWebConf()
	cfg.BrokerCfg = getConsumerConf()

	return cfg
}

func getConsumerConf() BrokerConfig {
	var cfg = BrokerConfig{}
	cfg.Brokers = strings.Split(envy.Get("BROKERS", "kafka-broker-broker:9092"), ";")
	cfg.ConsumeTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ProduceTopic = envy.Get("RESULTS_TOPIC", "results")
	cfg.ConsumerGroup = envy.Get("CONSUMER_GROUP", "agent_group")
	interval, err := strconv.Atoi(envy.Get("COMMIT_INTERVAL", "1"))
	if err != nil {
		log.Fatal().Err(err)
	}
	cfg.CommitInterval = uint(interval)
	cfg.SaramaCfg = sarama.NewConfig()
	cfg.SaramaCfg.Consumer.Return.Errors = true
	cfg.SaramaCfg.Producer.Return.Successes = true
	cfg.SaramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.SaramaCfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.SaramaCfg.Consumer.Offsets.AutoCommit.Interval = 100 * time.Millisecond

	return cfg
}

func getAppConf() App {
	var cfg = App{}

	format := envy.Get("TIME_FORMAT", "RFC3339")
	switch strings.ToUpper(format) {
	case "RFC3339":
		cfg.TimeFormat = enums.RFC3339
	case "RFC3339Nano":
		cfg.TimeFormat = enums.RFC3339Nano
	}
	cfg.LoggerCfg = gin.LoggerConfig{}
	size, err := strconv.Atoi(envy.Get("EXPRESSIONS_BUFF_SIZE", "10"))
	if err != nil {
		log.Fatal().Err(err).Msg("error while getting buff size")
	}
	cfg.TasksBuffSize = uint(size)

	return cfg
}

func getWebConf() HTTPConfig {
	var cfg = HTTPConfig{}

	cfg.Host = envy.Get("HTTP_SERVER_HOST", "0.0.0.0")
	cfg.Port = envy.Get("HTTP_SERVER_PORT", "8000")

	return cfg
}
