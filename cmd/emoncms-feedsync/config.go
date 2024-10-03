package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerInfo struct {
	Host   string `yaml:"host"`
	APIKey string `yaml:"apikey"`
}

type Config struct {
	FeedsFilter []string   `yaml:"feeds"`
	Local       ServerInfo `yaml:"local"`
	Remote      ServerInfo `yaml:"remote"`
	Interval    int        `yaml:"interval"`
}

func readConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
