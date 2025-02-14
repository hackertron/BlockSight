package config

import (
	"time"

	"github.com/spf13/viper"
)

type KafkaTopics struct {
	Blocks string
}

type Config struct {
	Alchemy struct {
		APIKey     string        `mapstructure:"api_key"`
		Network    string        `mapstructure:"network"`
		RetryCount int           `mapstructure:"retry_count"`
		RetryDelay time.Duration `mapstructure:"retry_delay"`
		RateLimit  int           `mapstructure:"rate_limit"`
	}
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topics  KafkaTopics
	}
	Metrics struct {
		Enabled bool `mapstructure:"enabled"`
		Port    int  `mapstructure:"port"`
	}
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
