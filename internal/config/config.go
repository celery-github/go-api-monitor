package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Name           string `yaml:"name" json:"name"`
	Method         string `yaml:"method" json:"method"`
	URL            string `yaml:"url" json:"url"`
	ExpectedStatus int    `yaml:"expected_status" json:"expected_status"`
}

type Config struct {
	IntervalSeconds int        `yaml:"interval_seconds"`
	TimeoutSeconds  int        `yaml:"timeout_seconds"`
	Endpoints       []Endpoint `yaml:"endpoints"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 60
	}
	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 5
	}
	if len(cfg.Endpoints) == 0 {
		return nil, errors.New("config must include at least one endpoint")
	}
	for i := range cfg.Endpoints {
		if cfg.Endpoints[i].ExpectedStatus == 0 {
			cfg.Endpoints[i].ExpectedStatus = 200
		}
		if cfg.Endpoints[i].Method == "" {
			cfg.Endpoints[i].Method = "GET"
		}
		if cfg.Endpoints[i].Name == "" {
			cfg.Endpoints[i].Name = cfg.Endpoints[i].URL
		}
	}

	return &cfg, nil
}
