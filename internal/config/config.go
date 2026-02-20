package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	goconfy "github.com/keksclan/goConfy"
)

// Config holds the application configuration loaded via goConfy (YAML + ENV macros).
type Config struct {
	Server ServerConfig `yaml:"server"`
	Auth   AuthConfig   `yaml:"auth"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// AuthConfig holds authentication/authorization settings.
type AuthConfig struct {
	Issuer           string `yaml:"issuer"`
	Audience         string `yaml:"audience"`
	JWKSURL          string `yaml:"jwks_url"`
	IntrospectionURL string `yaml:"introspection_url"`
	ClientID         string `yaml:"client_id"`
	ClientSecret     string `yaml:"client_secret"`

	// RequiredHeadersRaw is the comma-separated string from YAML/ENV.
	// Parsed into RequiredHeaders during Normalize().
	RequiredHeadersRaw string   `yaml:"required_headers"`
	RequiredHeaders    []string `yaml:"-"`
}

// Normalize implements goconfy.Normalizable – called automatically after decoding.
func (c *Config) Normalize() {
	// Parse required headers CSV into slice.
	if c.Auth.RequiredHeadersRaw != "" {
		parts := strings.Split(c.Auth.RequiredHeadersRaw, ",")
		headers := make([]string, 0, len(parts))
		for _, h := range parts {
			h = strings.TrimSpace(h)
			if h != "" {
				headers = append(headers, h)
			}
		}
		c.Auth.RequiredHeaders = headers
	}
}

// Validate implements goconfy.Validatable – called automatically after Normalize().
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.ReadTimeout == 0 {
		c.Server.ReadTimeout = 10 * time.Second
	}
	if c.Server.WriteTimeout == 0 {
		c.Server.WriteTimeout = 10 * time.Second
	}
	if c.Server.IdleTimeout == 0 {
		c.Server.IdleTimeout = 60 * time.Second
	}
	return nil
}

// Load reads config via goConfy from YAML file with ENV macro expansion.
// The YAML file path is determined by: configPath argument > CONFIG_PATH env > "config.yml".
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "config.yml"
	}

	cfg, err := goconfy.Load[Config](
		goconfy.WithFile(configPath),
		goconfy.WithStrictYAML(false),
	)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &cfg, nil
}

// AudienceList returns the parsed audience list.
// If Audience is "*", returns nil (any audience allowed).
// Otherwise splits by comma and trims spaces.
func (c *AuthConfig) AudienceList() []string {
	if c.Audience == "*" || c.Audience == "" {
		return nil
	}
	parts := strings.Split(c.Audience, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

// AudienceIsWildcard returns true if audience is set to "*".
func (c *AuthConfig) AudienceIsWildcard() bool {
	return c.Audience == "*"
}

// HasIntrospection returns true if introspection URL is configured.
func (c *AuthConfig) HasIntrospection() bool {
	return c.IntrospectionURL != ""
}

// HasJWKS returns true if JWKS URL is configured.
func (c *AuthConfig) HasJWKS() bool {
	return c.JWKSURL != ""
}
