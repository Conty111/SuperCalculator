package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	App            App
	HTTPConfig     HTTPConfig
	GRPCConfig     GRPCConfig
	BrokerCfg      BrokerConfig
	JSONConfigPath string
}

type App struct {
	Name      string
	LoggerCfg gin.LoggerConfig
}

type GRPCConfig struct {
	Host string
	Port string
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

func GetConfig(ctx context.Context) *Configuration {
	if config != nil {
		return config
	}

	cfg := getFromEnv()
	setJSONconfig(cfg, ctx.Value("index").(int))

	config = cfg

	return cfg
}

func getFromEnv() *Configuration {
	var cfg = &Configuration{}

	cfg.JSONConfigPath = envy.Get("JSON_CONFIG_PATH", "system_config.json")
	cfg.App = getAppConf()
	cfg.HTTPConfig = getWebConf()
	cfg.BrokerCfg = getConsumerConf()
	cfg.GRPCConfig = getGrpcConf()

	return cfg
}

func getGrpcConf() GRPCConfig {
	var cfg GRPCConfig
	cfg.Host = envy.Get("GRPC_HOST", "localhost")
	cfg.Port = envy.Get("GRPC_PORT", "5000")
	return cfg
}

func setJSONconfig(cfg *Configuration, num int) {
	file, err := os.Open(cfg.JSONConfigPath)
	if err != nil {

		log.Panic().Err(err).Msg("can't open json system_config")
	}
	defer file.Close()

	// Decode JSON from file
	decoder := json.NewDecoder(file)

	var jsonData models.JSONData
	if err := decoder.Decode(&jsonData); err != nil {
		log.Panic().Err(err).Msg("can't read json system_config")
	}
	if num > len(jsonData.Agents) {
		log.Fatal().Err(err).Msg("invalid argument or json system_config")
	}

	agentCfg := jsonData.Agents[num]

	cfg.App.Name = agentCfg.Name
	cfg.HTTPConfig.Port = strconv.Itoa(agentCfg.HttpPort)

	cfg.BrokerCfg.Partition = agentCfg.BrokerPartition
	cfg.BrokerCfg.ConsumerGroup = agentCfg.ConsumerGroup
	cfg.BrokerCfg.CommitInterval = agentCfg.BrokerCommitInterval
	log.Print(jsonData)
	brokers := make([]string, len(jsonData.Brokers))
	for i, broker := range jsonData.Brokers {
		brokers[i] = fmt.Sprintf("%s:%d", broker.Address, broker.Port)
	}
	cfg.BrokerCfg.Brokers = brokers

}

func getConsumerConf() BrokerConfig {
	var cfg = BrokerConfig{}
	cfg.ConsumeTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ProduceTopic = envy.Get("RES_TOPIC", "results")
	interval, err := strconv.Atoi(envy.Get("COMMIT_INTERVAL", "1"))
	if err != nil {
		log.Panic().Err(err).Msg("error converting COMMIT_INTERVAL")
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

	return cfg
}
