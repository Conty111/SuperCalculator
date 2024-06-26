package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/enums"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	App            *App
	DB             *DatabaseConfig
	HTTPConfig     *HTTPConfig
	BrokerCfg      *BrokerConfig
	JSONConfigPath string
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
	LoggerCfg          gin.LoggerConfig
	TimeoutResponse    time.Duration
	TimeToRetry        time.Duration
	Agents             []models.AgentConfig
	ApiToUse           enums.ApiType
	AuthPublicKeyPath  string
	AuthPrivateKeyPath string
	TokenTTL           time.Duration
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
	setJSONconfig(cfg)
	config = cfg

	return cfg
}

func setJSONconfig(cfg *Configuration) {
	file, err := os.Open(cfg.JSONConfigPath)
	if err != nil {
		d, _ := os.Getwd()
		log.Panic().Err(err).Str("current_dir", d).Msg("can't open json system_config")
	}
	defer file.Close()

	// Decode JSON from file
	decoder := json.NewDecoder(file)

	var jsonData models.JSONData
	if err := decoder.Decode(&jsonData); err != nil {
		log.Panic().Err(err).Msg("can't read json system_config")
	}

	brokers := make([]string, len(jsonData.Brokers))
	for i, broker := range jsonData.Brokers {
		brokers[i] = fmt.Sprintf("%s:%d", broker.Address, broker.Port)
	}

	cfg.BrokerCfg.Brokers = brokers
	cfg.App.Agents = jsonData.Agents
}

func getFromEnv() *Configuration {
	var cfg = &Configuration{}

	cfg.JSONConfigPath = envy.Get("JSON_CONFIG_PATH", "system_config.json")
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
	dbCfg.Path = envy.Get("PATH_TO_DB_FILE", "back-end/db/local.db")
	dbCfg.DBtype = envy.Get("DB_TYPE", "sqlite")
	port, err := strconv.Atoi(envy.Get("DB_PORT", "5432"))
	if err != nil {
		log.Panic().Err(err).Msg("cannot convert DB_PORT")
	}
	dbCfg.Port = port

	dbCfg.DSN = getDbDSN(dbCfg)

	return dbCfg
}

func getAppConf() *App {
	var cfg = App{}

	cfg.LoggerCfg = gin.LoggerConfig{}
	tResp, err := strconv.Atoi(envy.Get("TIMEOUT_RESPONSE", "5"))
	if err != nil {
		log.Error().Err(err).Msg("failed to get TIMEOUT_RESPONSE")
		tResp = 5
	}
	tRetry, err := strconv.Atoi(envy.Get("RETRY_INTERVAL", "5"))
	if err != nil {
		log.Error().Err(err).Msg("failed to get RETRY_INTERVAL")
		tRetry = 5
	}
	var apiType enums.ApiType
	apiType = enums.ApiType(envy.Get("AGENT_API_TYPE", "rest"))
	if apiType != enums.RestApi && apiType != enums.GrpcApi {
		log.Error().Err(errors.New("invalid api type")).Msg("should be rest or grpc")
		apiType = enums.RestApi
	}
	cfg.ApiToUse = apiType
	cfg.TimeToRetry = time.Duration(tRetry) * time.Second
	cfg.TimeoutResponse = time.Duration(tResp) * time.Second
	cfg.AuthPublicKeyPath = envy.Get("AUTH_PUBLIC_KEY_PATH", "")
	cfg.AuthPrivateKeyPath = envy.Get("AUTH_PRIVATE_KEY_PATH", "")
	ttl, err := strconv.Atoi(envy.Get("TOKEN_TTL", "700900"))
	if err != nil {
		log.Info().Msg("failed to parse token TTL from sys.env, setting default value")
		ttl = 700
	}
	cfg.TokenTTL = time.Second * time.Duration(ttl)

	return &cfg
}

func getBrokerConf() *BrokerConfig {
	var cfg = BrokerConfig{}

	cfg.ProduceTopic = envy.Get("TASKS_TOPIC", "tasks")
	cfg.ConsumeTopic = envy.Get("RES_TOPIC", "results")
	cfg.ConsumerGroup = envy.Get("CONSUMER_GROUP", "orkestrator_group")
	interval, err := strconv.Atoi(envy.Get("COMMIT_INTERVAL", "1"))
	if err != nil {
		log.Panic().Err(err).Msg("cannot convert COMMIT_INTERVAL")
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

	return &cfg
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
