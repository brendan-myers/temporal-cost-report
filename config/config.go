package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TemporalConfig represents the structure of temporal.yaml
type TemporalConfig struct {
	Env map[string]Environment `yaml:"env"`
}

// Environment represents a named environment configuration
type Environment struct {
	Address     string `yaml:"address"`
	Namespace   string `yaml:"namespace"`
	TLSCertPath string `yaml:"tls-cert-path"`
	TLSKeyPath  string `yaml:"tls-key-path"`
	APIKey      string `yaml:"api-key"`
}

// ConnectionConfig holds all connection parameters
type ConnectionConfig struct {
	Address     string
	Namespace   string
	TLSCertPath string
	TLSKeyPath  string
	APIKey      string
}

// DefaultConfigPath returns the default temporal.yaml path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "temporalio", "temporal.yaml")
}

// LoadEnvironment loads a named environment from the temporal.yaml file
func LoadEnvironment(envName string) (*Environment, error) {
	configPath := DefaultConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("temporal config file not found at %s", configPath)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config TemporalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse temporal.yaml: %w", err)
	}

	env, exists := config.Env[envName]
	if !exists {
		available := make([]string, 0, len(config.Env))
		for name := range config.Env {
			available = append(available, name)
		}
		return nil, fmt.Errorf("environment '%s' not found in temporal.yaml; available: %v", envName, available)
	}

	return &env, nil
}

// MergeConfig merges environment config with command-line overrides.
// Command-line flags take precedence over environment values.
func MergeConfig(env *Environment, flags ConnectionConfig) ConnectionConfig {
	result := ConnectionConfig{}

	// Start with environment values if present
	if env != nil {
		result.Address = env.Address
		result.Namespace = env.Namespace
		result.TLSCertPath = env.TLSCertPath
		result.TLSKeyPath = env.TLSKeyPath
		result.APIKey = env.APIKey
	}

	// Override with explicit flag values (non-empty strings)
	if flags.Address != "" {
		result.Address = flags.Address
	}
	if flags.Namespace != "" {
		result.Namespace = flags.Namespace
	}
	if flags.TLSCertPath != "" {
		result.TLSCertPath = flags.TLSCertPath
	}
	if flags.TLSKeyPath != "" {
		result.TLSKeyPath = flags.TLSKeyPath
	}
	if flags.APIKey != "" {
		result.APIKey = flags.APIKey
	}

	return result
}

// Validate checks that required fields are present and authentication is configured
func (c *ConnectionConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required: use --address flag or --env to load from temporal.yaml")
	}
	if c.Namespace == "" {
		return fmt.Errorf("namespace is required: use --namespace flag or --env to load from temporal.yaml")
	}

	hasTLS := c.TLSCertPath != "" && c.TLSKeyPath != ""
	hasAPIKey := c.APIKey != "" || os.Getenv("TEMPORAL_API_KEY") != ""

	if !hasTLS && !hasAPIKey {
		return fmt.Errorf("authentication required: provide --tls-cert-path and --tls-key-path for mTLS, or --api-key (or TEMPORAL_API_KEY env var) for API key auth")
	}

	// Warn if only one TLS path is provided
	if (c.TLSCertPath != "" && c.TLSKeyPath == "") || (c.TLSCertPath == "" && c.TLSKeyPath != "") {
		return fmt.Errorf("both --tls-cert-path and --tls-key-path must be provided together")
	}

	return nil
}

// UsesMTLS returns true if mTLS authentication should be used
func (c *ConnectionConfig) UsesMTLS() bool {
	return c.TLSCertPath != "" && c.TLSKeyPath != ""
}
