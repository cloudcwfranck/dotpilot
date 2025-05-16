package core

import (
	"os"
	"path/filepath"

	"github.com/dotpilot/utils"
)

// RunHooks runs hooks based on the environment
func RunHooks(dotpilotDir, environment, hookName string) error {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Define hook files in order:
	// 1. Common
	// 2. Environment-specific
	// 3. Machine-specific

	var hookFiles []string
	hookFiles = append(hookFiles, filepath.Join(dotpilotDir, "common", hookName))
	if environment != "" {
		hookFiles = append(hookFiles, filepath.Join(dotpilotDir, "envs", environment, hookName))
	}
	hookFiles = append(hookFiles, filepath.Join(dotpilotDir, "machine", hostname, hookName))

	// Run hooks
	for _, hookFile := range hookFiles {
		if err := runHook(hookFile); err != nil {
			return err
		}
	}

	return nil
}

// runHook runs a single hook script
func runHook(hookFile string) error {
	// Check if hook file exists
	if _, err := os.Stat(hookFile); os.IsNotExist(err) {
		utils.Logger.Debug().Msgf("Hook file does not exist: %s", hookFile)
		return nil
	}

	// Make hook file executable
	if err := os.Chmod(hookFile, 0755); err != nil {
		return err
	}

	// Execute hook
	utils.Logger.Info().Msgf("Running hook: %s", hookFile)
	output, err := utils.ExecuteCommand(hookFile)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("Hook failed: %s\nOutput: %s", hookFile, output)
		return err
	}

	utils.Logger.Info().Msgf("Hook succeeded: %s", hookFile)
	return nil
}
