package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func GetConfigPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error getting user config dir: %w", err)
	}

	appConfigFolder := filepath.Join(userConfigDir, "repo-archiver")

	if _, err := os.Stat(appConfigFolder); os.IsNotExist(err) {
		log.Debug().Str("path", appConfigFolder).Msg("Config folder does not exist, creating it")
		err = os.Mkdir(appConfigFolder, 0700)
		if err != nil {
			return "", fmt.Errorf("error creating config folder: %w", err)
		}
	}

	return filepath.Join(appConfigFolder, "config.yaml"), nil
}

func Load() (Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting config path: %w", err)
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = Save(Config{})
			if err != nil {
				return Config{}, fmt.Errorf("error saving config file: %w", err)
			}
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		err := configFile.Close()
		if err != nil {
			log.Warn().Err(err).Msg("error closing config file")
		}
	}()

	var config Config
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding config file: %w", err)
	}

	return config, nil
}

func Save(config Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("error getting config path: %w", err)
	}

	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}

	defer func() {
		err := configFile.Close()
		if err != nil {
			log.Warn().Err(err).Msg("error closing config file")
		}
	}()

	err = yaml.NewEncoder(configFile).Encode(config)
	if err != nil {
		return fmt.Errorf("error encoding config file: %w", err)
	}

	return nil
}
