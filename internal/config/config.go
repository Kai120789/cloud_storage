package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

var AppConfig *Config

type Config struct {
	ServerAddress string
	DBDSN         string
	LogLevel      string
	SecretKey     string
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "URL and port to run server")
	flag.StringVar(&cfg.DBDSN, "d", "", "DBDSN for database")

	cfg.SecretKey = getEnvStringOrDefault("SECRET_KEY", "default")
	cfg.ServerAddress = getEnvStringOrDefault("SERVER_ADDRESS", "localhost:8080")
	cfg.DBDSN = getEnvStringOrDefault("DBDSN", "")

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	flag.Parse()

	AppConfig = cfg

	return cfg, nil
}

func getEnvStringOrDefault(name, defaultValue string) string {
	if envString := os.Getenv(name); envString != "" {
		return envString
	}

	return defaultValue
}

func getEnvIntOrDefault(name string, defaultValue int) (*int, error) {
	if envInt := os.Getenv(name); envInt != "" {
		intEnvInt, err := strconv.Atoi(envInt)
		if err != nil {
			return nil, err
		}
		return &intEnvInt, nil
	}

	return &defaultValue, nil
}
