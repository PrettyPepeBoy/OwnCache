package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env             string        `yaml:"env"`
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	HTTPServer      `yaml:"http_server"`
}

type HTTPServer struct {
	Port        string        `yaml:"port"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	Timeout     time.Duration `yaml:"timout"`
}

func MustSetupConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("file does not exist :%s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file :%s", configPath)
	}

	return &cfg
}
