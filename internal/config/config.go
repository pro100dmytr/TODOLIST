package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  struct {
		Address     string        `yaml:"address" env-default:"localhost:1010"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
		User        string        `yaml:"user" env-required:"true"`
		Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
	} `yaml:"http_server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
