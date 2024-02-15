package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version string `yaml:"apiVersion"`
	Filter  []string
}

func Read() (*Config, error) {
	config := Config{}
	bytes, err := os.ReadFile("gochangelog.config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Dump: \n%v\n\n", config)
	return &config, nil
}
