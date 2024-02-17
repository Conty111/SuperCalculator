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
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	SSLMode  string
	DSN      string
	Path     string
	DBtype   string
}

type App struct {
	LoggerCfg  gin.LoggerConfig
	AgentCount uint
}

type HTTPConfig struct {
	Host           string
	Port           string
	AgentAddresses []string
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
	globalEnvPath := envy.Get("GLOBAL_ENV", "../../.env")
	err := envy.Load(globalEnvPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load global env")
	}
	cfg.App = getAppConf()
	cfg.DB = getDBConfig()
	cfg.HTTPConfig = getWebConf()
	cfg.BrokerCfg = getBrokerConf()

	return cfg
}

func getDBConfig() *DatabaseConfig {
	dbCfg := &DatabaseConfig{}

	dbCfg.Host = envy.Get("DB_HOST", "localhost")
	dbCfg.User = envy.Get("DB_USER", "user")
	dbCfg.Password = envy.Get("DB_PASSWORD", "sqlite")
	dbCfg.DBName = envy.Get("DB_NAME", "test")
	dbCfg.SSLMode = envy.Get("DB_SSLMODE", "disable")
	dbCfg.Path = envy.Get("PATH_TO_DB", "../db/test.db")
	dbCfg.DBtype = envy.Get("DB_TYPE", "sqlite")
	dbCfg.Port = getDbPort()
	dbCfg.DSN = getDbDSN(dbCfg)

	return dbCfg
}

func getAppConf() *App {
	var cfg = App{}

	cfg.LoggerCfg = gin.LoggerConfig{}
	agentCount, err := strconv.Atoi(envy.Get("AGENT_COUNT", "1"))
	if err != nil {
		log.Fatal().Err(err).Msg("error getting AGENT_COUNT")
	}
	cfg.AgentCount = uint(agentCount)

	return &cfg
}

func getBrokerConf() *BrokerConfig {
	var cfg = BrokerConfig{}
	cfg.Brokers = strings.Split(envy.Get("BROKERS", "kafka-broker-broker:9092"), ";")
	cfg.ProduceTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ConsumeTopic = envy.Get("RES_TOPIC", "results")
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

	return &cfg
}

func getWebConf() *HTTPConfig {
	var cfg = HTTPConfig{}

	cfg.Host = envy.Get("HTTP_SERVER_HOST", "0.0.0.0")
	cfg.Port = envy.Get("HTTP_SERVER_PORT", "8000")
	hosts := envy.Get("HTTP_AGENT_ADDRESSES", "localhost")
	cfg.AgentAddresses = strings.Split(hosts, ";")

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
	switch dbConfig.DBtype {
	case "sqlite":
		return dbConfig.Path
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			dbConfig.Host, dbConfig.User, dbConfig.Password,
			dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode,
		)
	default:
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			dbConfig.Host, dbConfig.User, dbConfig.Password,
			dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode,
		)
	}
}
