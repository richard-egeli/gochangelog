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

var validConfigFilenames = []string{
	"go.changelog.config.yaml",
	"gochangelog.config.yaml",
	"changelog.config.yaml",
	"go.changelog.yaml",
	"changelog.yaml",
}

func Read() (*YAML, error) {
	config := YAML{}
	for index, filename := range validConfigFilenames {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			if index < len(validConfigFilenames)-1 {
				continue
			}
			return nil, err
		}

		err = yaml.Unmarshal(bytes, &config)
		if err != nil {
			if index < len(validConfigFilenames)-1 {
				continue
			}

			return nil, err
		}

		break
	}

	return &config, nil
}
