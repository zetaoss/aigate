package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func normalizeLogLevel(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

type Config struct {
	Server    ServerConfig              `yaml:"server"`
	Models    []string                  `yaml:"models"`
	Providers map[string]ProviderConfig `yaml:"providers"`
	Fallback  FallbackConfig            `yaml:"fallback"`
}

type ServerConfig struct {
	Port                    int    `yaml:"port"`
	LogLevel                string `yaml:"logLevel"`
	ValidateModelsOnStartup *bool  `yaml:"validateModelsOnStartup"`
}

type ProviderConfig struct {
	APIKey string `yaml:"apiKey"`
}

type FallbackConfig struct {
	MaxFallbacks      *int          `yaml:"maxFallbacks"`
	RetryOnCodes      []int         `yaml:"retryOnCodes"`
	PerAttemptTimeout time.Duration `yaml:"perAttemptTimeout"`
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file %q: %w", path, err)
	}

	if len(cfg.Models) == 0 {
		return nil, fmt.Errorf("config must define at least one model")
	}

	if len(cfg.Providers) == 0 {
		return nil, fmt.Errorf("config must define at least one provider")
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return nil, fmt.Errorf("config server.port must be between 1 and 65535")
	}
	if cfg.Server.LogLevel == "" {
		cfg.Server.LogLevel = "info"
	}
	switch normalized := normalizeLogLevel(cfg.Server.LogLevel); normalized {
	case "debug", "info", "error":
		cfg.Server.LogLevel = normalized
	default:
		return nil, fmt.Errorf("config server.logLevel must be one of: debug, info, error")
	}
	if cfg.Server.ValidateModelsOnStartup == nil {
		defaultValue := true
		cfg.Server.ValidateModelsOnStartup = &defaultValue
	}
	if cfg.Fallback.MaxFallbacks == nil {
		defaultFallbacks := -1
		cfg.Fallback.MaxFallbacks = &defaultFallbacks
	}
	if *cfg.Fallback.MaxFallbacks < -1 {
		return nil, fmt.Errorf("config fallback.maxFallbacks must be >= -1")
	}
	if len(cfg.Fallback.RetryOnCodes) == 0 {
		cfg.Fallback.RetryOnCodes = []int{429, 500, 502, 503, 504}
	}
	for _, code := range cfg.Fallback.RetryOnCodes {
		if code < 100 || code > 599 {
			return nil, fmt.Errorf("config fallback.retryOnCodes must be valid HTTP status codes")
		}
	}
	if cfg.Fallback.PerAttemptTimeout == 0 {
		cfg.Fallback.PerAttemptTimeout = 30 * time.Second
	}
	if cfg.Fallback.PerAttemptTimeout < 0 {
		return nil, fmt.Errorf("config fallback.perAttemptTimeout must be >= 0")
	}

	for providerName := range cfg.Providers {
		if providerName == "" {
			return nil, fmt.Errorf("provider names cannot be empty")
		}
	}

	return &cfg, nil
}

func RedactedYAML(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("config is nil")
	}

	redacted := *cfg
	redacted.Providers = make(map[string]ProviderConfig, len(cfg.Providers))
	for name, provider := range cfg.Providers {
		copyProvider := provider
		if strings.TrimSpace(copyProvider.APIKey) != "" {
			copyProvider.APIKey = "[redacted]"
		}
		redacted.Providers[name] = copyProvider
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(&redacted)
	_ = enc.Close()
	if err != nil {
		return "", fmt.Errorf("marshal redacted config: %w", err)
	}
	return buf.String(), nil
}
