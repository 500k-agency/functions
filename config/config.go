package config

import (
	"fmt"
	"os"

	"github.com/500k-agency/function/lib/connect"
	"github.com/500k-agency/function/product"

	"github.com/BurntSushi/toml"
)

// Config holds all the configuration fields needed within the application
type Config struct {
	Environment string `toml:"environment"`

	// [connect]
	Connect connect.Configs `toml:"connect"`

	// [products]
	Products []product.Config `toml:"products"`
}

// NewFromSecrets instantiates the config struct from secrets
func NewFromSecrets() (*Config, error) {
	conf := Config{}

	if _, err := os.Stat("/etc/secrets/latest"); os.IsNotExist(err) {
		return nil, err
	}
	if _, err := toml.DecodeFile("/etc/secrets/latest", &conf); err != nil {
		return nil, fmt.Errorf("unable to load config file: %w", err)
	}

	return &conf, nil
}

// NewFromFile instantiates the config struct
func NewFromFile(fileConfig, envConfig string) (*Config, error) {
	file := fileConfig
	if file == "" {
		file = envConfig
	}

	var conf Config
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, fmt.Errorf("unable to load config file: %w", err)
	}

	return &conf, nil
}
