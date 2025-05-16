package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotpilot/utils"
)

// ApplyConfigurations applies all configurations based on the environment
func ApplyConfigurations(dotpilotDir, environment string) error {
	return ApplyConfigurationsWithOptions(dotpilotDir, environment, true, true)
}

// ApplyConfigurationsWithOptions applies all configurations with specified options
func ApplyConfigurationsWithOptions(dotpilotDir, environment string, backup, diffPrompt bool) error {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Apply configurations in order:
	// 1. Common
	// 2. Environment-specific
	// 3. Machine-specific

	// 1. Apply common configurations
	commonDir := filepath.Join(dotpilotDir, "common")
	if err := applyConfigDir(commonDir, backup, diffPrompt); err != nil {
		return err
	}

	// 2. Apply environment-specific configurations
	if environment != "" {
		envDir := filepath.Join(dotpilotDir, "envs", environment)
		if err := applyConfigDir(envDir, backup, diffPrompt); err != nil {
			return err
		}
	}

	// 3. Apply machine-specific configurations
	machineDir := filepath.Join(dotpilotDir, "machine", hostname)
	if err := applyConfigDir(machineDir, backup, diffPrompt); err != nil {
		return err
	}

	return nil
}

// applyConfigDir applies configurations from a specific directory
func applyConfigDir(configDir string, backup, diffPrompt bool) error {
	// Check if directory exists
	_, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		utils.Logger.Debug().Msgf("Configuration directory does not exist: %s", configDir)
		return nil
	}

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Walk through the configuration directory
	return filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == configDir {
			return nil
		}

		// Get the relative path from the configuration directory
		relPath, err := filepath.Rel(configDir, path)
		if err != nil {
			return err
		}

		// Skip special directories and files
		if strings.HasPrefix(relPath, ".git") {
			return nil
		}
		if relPath == "README.md" {
			return nil
		}

		// Construct the target path in the home directory
		targetPath := filepath.Join(home, relPath)

		// Handle directory
		if info.IsDir() {
			if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		// Check if target already exists and is not a symlink to our path
		targetInfo, err := os.Lstat(targetPath)
		if err == nil {
			isSymlink := targetInfo.Mode()&os.ModeSymlink != 0
			
			if isSymlink {
				// Check if symlink points to our dotpilot path
				linkTarget, err := os.Readlink(targetPath)
				if err == nil && linkTarget == path {
					utils.Logger.Debug().Msgf("Symlink already exists: %s -> %s", targetPath, path)
					return nil
				}
			}

			// It exists but isn't a correct symlink, prompt for diff if needed
			if diffPrompt {
				if _, err := os.Stat(targetPath); err == nil {
					diff, err := FileDiff(targetPath, path)
					if err != nil {
						utils.Logger.Warn().Err(err).Msgf("Failed to get diff for %s", targetPath)
					} else {
						fmt.Printf("Diff for %s:\n%s\n", targetPath, diff)
						
						if !utils.PromptYesNo(fmt.Sprintf("Apply changes to %s?", targetPath)) {
							utils.Logger.Info().Msgf("Skipping %s", targetPath)
							return nil
						}
					}
				}
			}

			// Backup if requested
			if backup {
				backupPath, err := BackupFile(targetPath)
				if err != nil {
					utils.Logger.Warn().Err(err).Msgf("Failed to backup %s", targetPath)
				} else if backupPath != "" {
					utils.Logger.Info().Msgf("Backed up %s to %s", targetPath, backupPath)
				}
			}

			// Remove the target if it exists
			if err := os.Remove(targetPath); err != nil {
				return err
			}
		}

		// Create symlink
		utils.Logger.Debug().Msgf("Creating symlink: %s -> %s", targetPath, path)
		if err := os.Symlink(path, targetPath); err != nil {
			return err
		}

		// Update tracking list
		relTarget, err := filepath.Rel(home, targetPath)
		if err == nil {
			AddTrackingPath(relTarget)
		}

		return nil
	})
}
