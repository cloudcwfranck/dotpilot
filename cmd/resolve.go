package cmd

import (
        "os"
        "path/filepath"

        "github.com/dotpilot/core"
        "github.com/dotpilot/utils"
        "github.com/spf13/cobra"
)

var (
        resolveStrategy string
)

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
        Use:   "resolve",
        Short: "Resolve conflicts between local and tracked dotfiles",
        Long: `Detect and resolve conflicts between local dotfiles and their
tracked versions in the dotpilot repository.

Strategies available:
- interactive: Prompts for each conflict (default)
- keep-local: Keep the local version of conflicting files
- keep-remote: Keep the remote version of conflicting files
- merge: Attempt to merge changes using a merge tool
- backup-both: Keep both versions with backups

For example:
  dotpilot resolve
  dotpilot resolve --strategy=keep-remote
  dotpilot resolve --strategy=merge`,
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

                // Parse the strategy
                var strategy core.ConflictResolutionStrategy
                switch resolveStrategy {
                case "interactive":
                        strategy = core.StrategyInteractive
                case "keep-local":
                        strategy = core.StrategyKeepLocal
                case "keep-remote":
                        strategy = core.StrategyKeepRemote
                case "merge":
                        strategy = core.StrategyMerge
                case "backup-both":
                        strategy = core.StrategyBackupBoth
                default:
                        utils.Logger.Warn().Msgf("Unknown conflict strategy: %s, using interactive", resolveStrategy)
                        strategy = core.StrategyInteractive
                }

                utils.Logger.Info().Msgf("Checking for conflicts with strategy: %s", strategy)
                if err := core.ResolveConflicts(dotpilotDir, strategy); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to resolve conflicts")
                        os.Exit(1)
                }

                utils.Logger.Info().Msg("Conflict resolution completed successfully")
        },
}

func init() {
        resolveCmd.Flags().StringVar(&resolveStrategy, "strategy", "interactive",
                "Conflict resolution strategy: interactive, keep-local, keep-remote, merge, or backup-both")

        // Add completion for strategy flag
        if err := resolveCmd.RegisterFlagCompletionFunc("strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                strategies := []string{
                        "interactive",   // Prompt for each conflict
                        "keep-local",    // Keep local versions
                        "keep-remote",   // Keep remote versions  
                        "merge",         // Try to merge changes
                        "backup-both",   // Keep both versions
                }
                return strategies, cobra.ShellCompDirectiveNoFileComp
        }); err != nil {
                utils.Logger.Debug().Err(err).Msg("Failed to register strategy flag completion")
        }

        rootCmd.AddCommand(resolveCmd)
}