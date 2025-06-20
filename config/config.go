package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the structure of your YAML file
type Config struct {
	Vus string `yaml:"vus"`
	Duration string `yaml:"duration"`
}

func ReadConfig() (*Config, error) {
	// Open the YAML file
	file, err := os.Open("config/default.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the YAML file
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	log.Println("Config loaded successfully:", config)
	return &config, nil
}