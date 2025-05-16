package cmd

import (
	"os"
	"path/filepath"

	"github.com/dotpilot/core"
	"github.com/dotpilot/utils"
	"github.com/spf13/cobra"
)

var (
	skipCommon    bool
	skipEnv       bool
	skipMachine   bool
	skipSetupScripts bool
	forceOverwrite bool
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Apply dotfiles and run setup scripts",
	Long: `Bootstrap applies dotfiles from common/, envs/<env>/, and machine/<hostname>/, 
then runs any setup scripts like install_packages.sh.

This command is typically used when setting up a new machine or after significant changes.

For example:
  dotpilot bootstrap
  dotpilot bootstrap --skip-setup-scripts
  dotpilot bootstrap --force`,
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

		// Get hostname for machine-specific configurations
		hostname, err := os.Hostname()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to get hostname")
			hostname = "unknown"
		}

		// Get current environment
		cfg := core.GetConfig()
		environment := cfg.CurrentEnvironment
		if environment == "" {
			environment = "default"
		}

		// Initialize operation manager for progress tracking
		operationManager := utils.NewOperationManager()

		// Apply configurations from different sources
		utils.Logger.Info().Msg("Starting bootstrap process...")

		// 1. Apply common configurations
		if !skipCommon {
			commonOp := operationManager.AddOperation("common", "Applying common dotfiles...", utils.Bar)
			commonOp.Start()

			commonDir := filepath.Join(dotpilotDir, "common")
			if _, err := os.Stat(commonDir); os.IsNotExist(err) {
				utils.Logger.Info().Msg("No common directory found, creating...")
				if err := os.MkdirAll(commonDir, 0755); err != nil {
					commonOp.Stop()
					utils.Logger.Error().Err(err).Msg("Failed to create common directory")
					os.Exit(1)
				}
			}

			if err := core.ApplyDirectoryConfigs(commonDir, home, forceOverwrite); err != nil {
				commonOp.Stop()
				utils.Logger.Error().Err(err).Msg("Failed to apply common configurations")
				os.Exit(1)
			}
			
			commonOp.SetState(utils.StateSuccess)
			commonOp.Stop()
		}

		// 2. Apply environment-specific configurations
		if !skipEnv && environment != "default" {
			envOp := operationManager.AddOperation("env", "Applying environment-specific dotfiles...", utils.Bar)
			envOp.Start()

			envDir := filepath.Join(dotpilotDir, "envs", environment)
			if _, err := os.Stat(envDir); os.IsNotExist(err) {
				utils.Logger.Info().Msgf("No configuration for environment '%s' found, creating...", environment)
				if err := os.MkdirAll(envDir, 0755); err != nil {
					envOp.Stop()
					utils.Logger.Error().Err(err).Msg("Failed to create environment directory")
					os.Exit(1)
				}
				envOp.SetState(utils.StateInfo)
				envOp.Stop()
			} else {
				if err := core.ApplyDirectoryConfigs(envDir, home, forceOverwrite); err != nil {
					envOp.Stop()
					utils.Logger.Error().Err(err).Msg("Failed to apply environment-specific configurations")
					os.Exit(1)
				}
				envOp.SetState(utils.StateSuccess)
				envOp.Stop()
			}
		}

		// 3. Apply machine-specific configurations
		if !skipMachine {
			machineOp := operationManager.AddOperation("machine", "Applying machine-specific dotfiles...", utils.Bar)
			machineOp.Start()

			machineDir := filepath.Join(dotpilotDir, "machine", hostname)
			if _, err := os.Stat(machineDir); os.IsNotExist(err) {
				utils.Logger.Info().Msgf("No configuration for machine '%s' found, creating...", hostname)
				if err := os.MkdirAll(machineDir, 0755); err != nil {
					machineOp.Stop()
					utils.Logger.Error().Err(err).Msg("Failed to create machine directory")
					os.Exit(1)
				}
				machineOp.SetState(utils.StateInfo)
				machineOp.Stop()
			} else {
				if err := core.ApplyDirectoryConfigs(machineDir, home, forceOverwrite); err != nil {
					machineOp.Stop()
					utils.Logger.Error().Err(err).Msg("Failed to apply machine-specific configurations")
					os.Exit(1)
				}
				machineOp.SetState(utils.StateSuccess)
				machineOp.Stop()
			}
		}

		// 4. Run setup scripts
		if !skipSetupScripts {
			scriptsOp := operationManager.AddOperation("scripts", "Running setup scripts...", utils.Pulse)
			scriptsOp.Start()

			// Run common setup scripts
			if !skipCommon {
				commonScriptPath := filepath.Join(dotpilotDir, "common", "install_packages.sh")
				if _, err := os.Stat(commonScriptPath); err == nil {
					utils.Logger.Info().Msg("Running common setup script...")
					if err := core.RunScript(commonScriptPath); err != nil {
						scriptsOp.SetState(utils.StateWarning)
						utils.Logger.Warn().Err(err).Msg("Error running common setup script")
						// Continue anyway
					}
				}
			}

			// Run environment-specific setup scripts
			if !skipEnv && environment != "default" {
				envScriptPath := filepath.Join(dotpilotDir, "envs", environment, "install_packages.sh")
				if _, err := os.Stat(envScriptPath); err == nil {
					utils.Logger.Info().Msg("Running environment setup script...")
					if err := core.RunScript(envScriptPath); err != nil {
						scriptsOp.SetState(utils.StateWarning)
						utils.Logger.Warn().Err(err).Msg("Error running environment setup script")
						// Continue anyway
					}
				}
			}

			// Run machine-specific setup scripts
			if !skipMachine {
				machineScriptPath := filepath.Join(dotpilotDir, "machine", hostname, "install_packages.sh")
				if _, err := os.Stat(machineScriptPath); err == nil {
					utils.Logger.Info().Msg("Running machine-specific setup script...")
					if err := core.RunScript(machineScriptPath); err != nil {
						scriptsOp.SetState(utils.StateWarning)
						utils.Logger.Warn().Err(err).Msg("Error running machine-specific setup script")
						// Continue anyway
					}
				}
			}

			scriptsOp.SetState(utils.StateSuccess)
			scriptsOp.Stop()
		}

		utils.Logger.Info().Msg("Bootstrap completed successfully!")
	},
}

func init() {
	// Add flags
	bootstrapCmd.Flags().BoolVar(&skipCommon, "skip-common", false, "Skip applying common dotfiles")
	bootstrapCmd.Flags().BoolVar(&skipEnv, "skip-env", false, "Skip applying environment-specific dotfiles")
	bootstrapCmd.Flags().BoolVar(&skipMachine, "skip-machine", false, "Skip applying machine-specific dotfiles")
	bootstrapCmd.Flags().BoolVar(&skipSetupScripts, "skip-setup-scripts", false, "Skip running setup scripts")
	bootstrapCmd.Flags().BoolVar(&forceOverwrite, "force", false, "Force overwrite existing files without prompting")
}