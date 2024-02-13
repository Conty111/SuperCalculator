package config

import (
	"github.com/Conty111/SuperCalculator/back-end/agent/internal/enums"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type Configuration struct {
	App        App
	DB         DB
	HTTPConfig HTTPConfig
}

type DB struct {
	Port     uint
	Host     string
	Password string
	User     string
	DBName   string
}

type App struct {
	TimeFormat        enums.TimeFormat
	AuthPublicKeyPath string
	LoggerCfg         gin.LoggerConfig
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
	cfg.DB = getDatabaseConf()
	cfg.HTTPConfig = getWebConf()

	return cfg
}

func getDatabaseConf() DB {
	var cfg = DB{}

	cfg.Host = envy.Get("DB_HOST", "0.0.0.0")
	port, err := strconv.Atoi(envy.Get("DB_PORT", "27017"))
	if err != nil {
		log.Fatal().Err(err)
	}
	cfg.Port = uint(port)
	cfg.DBName = envy.Get("DB_NAME", "test")
	cfg.User = envy.Get("DB_USER", "mongo-repo")
	cfg.Password = envy.Get("DB_PASSWORD", "mongo-repo")

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
	cfg.AuthPublicKeyPath = envy.Get("AUTH_PUBLIC_KEY_PATH", "cert/ec-prime256v1-pub-key.pem")

	return cfg
}

func getWebConf() HTTPConfig {
	var cfg = HTTPConfig{}

	cfg.Host = envy.Get("HTTP_SERVER_HOST", "0.0.0.0")
	cfg.Port = envy.Get("HTTP_SERVER_PORT", "8000")

	return cfg
}
