package config

import (
	"context"
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/system_config"
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
	BrokerCfg      BrokerConfig
	JSONConfigPath string
}

type App struct {
	Name      string
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

	return cfg
}

func setJSONconfig(cfg *Configuration, num int) {
	file, err := os.Open(cfg.JSONConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("can't open json system_config")
	}
	defer file.Close()

	// Decode JSON from file
	decoder := json.NewDecoder(file)

	var jsonData system_config.JSONData
	if err := decoder.Decode(&jsonData); err != nil {
		log.Fatal().Err(err).Msg("can't read json system_config")
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
	brokers := make([]string, len(jsonData.Brokers))
	for i, broker := range jsonData.Brokers {
		brokers[i] = broker.Address
	}
	cfg.BrokerCfg.Brokers = brokers

}

func getConsumerConf() BrokerConfig {
	var cfg = BrokerConfig{}
	cfg.ConsumeTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ProduceTopic = envy.Get("RES_TOPIC", "results")
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
	//cfg.Port = envy.Get("HTTP_SERVER_PORT", "8000")

	return cfg
}
