package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		MetricsPort int    `yaml:"metricsPort"`
	} `yaml:"server"`
	Tzkt struct {
		Url string `yaml:"url"`
	} `yaml:"tzkt"`
	Db struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`
}

// var config Config

func LoadConfig(configPath string) (Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.Server.Port == 0 {
		return Config{}, errors.New("server port is required")
	}

	if cfg.Server.MetricsPort == 0 {
		return Config{}, errors.New("server metrics port is required")
	}

	if cfg.Tzkt.Url == "" {
		return Config{}, errors.New("tzkt url is required")
	}

	if cfg.Db.Host == "" {
		return Config{}, errors.New("db host is required")
	}

	if cfg.Db.Port == 0 {
		return Config{}, errors.New("db port is required")
	}

	if cfg.Db.User == "" {
		return Config{}, errors.New("db user is required")
	}

	if cfg.Db.Password == "" {
		return Config{}, errors.New("db password is required")
	}

	if cfg.Db.Database == "" {
		return Config{}, errors.New("db database is required")
	}

	// config = cfg
	return cfg, nil
}
