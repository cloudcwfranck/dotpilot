package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotpilot/utils"
)

// TrackFile tracks a file or directory in dotpilot
func TrackFile(source, destination, dotpilotDir string, overwrite bool) error {
	// Check if source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create destination directory
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Check if destination already exists
	_, err = os.Stat(destination)
	if err == nil && !overwrite {
		return fmt.Errorf("destination already exists: %s", destination)
	}

	// Handle directory
	if sourceInfo.IsDir() {
		return trackDirectory(source, destination, overwrite)
	}

	// Handle file
	return trackSingleFile(source, destination, overwrite)
}

// trackDirectory tracks a directory and its contents
func trackDirectory(source, destination string, overwrite bool) error {
	// Create destination directory
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	// Walk through the source directory
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path from the source directory
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// Skip the root directory
		if relPath == "." {
			return nil
		}

		// Construct the destination path
		destPath := filepath.Join(destination, relPath)

		// Handle directory
		if info.IsDir() {
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		// Handle file
		return trackSingleFile(path, destPath, overwrite)
	})
}

// trackSingleFile tracks a single file
func trackSingleFile(source, destination string, overwrite bool) error {
	// Get source info
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Check if destination already exists
	_, err = os.Stat(destination)
	if err == nil && !overwrite {
		return fmt.Errorf("destination already exists: %s", destination)
	}

	// Create destination directory
	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Copy file
	if err := copyFile(source, destination, sourceInfo.Mode()); err != nil {
		return err
	}

	// Create symlink, first backup existing file if necessary
	linkSource := destination
	linkDest := source

	// Check if source is already a symlink
	linkInfo, err := os.Lstat(source)
	if err == nil && linkInfo.Mode()&os.ModeSymlink != 0 {
		// If it's already a symlink, check if it points to our destination
		linkTarget, err := os.Readlink(source)
		if err == nil && linkTarget == destination {
			utils.Logger.Debug().Msgf("Symlink already exists: %s -> %s", source, destination)
			return nil
		}
	}

	// Backup existing file if it's not already a symlink to our destination
	if err == nil && linkInfo.Mode()&os.ModeSymlink == 0 {
		backupPath := source + ".dotpilot.bak." + time.Now().Format("20060102150405")
		utils.Logger.Debug().Msgf("Backing up %s to %s", source, backupPath)
		if err := os.Rename(source, backupPath); err != nil {
			return err
		}
	}

	// Create symlink
	utils.Logger.Debug().Msgf("Creating symlink: %s -> %s", linkDest, linkSource)
	if err := os.Symlink(linkSource, linkDest); err != nil {
		return err
	}

	// Update tracking list
	relSource, err := filepath.Rel(os.Getenv("HOME"), source)
	if err == nil {
		AddTrackingPath(relSource)
	}

	return nil
}

// copyFile copies a file from source to destination
func copyFile(source, destination string, mode os.FileMode) error {
	// Open source file
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, sourceFile)
	return err
}

// BackupFile creates a backup of a file
func BackupFile(path string) (string, error) {
	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", nil
	}

	// Create backup path
	backupPath := path + ".dotpilot.bak." + time.Now().Format("20060102150405")
	
	// Copy file
	sourceInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	err = copyFile(path, backupPath, sourceInfo.Mode())
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

// FileDiff returns the diff between two files
func FileDiff(file1, file2 string) (string, error) {
	// Read files
	content1, err := ioutil.ReadFile(file1)
	if err != nil {
		return "", err
	}

	content2, err := ioutil.ReadFile(file2)
	if err != nil {
		return "", err
	}

	// Compare line by line
	lines1 := strings.Split(string(content1), "\n")
	lines2 := strings.Split(string(content2), "\n")

	diff := ""
	maxLines := len(lines1)
	if len(lines2) > maxLines {
		maxLines = len(lines2)
	}

	for i := 0; i < maxLines; i++ {
		if i >= len(lines1) {
			diff += fmt.Sprintf("+ %s\n", lines2[i])
		} else if i >= len(lines2) {
			diff += fmt.Sprintf("- %s\n", lines1[i])
		} else if lines1[i] != lines2[i] {
			diff += fmt.Sprintf("- %s\n+ %s\n", lines1[i], lines2[i])
		}
	}

	if diff == "" {
		return "Files are identical", nil
	}

	return diff, nil
}
