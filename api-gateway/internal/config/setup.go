// Package config is a nice package
package config

import (
	"fmt"

	"github.com/spf13/viper"
)
type Config struct {
	Clients    Clients    `mapstructure:"clients"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
}

type HTTPServer struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type Clients struct {
	Orders Orders `mapstructure:"orders"`
}

type Orders struct {
	GRPCPort int    `mapstructure:"grpc_port"`
	GRPCHost string `mapstructure:"grpc_host"`
}


func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &cfg, nil
}

func MustLoad(path string) (cfg *Config) {
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}
	return
}
