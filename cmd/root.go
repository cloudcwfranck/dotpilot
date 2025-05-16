package cmd

import (
        "fmt"
        "os"
        "path/filepath"

        "github.com/dotpilot/core"
        "github.com/dotpilot/utils"
        "github.com/spf13/cobra"
)

var (
        cfgFile string
        verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
        Use:   "dotpilot",
        Short: "Manage and sync dotfiles across multiple machines",
        Long: `DotPilot is a cross-platform tool to manage and sync dotfiles across 
multiple machines with environment-specific overrides.

It uses a Git-backed system to track changes to dotfiles, supports scoped
environments (e.g., dev, prod, hardened), and includes machine-specific
configurations.`,
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
                // Set up logging level
                if verbose {
                        utils.SetLogLevel("debug")
                }
        },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
        return rootCmd.Execute()
}

func init() {
        cobra.OnInitialize(initConfig)

        // Global flags
        rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dotpilotrc)")
        rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

        // Setup bash completion
        rootCmd.CompletionOptions.DisableDefaultCmd = false
        rootCmd.CompletionOptions.DisableNoDescFlag = false

        // Add subcommands
        rootCmd.AddCommand(initCmd)
        rootCmd.AddCommand(trackCmd)
        rootCmd.AddCommand(syncCmd)
        rootCmd.AddCommand(bootstrapCmd)
        rootCmd.AddCommand(statusCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
        if cfgFile != "" {
                // Use config file from the flag
                core.LoadConfig(cfgFile)
        } else {
                // Find home directory
                home, err := os.UserHomeDir()
                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }

                // Search for config in home directory
                defaultConfigPath := filepath.Join(home, ".dotpilotrc")
                if _, err := os.Stat(defaultConfigPath); err == nil {
                        core.LoadConfig(defaultConfigPath)
                } else {
                        utils.Logger.Debug().Msg("No config file found, using defaults")
                        core.InitDefaultConfig()
                }
        }
}
