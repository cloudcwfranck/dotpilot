package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotpilot/utils"
)

// InstallPackages installs packages based on the environment and OS
func InstallPackages(dotpilotDir, environment, overridePackageSystem string) error {
	// Get OS info
	osInfo := utils.GetOSInfo()
	packageSystem := osInfo.PackageManager
	
	// Override package system if specified
	if overridePackageSystem != "" {
		packageSystem = overridePackageSystem
	}

	utils.Logger.Info().Msgf("Detected OS: %s, Package System: %s", osInfo.Name, packageSystem)

	// Get package files in order:
	// 1. Common
	// 2. Environment-specific
	// 3. Machine-specific

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Define package file names based on package system
	var packageFiles []string
	switch packageSystem {
	case "apt":
		packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "common", "packages.apt"))
		if environment != "" {
			packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "envs", environment, "packages.apt"))
		}
		packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "machine", hostname, "packages.apt"))
	case "brew":
		packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "common", "packages.brew"))
		if environment != "" {
			packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "envs", environment, "packages.brew"))
		}
		packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "machine", hostname, "packages.brew"))
	case "yay":
		packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "common", "packages.yay"))
		if environment != "" {
			packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "envs", environment, "packages.yay"))
		}
		packageFiles = append(packageFiles, filepath.Join(dotpilotDir, "machine", hostname, "packages.yay"))
	default:
		return fmt.Errorf("unsupported package system: %s", packageSystem)
	}

	// Read package files and install packages
	for _, packageFile := range packageFiles {
		if err := installPackagesFromFile(packageFile, packageSystem); err != nil {
			return err
		}
	}

	return nil
}

// installPackagesFromFile installs packages from a file
func installPackagesFromFile(packageFile, packageSystem string) error {
	// Check if package file exists
	if _, err := os.Stat(packageFile); os.IsNotExist(err) {
		utils.Logger.Debug().Msgf("Package file does not exist: %s", packageFile)
		return nil
	}

	// Read package file
	data, err := os.ReadFile(packageFile)
	if err != nil {
		return err
	}

	// Parse packages
	lines := strings.Split(string(data), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		packages = append(packages, line)
	}

	if len(packages) == 0 {
		utils.Logger.Debug().Msgf("No packages to install from %s", packageFile)
		return nil
	}

	utils.Logger.Info().Msgf("Installing %d packages from %s", len(packages), packageFile)

	// Build installation command
	var cmd string
	var args []string
	switch packageSystem {
	case "apt":
		cmd = "apt-get"
		args = append([]string{"install", "-y"}, packages...)
	case "brew":
		cmd = "brew"
		args = append([]string{"install"}, packages...)
	case "yay":
		cmd = "yay"
		args = append([]string{"-S", "--noconfirm"}, packages...)
	default:
		return fmt.Errorf("unsupported package system: %s", packageSystem)
	}

	// Run installation command
	output, err := utils.ExecuteCommand(cmd, args...)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("Failed to install packages: %s", output)
		return err
	}

	utils.Logger.Info().Msgf("Successfully installed packages from %s", packageFile)
	return nil
}
