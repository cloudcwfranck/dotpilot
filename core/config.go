package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dotpilot/utils"
)

// Config represents the configuration of dotpilot
type Config struct {
	RemoteRepository   string                 `json:"remote_repository"`
	CurrentEnvironment string                 `json:"current_environment"`
	TrackingPaths      []string               `json:"tracking_paths"`
	Options            map[string]interface{} `json:"options"`
}

var currentConfig Config

// LoadConfig loads the configuration from the file
func LoadConfig(configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &currentConfig)
	if err != nil {
		return err
	}

	utils.Logger.Debug().Msgf("Loaded config from %s", configPath)
	return nil
}

// SaveConfig saves the current configuration to the file
func SaveConfig(configPath string) error {
	data, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	utils.Logger.Debug().Msgf("Saved config to %s", configPath)
	return nil
}

// GetConfig returns the current configuration
func GetConfig() Config {
	return currentConfig
}

// SetConfig sets the current configuration
func SetConfig(config Config) {
	currentConfig = config
}

// InitDefaultConfig initializes a default configuration
func InitDefaultConfig() {
	currentConfig = Config{
		RemoteRepository:   "",
		CurrentEnvironment: "default",
		TrackingPaths:      []string{},
		Options: map[string]interface{}{
			"backup_before_overwrite": true,
			"prompt_on_diff":          true,
		},
	}
}

// CreateDefaultConfigFile creates a default configuration file
func CreateDefaultConfigFile(remoteRepo, environment string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".dotpilotrc")
	
	// Initialize config
	currentConfig = Config{
		RemoteRepository:   remoteRepo,
		CurrentEnvironment: environment,
		TrackingPaths:      []string{},
		Options: map[string]interface{}{
			"backup_before_overwrite": true,
			"prompt_on_diff":          true,
		},
	}

	// Save config
	return SaveConfig(configPath)
}

// UpdateEnvironment updates the current environment in the configuration
func UpdateEnvironment(environment string) error {
	currentConfig.CurrentEnvironment = environment

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".dotpilotrc")
	return SaveConfig(configPath)
}

// AddTrackingPath adds a path to the tracked paths list
func AddTrackingPath(path string) error {
	// Check if the path is already tracked
	for _, p := range currentConfig.TrackingPaths {
		if p == path {
			return nil
		}
	}

	currentConfig.TrackingPaths = append(currentConfig.TrackingPaths, path)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".dotpilotrc")
	return SaveConfig(configPath)
}
