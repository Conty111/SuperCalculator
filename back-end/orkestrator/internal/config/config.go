package config

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type Configuration struct {
	App        *App
	DB         *DatabaseConfig
	HTTPConfig *HTTPConfig
	BrokerCfg  *BrokerConfig
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

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	SSLMode  string
	DSN      string
}

type App struct {
	LoggerCfg gin.LoggerConfig
}

type HTTPConfig struct {
	Host string
	Port string
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
	cfg.DB = getDBConfig()
	cfg.HTTPConfig = getWebConf()

	return cfg
}

func getDBConfig() *DatabaseConfig {
	dbCfg := &DatabaseConfig{}

	dbCfg.Host = envy.Get("DB_HOST", "localhost")
	dbCfg.User = envy.Get("DB_USER", "user")
	dbCfg.Password = envy.Get("DB_PASSWORD", "sqlite")
	dbCfg.DBName = envy.Get("DB_NAME", "test")
	dbCfg.SSLMode = envy.Get("DB_SSLMODE", "disable")
	dbCfg.Port = getDbPort()
	dbCfg.DSN = getDbDSN(dbCfg)

	return dbCfg
}

func getAppConf() *App {
	var cfg = App{}

	cfg.LoggerCfg = gin.LoggerConfig{}

	return &cfg
}

func getConsumerConf() BrokerConfig {
	var cfg = BrokerConfig{}
	cfg.Brokers = strings.Split(envy.Get("BROKERS", "kafka-broker-broker:9092"), ";")
	cfg.ProduceTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ConsumeTopic = envy.Get("RESULTS_TOPIC", "results")
	cfg.ConsumerGroup = envy.Get("CONSUMER_GROUP", "orkestrator_group")
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

func getWebConf() *HTTPConfig {
	var cfg = HTTPConfig{}

	cfg.Host = envy.Get("HTTP_SERVER_HOST", "0.0.0.0")
	cfg.Port = envy.Get("HTTP_SERVER_PORT", "8000")

	return &cfg
}

func getDbPort() int {
	port, err := strconv.Atoi(envy.Get("DB_PORT", "5432"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get DB_PORT")
	}
	return port
}

func getDbDSN(dbConfig *DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode,
	)
}
