package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var Instance *Config

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

func Load() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("error getting config path: %w", err)
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			Instance = &Config{}
			err = Save()
			if err != nil {
				return fmt.Errorf("error saving config file: %w", err)
			}
			return nil
		}
		return fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		err := configFile.Close()
		if err != nil {
			log.Warn().Err(err).Msg("error closing config file")
		}
	}()

	err = yaml.NewDecoder(configFile).Decode(&Instance)
	if err != nil {
		return fmt.Errorf("error decoding config file: %w", err)
	}

	return nil
}

func Save() error {
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

	err = yaml.NewEncoder(configFile).Encode(Instance)
	if err != nil {
		return fmt.Errorf("error encoding config file: %w", err)
	}

	return nil
}
