// Package config loads and validates the open-shark TOML configuration.
package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/Guilherme-Costa-Dev/open-shark/internal/protocol"
)

const stages = 6

// Config mirrors config.example.toml.
type Config struct {
	DPI     DPI     `toml:"dpi"`
	Polling Polling `toml:"polling"`
}

type DPI struct {
	Stages []int `toml:"stages"`
	Active int   `toml:"active"`
}

type Polling struct {
	Rate int `toml:"rate"`
}

// Load reads and validates the configuration at path.
func Load(path string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(path, &c); err != nil {
		return c, fmt.Errorf("read %s: %w", path, err)
	}
	if err := c.Validate(); err != nil {
		return c, fmt.Errorf("%s: %w", path, err)
	}
	return c, nil
}

// Validate checks every field against the values the mouse accepts.
func (c Config) Validate() error {
	if len(c.DPI.Stages) != stages {
		return fmt.Errorf("dpi.stages must have exactly %d values, got %d", stages, len(c.DPI.Stages))
	}
	for i, d := range c.DPI.Stages {
		if d < protocol.MinDPI || d > protocol.MaxDPI {
			return fmt.Errorf("dpi.stages[%d] = %d is out of range %d..%d", i, d, protocol.MinDPI, protocol.MaxDPI)
		}
	}
	if c.DPI.Active < 1 || c.DPI.Active > stages {
		return fmt.Errorf("dpi.active = %d is out of range 1..%d", c.DPI.Active, stages)
	}

	switch c.Polling.Rate {
	case 125, 250, 500, 1000:
	default:
		return fmt.Errorf("polling.rate = %d is unsupported (use 125, 250, 500 or 1000)", c.Polling.Rate)
	}
	return nil
}

// DefaultPath returns the config path. When the command is run through sudo,
// $HOME points at root's home directory, so the home of the user who invoked
// sudo is preferred instead.
func DefaultPath() string {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		if u, err := user.Lookup(sudoUser); err == nil && u.HomeDir != "" {
			return filepath.Join(u.HomeDir, ".config", "open-shark", "config.toml")
		}
	}
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "open-shark", "config.toml")
	}
	return "config.toml"
}
