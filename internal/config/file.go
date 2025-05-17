package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SyncFile string `yaml:"syncfile"`
}

type ConfigFile struct {
	path   string
	Config Config
}

func (c *ConfigFile) Close() error {
	d, err := yaml.Marshal(&c.Config)
	if err != nil {
		return err
	}

	err = os.WriteFile(c.path, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func OpenConfigFile(path string) (*ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		slog.Info("config file does not exist", "path", path)
		return &ConfigFile{
			path,
			Config{
				SyncFile: "",
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &ConfigFile{
		path,
		config,
	}, nil
}
