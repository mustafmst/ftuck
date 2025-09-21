package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SyncFile string `yaml:"syncfile"`
}

func (c *Config) GetSyncFilePath() string {
	return c.SyncFile
}

type ConfigFile struct {
	path   string
	Config Config
}

func (c *ConfigFile) Save() error {
	slog.Debug("saving config file", "path", c.path)
	
	d, err := yaml.Marshal(&c.Config)
	if err != nil {
		slog.Error("failed to marshal config", "path", c.path, "error", err)
		return err
	}

	err = os.WriteFile(c.path, d, 0644)
	if err != nil {
		slog.Error("failed to write config file", "path", c.path, "error", err)
		return err
	}

	slog.Debug("config file saved successfully", "path", c.path)
	return nil
}

func OpenConfigFile(path string) (*ConfigFile, error) {
	slog.Debug("attempting to open config file", "path", path)
	
	data, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		slog.Info("config file does not exist, creating new one", "path", path)
		return &ConfigFile{
			path,
			Config{
				SyncFile: "",
			},
		}, nil
	}
	if err != nil {
		slog.Error("failed to read config file", "path", path, "error", err)
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error("failed to parse config file", "path", path, "error", err)
		return nil, err
	}

	slog.Debug("config file loaded successfully", "path", path, "sync_file", config.SyncFile)
	return &ConfigFile{
		path,
		config,
	}, nil
}
