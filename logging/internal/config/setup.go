// Package config is a nice package
package config

import "github.com/spf13/viper"

type Config struct {
	Kafka `mapstructure:"kafka"`
	Logging `mapstructure:"logging"`
}

type Kafka struct {
	BrokerAdresses []string `mapstructure:"broker_adresses"`
	ConsumerGroup string `mapstructure:"consumer_group"`
	Topic string `mapstructure:"topic"`
	NumConsumers int `mapstructure:"num_consumers"`
}

type Logging struct {
	LogBatchSize int `mapstructure:"log_batch_size"`
}

func MustLoad(pathToFile string) *Config {
	cfg, err := Load(pathToFile)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load(pathToFile string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(pathToFile)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
