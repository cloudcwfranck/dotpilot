package cmd

import (
        "fmt"
        "os"

        "github.com/dotpilot/core"
        "github.com/dotpilot/utils"
        "github.com/spf13/cobra"
)

var (
        remoteRepo    string
        environment   string
        forceInit     bool
        skipPackages  bool
        skipHooks     bool
        packageSystem string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
        Use:   "init",
        Short: "Initialize dotpilot with a remote repository",
        Long: `Initialize dotpilot by cloning the specified remote repository,
setting up configurations, and optionally installing packages and running hooks.

For example:
  dotpilot init --remote https://github.com/username/dotfiles.git --env dev`,
        Run: func(cmd *cobra.Command, args []string) {
                if remoteRepo == "" {
                        utils.Logger.Error().Msg("Remote repository URL is required")
                        cmd.Help()
                        os.Exit(1)
                }

                // Get the home directory
                home, err := os.UserHomeDir()
                if err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to get home directory")
                        os.Exit(1)
                }

                // Create .dotpilot directory
                dotpilotDir := fmt.Sprintf("%s/.dotpilot", home)
                if _, err := os.Stat(dotpilotDir); !os.IsNotExist(err) && !forceInit {
                        utils.Logger.Error().Msg("Dotpilot directory already exists. Use --force to reinitialize")
                        os.Exit(1)
                }

                if forceInit && !os.IsNotExist(err) {
                        utils.Logger.Info().Msg("Removing existing dotpilot directory...")
                        if err := os.RemoveAll(dotpilotDir); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to remove existing dotpilot directory")
                                os.Exit(1)
                        }
                }

                // Initialize dotpilot
                utils.Logger.Info().Msgf("Initializing dotpilot with repository: %s", remoteRepo)
                if err := core.InitializeRepo(remoteRepo, dotpilotDir, environment); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to initialize repository")
                        os.Exit(1)
                }

                // Apply configurations
                utils.Logger.Info().Msg("Applying configurations...")
                if err := core.ApplyConfigurations(dotpilotDir, environment); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to apply configurations")
                        os.Exit(1)
                }

                // Run pre-installation hooks
                if !skipHooks {
                        utils.Logger.Info().Msg("Running pre-installation hooks...")
                        if err := core.RunHooks(dotpilotDir, environment, "preinstall.sh"); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to run pre-installation hooks")
                                os.Exit(1)
                        }
                }

                // Install packages
                if !skipPackages {
                        utils.Logger.Info().Msg("Installing packages...")
                        if err := core.InstallPackages(dotpilotDir, environment, packageSystem); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to install packages")
                                os.Exit(1)
                        }
                }

                // Run post-installation hooks
                if !skipHooks {
                        utils.Logger.Info().Msg("Running post-installation hooks...")
                        if err := core.RunHooks(dotpilotDir, environment, "postinstall.sh"); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to run post-installation hooks")
                                os.Exit(1)
                        }
                }

                utils.Logger.Info().Msg("Dotpilot initialized successfully!")
        },
}

func init() {
        initCmd.Flags().StringVar(&remoteRepo, "remote", "", "URL of the remote Git repository (required)")
        initCmd.Flags().StringVar(&environment, "env", "default", "Environment to use (e.g., dev, prod)")
        initCmd.Flags().BoolVar(&forceInit, "force", false, "Force reinitialization if dotpilot is already initialized")
        initCmd.Flags().BoolVar(&skipPackages, "skip-packages", false, "Skip package installation")
        initCmd.Flags().BoolVar(&skipHooks, "skip-hooks", false, "Skip running hooks")
        initCmd.Flags().StringVar(&packageSystem, "package-system", "", "Override automatic package system detection (apt, brew, yay)")

        initCmd.MarkFlagRequired("remote")
        
        // Add completion for environment flag
        if err := initCmd.RegisterFlagCompletionFunc("env", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                return []string{"default", "dev", "prod", "test", "hardened"}, cobra.ShellCompDirectiveNoFileComp
        }); err != nil {
                utils.Logger.Debug().Err(err).Msg("Failed to register environment flag completion")
        }
        
        // Add completion for package system flag
        if err := initCmd.RegisterFlagCompletionFunc("package-system", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                return []string{"apt", "brew", "yay", "dnf", "pacman", "zypper"}, cobra.ShellCompDirectiveNoFileComp
        }); err != nil {
                utils.Logger.Debug().Err(err).Msg("Failed to register package-system flag completion")
        }
}
