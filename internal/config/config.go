package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Config represents the configuration structure
type Config struct {
	Port int `yaml:"port"`
	DBUser       string `yaml:"db_user"`
	DBPassword   string `yaml:"db_password"`
	DBHost       string `yaml:"db_host"`
	DBPort       int    `yaml:"db_port"`
	DBName       string `yaml:"db_name"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// LoadConfig loads the configuration from the YAML file
func LoadConfig() (Config, error) {
	config := Config{}

	file, err := os.Open("config.yaml")
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

