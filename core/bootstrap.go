package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dotpilot/utils"
)

// ApplyDirectoryConfigs applies all configurations from the given directory
// to the destination directory (typically home directory)
func ApplyDirectoryConfigs(sourceDir, destDir string, forceOverwrite bool) error {
	// Check if the source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	// List all files and directories in the source directory
	entries, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %s: %w", sourceDir, err)
	}

	// Process each entry
	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		
		// Skip hidden files/directories and install_packages.sh (handled separately)
		if entry.Name()[0] == '.' || entry.Name() == "install_packages.sh" {
			continue
		}

		// Determine destination path
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// For directories, recursively apply configurations
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %s: %w", destPath, err)
			}

			if err := ApplyDirectoryConfigs(sourcePath, destPath, forceOverwrite); err != nil {
				return err
			}
		} else {
			// For files, create symlinks
			if err := CreateSymlink(sourcePath, destPath, forceOverwrite); err != nil {
				return fmt.Errorf("failed to create symlink for %s: %w", entry.Name(), err)
			}
			utils.Logger.Debug().Msgf("Created symlink: %s -> %s", destPath, sourcePath)
		}
	}

	return nil
}

// CreateSymlink creates a symlink from source to dest
// If dest already exists and forceOverwrite is true, it will be replaced
// Otherwise, the user will be prompted to confirm the overwrite
func CreateSymlink(source, dest string, forceOverwrite bool) error {
	// Check if destination already exists
	if _, err := os.Stat(dest); err == nil {
		// If forceOverwrite is false, prompt the user
		if !forceOverwrite {
			utils.Logger.Warn().Msgf("File already exists: %s", dest)
			if !PromptYesNo(fmt.Sprintf("Overwrite existing file: %s?", dest)) {
				utils.Logger.Info().Msgf("Skipping %s", dest)
				return nil
			}
		}
		
		// Create a backup of the existing file
		backupPath := dest + ".backup"
		utils.Logger.Debug().Msgf("Creating backup of %s to %s", dest, backupPath)
		if err := os.Rename(dest, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of %s: %w", dest, err)
		}
	}

	// Create parent directory if it doesn't exist
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s: %w", destDir, err)
	}

	// Create the symlink
	utils.Logger.Debug().Msgf("Creating symlink: %s -> %s", dest, source)
	if err := os.Symlink(source, dest); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// RunScript executes the given script with bash
func RunScript(scriptPath string) error {
	utils.Logger.Debug().Msgf("Running script: %s", scriptPath)

	// Make script executable if it's not already
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return fmt.Errorf("failed to make script executable: %w", err)
	}

	// Run the script with bash
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}

	return nil
}

// PromptYesNo asks the user a yes/no question and returns true if the answer is yes
func PromptYesNo(question string) bool {
	var response string
	utils.Logger.Info().Msgf("%s (y/n): ", question)
	fmt.Scanln(&response)
	return response == "y" || response == "Y" || response == "yes" || response == "Yes"
}