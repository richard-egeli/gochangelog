package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YAML struct {
	RepoURL  string   `yaml:"repoURL"`
	Provider string   `yaml:"provider"`
	Output   string   `yaml:"output"`
	Version  string   `yaml:"apiVersion"`
	Filter   []string `yaml:"filter"`
}

func Read() (*YAML, error) {
	config := YAML{}
	bytes, err := os.ReadFile("gochangelog.config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
