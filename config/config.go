package config

import (
"fmt"
"os"

"github.com/youngprinnce/product-microservice/internal/logger"
"gopkg.in/yaml.v2"
)

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Env     string `yaml:"env"`
}

type Database struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DbName   string `yaml:"db_name"`
}

type Server struct {
	Listen string `yaml:"listen"`
	Port   string `yaml:"port"`
}

type Config struct {
	App      App      `yaml:"app"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

var conf Config

// Load loads configuration from environment or default file
func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "etc/config.yaml"
	}
	
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &conf, nil
}

// LoadConfig loads configuration from specified path (backwards compatibility)
func LoadConfig(path string) *Config {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("yamlFile.Get err   #%v ", err))
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unmarshal: %v", err))
	}
	return &conf
}

// GetConfig returns the loaded configuration
func GetConfig() *Config {
	return &conf
}
