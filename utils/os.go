package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OSInfo contains information about the operating system
type OSInfo struct {
	Name           string
	Version        string
	PackageManager string
}

// GetOSInfo returns information about the current operating system
func GetOSInfo() OSInfo {
	info := OSInfo{
		Name:           runtime.GOOS,
		Version:        "",
		PackageManager: "",
	}

	switch info.Name {
	case "darwin":
		info.Name = "macOS"
		info.PackageManager = "brew"
		if _, err := exec.LookPath("brew"); err != nil {
			Logger.Warn().Msg("Homebrew not found, some features may not work")
		}
	case "linux":
		// Detect Linux distribution
		info = detectLinuxDistro()
	case "windows":
		info.Name = "Windows"
		// No default package manager
	}

	return info
}

// detectLinuxDistro detects the Linux distribution and package manager
func detectLinuxDistro() OSInfo {
	info := OSInfo{
		Name:           "Linux",
		Version:        "",
		PackageManager: "",
	}

	// Check if /etc/os-release exists
	if _, err := os.Stat("/etc/os-release"); err == nil {
		data, err := os.ReadFile("/etc/os-release")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ID=") {
					info.Name = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
				} else if strings.HasPrefix(line, "VERSION_ID=") {
					info.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
				}
			}
		}
	}

	// Detect package manager
	switch info.Name {
	case "ubuntu", "debian", "pop", "elementary", "mint":
		info.PackageManager = "apt"
	case "fedora", "centos", "rhel":
		info.PackageManager = "dnf"
	case "arch", "manjaro", "endeavouros":
		// Check if yay is installed
		if _, err := exec.LookPath("yay"); err == nil {
			info.PackageManager = "yay"
		} else {
			info.PackageManager = "pacman"
		}
	case "opensuse", "suse":
		info.PackageManager = "zypper"
	default:
		// Try to detect package manager based on available commands
		for _, pm := range []string{"apt", "dnf", "pacman", "yay", "zypper"} {
			if _, err := exec.LookPath(pm); err == nil {
				info.PackageManager = pm
				break
			}
		}
	}

	return info
}
