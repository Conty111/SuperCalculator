package config

import (
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
	LoggerCfg gin.LoggerConfig
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

	globalEnv := envy.Get("GLOBAL_ENV", "../../.env")
	err := envy.Load(globalEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load global env")
	}
	cfg.App = getAppConf()
	cfg.HTTPConfig = getWebConf()
	cfg.BrokerCfg = getConsumerConf()

	return cfg
}

func getConsumerConf() BrokerConfig {
	var cfg = BrokerConfig{}
	cfg.Brokers = strings.Split(envy.Get("BROKERS", "kafka-broker-broker:9092"), ";")
	cfg.ConsumeTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ProduceTopic = envy.Get("RES_TOPIC", "results")
	cfg.ConsumerGroup = envy.Get("CONSUMER_AGENT_GROUP", "agent_group")
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
	cfg.SaramaCfg.Consumer.Offsets.AutoCommit.Interval = 500 * time.Millisecond

	return cfg
}

func getAppConf() App {
	var cfg = App{}

	cfg.LoggerCfg = gin.LoggerConfig{}

	return cfg
}

func getWebConf() HTTPConfig {
	var cfg = HTTPConfig{}

	cfg.Host = envy.Get("HTTP_SERVER_HOST", "0.0.0.0")
	cfg.Port = envy.Get("HTTP_SERVER_PORT", "8000")

	return cfg
}
