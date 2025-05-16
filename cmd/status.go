package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotpilot/core"
	"github.com/dotpilot/utils"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of dotpilot",
	Long: `Show the current status of the dotpilot repository,
including the current environment, tracked files, and git status.

For example:
  dotpilot status`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get home directory
		home, err := os.UserHomeDir()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to get home directory")
			os.Exit(1)
		}

		// Check if dotpilot is initialized
		dotpilotDir := filepath.Join(home, ".dotpilot")
		if _, err := os.Stat(dotpilotDir); os.IsNotExist(err) {
			utils.Logger.Error().Msg("Dotpilot is not initialized. Run 'dotpilot init' first.")
			os.Exit(1)
		}

		// Get current environment
		cfg := core.GetConfig()
		environment := cfg.CurrentEnvironment
		if environment == "" {
			environment = "default"
		}

		// Get hostname
		hostname, err := os.Hostname()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to get hostname")
			hostname = "unknown"
		}

		// Get OS info
		osInfo := utils.GetOSInfo()

		// Print general status
		fmt.Println("=== DotPilot Status ===")
		fmt.Printf("Current environment: %s\n", environment)
		fmt.Printf("Machine hostname: %s\n", hostname)
		fmt.Printf("Operating system: %s\n", osInfo.Name)
		fmt.Printf("Package system: %s\n", osInfo.PackageManager)
		fmt.Println()

		// Check for uncommitted changes
		hasChanges, err := core.HasUncommittedChanges(dotpilotDir)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to check for uncommitted changes")
			os.Exit(1)
		}

		// Print Git status
		fmt.Println("=== Git Status ===")
		if hasChanges {
			fmt.Println("Repository has uncommitted changes.")
			changes, err := core.GetGitStatus(dotpilotDir)
			if err != nil {
				utils.Logger.Error().Err(err).Msg("Failed to get git status")
			} else {
				fmt.Print(changes)
			}
		} else {
			fmt.Println("Repository is clean, no uncommitted changes.")
		}

		// Get remote status
		behindAhead, err := core.GetRemoteStatus(dotpilotDir)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to get remote status")
		} else {
			if behindAhead.Behind > 0 {
				fmt.Printf("Local is behind remote by %d commits.\n", behindAhead.Behind)
			}
			if behindAhead.Ahead > 0 {
				fmt.Printf("Local is ahead of remote by %d commits.\n", behindAhead.Ahead)
			}
			if behindAhead.Behind == 0 && behindAhead.Ahead == 0 {
				fmt.Println("Local is in sync with remote.")
			}
		}
		fmt.Println()

		// Print tracked files
		fmt.Println("=== Tracked Files ===")
		trackedFiles, err := core.GetTrackedFiles(dotpilotDir)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to get tracked files")
		} else {
			if len(trackedFiles) == 0 {
				fmt.Println("No files are currently tracked.")
			} else {
				for _, file := range trackedFiles {
					fmt.Printf("- %s\n", file)
				}
			}
		}
	},
}

func init() {
	// No additional flags needed for status command
}
