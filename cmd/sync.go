package cmd

import (
        "os"
        "path/filepath"

        "github.com/dotpilot/core"
        "github.com/dotpilot/utils"
        "github.com/spf13/cobra"
)

var (
        noPull            bool
        noPush            bool
        noBackup          bool
        noDiffPrompt      bool
        dryRun            bool
        resolveConflicts  bool
        conflictStrategy  string
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
        Use:   "sync",
        Short: "Sync dotfiles with remote repository",
        Long: `Sync dotfiles between the local dotpilot repository and the remote repository.
By default, this will pull changes from the remote, apply them to the local system,
and push any local changes back to the remote.

For example:
  dotpilot sync
  dotpilot sync --no-push
  dotpilot sync --dry-run
  dotpilot sync --resolve-conflicts --strategy=interactive`,
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

                // Sync process
                utils.Logger.Info().Msg("Starting sync process...")

                // Check for uncommitted changes
                hasChanges, err := core.HasUncommittedChanges(dotpilotDir)
                if err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to check for uncommitted changes")
                        os.Exit(1)
                }

                if hasChanges {
                        utils.Logger.Info().Msg("Uncommitted changes detected, committing...")
                        if err := core.CommitChanges(dotpilotDir, "Auto-commit before sync"); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to commit changes")
                                os.Exit(1)
                        }
                }

                // Pull changes
                if !noPull {
                        utils.Logger.Info().Msg("Pulling changes from remote...")
                        if dryRun {
                                utils.Logger.Info().Msg("[DRY RUN] Would pull changes from remote")
                        } else {
                                if err := core.PullChanges(dotpilotDir); err != nil {
                                        utils.Logger.Error().Err(err).Msg("Failed to pull changes")
                                        os.Exit(1)
                                }

                                // Run post-pull hooks
                                utils.Logger.Info().Msg("Running post-pull hooks...")
                                if err := core.RunHooks(dotpilotDir, environment, "postpull.sh"); err != nil {
                                        utils.Logger.Error().Err(err).Msg("Failed to run post-pull hooks")
                                        // Continue anyway
                                }
                        }
                }

                // Resolve conflicts if requested
                if resolveConflicts {
                        utils.Logger.Info().Msgf("Resolving conflicts with strategy: %s", conflictStrategy)
                        
                        if dryRun {
                                utils.Logger.Info().Msg("[DRY RUN] Would resolve conflicts")
                        } else {
                                // Parse the strategy
                                var strategy core.ConflictResolutionStrategy
                                switch conflictStrategy {
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
                                        utils.Logger.Warn().Msgf("Unknown conflict strategy: %s, using interactive", conflictStrategy)
                                        strategy = core.StrategyInteractive
                                }
                                
                                if err := core.ResolveConflicts(dotpilotDir, strategy); err != nil {
                                        utils.Logger.Error().Err(err).Msg("Failed to resolve conflicts")
                                        os.Exit(1)
                                }
                        }
                }

                // Apply configurations
                utils.Logger.Info().Msg("Applying configurations...")
                if dryRun {
                        utils.Logger.Info().Msg("[DRY RUN] Would apply configurations")
                } else {
                        backupEnabled := !noBackup
                        diffPromptEnabled := !noDiffPrompt
                        if err := core.ApplyConfigurationsWithOptions(dotpilotDir, environment, backupEnabled, diffPromptEnabled); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to apply configurations")
                                os.Exit(1)
                        }
                }

                // Push changes
                if !noPush {
                        utils.Logger.Info().Msg("Pushing changes to remote...")
                        if dryRun {
                                utils.Logger.Info().Msg("[DRY RUN] Would push changes to remote")
                        } else {
                                if err := core.PushChanges(dotpilotDir); err != nil {
                                        utils.Logger.Error().Err(err).Msg("Failed to push changes")
                                        os.Exit(1)
                                }
                        }
                }

                utils.Logger.Info().Msg("Sync completed successfully!")
        },
}

func init() {
        syncCmd.Flags().BoolVar(&noPull, "no-pull", false, "Skip pulling changes from remote")
        syncCmd.Flags().BoolVar(&noPush, "no-push", false, "Skip pushing changes to remote")
        syncCmd.Flags().BoolVar(&noBackup, "no-backup", false, "Skip backing up files before overwriting")
        syncCmd.Flags().BoolVar(&noDiffPrompt, "no-diff-prompt", false, "Skip prompting for diffs before applying changes")
        syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
        
        // Advanced conflict resolution flags
        syncCmd.Flags().BoolVar(&resolveConflicts, "resolve-conflicts", false, "Detect and resolve conflicts between local and remote files")
        syncCmd.Flags().StringVar(&conflictStrategy, "strategy", "interactive", 
                "Conflict resolution strategy: interactive, keep-local, keep-remote, merge, or backup-both")
}
