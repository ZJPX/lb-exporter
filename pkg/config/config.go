package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

type Credential struct {
	AccessKey       string `json:"accessKey"`
	AccessKeySecret string `json:"accessKeySecret"`
	Region          string `json:"region"`
	Cloud           string `json:"cloud"`
}

type CronSpec struct {
	Spec string `json:"spec"`
}

// Config exporter config
type Config struct {
	Credentials map[string]Credential `json:"credentials"`
	// todo: add extra labels
	Labels  map[string]string    `json:"labels,omitempty"`
	Metrics map[string][]*Metric `json:"metrics"` // mapping for namespace and metrics
	// InstanceInfos []string             `json:"instanceInfos"`
	Cron      CronSpec `json:"cron"`
	RateLimit int      `json:"rate_limit"`
}

func (c *Config) SetDefaults() {
	for key, _ := range c.Credentials {
		if c.Credentials[key].Region == "" {
			credential := c.Credentials[key]
			credential.Region = "cn-hangzhou"
			c.Credentials[key] = credential
		}
	}
	for _, metrics := range c.Metrics {
		for i := range metrics {
			metrics[i].setDefaults()
		}
	}
}

// Parse parse config from file
func Parse(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	cfg.SetDefaults()
	return &cfg, nil
}
