package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	File string `yaml:"file"`
	Delimiter string `yaml:"delimiter"`
}

type Patch struct {
	File *string
	Delimiter *string
}

func getConfigFromDir(dir string) (*Config, error) {
	path := filepath.Join(dir, ".boltcli")
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}

func mergeConfig(configs ...*Config) *Config {
	result := Config{}
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		if len(cfg.File) > 0 {
			result.File = cfg.File
		}
		if len(cfg.Delimiter) > 0 {
			result.Delimiter = cfg.Delimiter
		}
	}
	return &result
}

func PatchConfig(patch *Patch, global bool) error {
	var configDir string
	var err error

	if global {
		configDir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	} else {
		configDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	fmt.Println(configDir)
	cfg, err := getConfigFromDir(configDir)
	if err != nil {
		return err
	}

	if cfg == nil {
		cfg = &Config{}
	}
	if patch.File != nil {
		cfg.File = *patch.File
	}
	if patch.Delimiter != nil {
		cfg.Delimiter = *patch.Delimiter
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(configDir, ".boltcli"), data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetConfigGlobal() (*Config, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfg, err := getConfigFromDir(configDir)
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		return cfg, nil
	}

	return getDefaultConfig(), nil
}

func GetConfig() (*Config, error) {
	configDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	local, err := getConfigFromDir(configDir)
	if err != nil {
		return nil, err
	}

	configDir, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	global, err := getConfigFromDir(configDir)
	if err != nil {
		return nil, err
	}

	return mergeConfig(getDefaultConfig(), global, local), nil
}

func getDefaultConfig() *Config {
	return &Config{
		Delimiter: "/",
		File: "",
	}
}
