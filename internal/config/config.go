// Package config defines the config structure and operations for summon
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Hotkey struct {
	Keys string `toml:"keys"`
	App  string `toml:"app"`
}

type Config struct {
	Hotkeys []Hotkey `toml:"binding"`
}
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// first run, no config yet
			if err = createDefaultConfig(); err != nil {
				return nil, err
			}
			// decode again
			if _, err = toml.DecodeFile(path, &cfg); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &cfg, nil
}

func createDefaultConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(&Config{
		Hotkeys: []Hotkey{
			Hotkey{Keys: "option+f", App: "Finder"},
		},
	})
}

func AppendBinding(keys, app string) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\n[[binding]]\nkeys = %q\napp = %q\n", keys, app)
	return err
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/.config/summon/config.toml", nil
}
