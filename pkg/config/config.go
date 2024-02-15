package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YAML struct {
	Output  string   `yaml:"output"`
	Version string   `yaml:"apiVersion"`
	Filter  []string `yaml:"filter"`
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

	fmt.Printf("Dump: \n%v\n\n", config)
	return &config, nil
}
