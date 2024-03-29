package config

import (
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *domain.Config) error {
	if err := envconfig.Process("db", &cfg.Db); err != nil {
		return err
	}
	return nil
}

func unmarshal(cfg *domain.Config) error {
	if err := viper.UnmarshalKey("expression", &cfg.Expression); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("app_http", &cfg.AppHttpServer); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("ampq_queue_struct", &cfg.AMPQStructCfg); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("ampq_queue_server", &cfg.AMPQServerConfig); err != nil {
		return err
	}
	return nil
}

func Init(configDir string) (*domain.Config, error) {
	viper.SetConfigName("config")
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var cfg domain.Config
	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	logrus.Infof("Config: %+v", cfg)
	return &cfg, nil
}
