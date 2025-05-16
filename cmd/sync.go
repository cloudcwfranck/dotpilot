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
        noPull            bool
        noPush            bool
        noBackup          bool
        noDiffPrompt      bool
        dryRun            bool
        resolveConflicts  bool
        conflictStrategy  string
        noProgress        bool // Whether to disable progress indicators
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
                
                // Initialize operation manager for progress tracking
                var operationManager *utils.OperationManager
                if !noProgress && !dryRun {
                    operationManager = utils.NewOperationManager()
                }

                // Check for uncommitted changes
                hasChanges, err := core.HasUncommittedChanges(dotpilotDir)
                if err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to check for uncommitted changes")
                        os.Exit(1)
                }

                if hasChanges {
                        utils.Logger.Info().Msg("Uncommitted changes detected, committing...")
                        
                        // Create progress for commit operation
                        var commitOp *utils.Operation
                        if operationManager != nil {
                            commitOp = operationManager.AddOperation("commit", "Committing changes...", utils.Spinner)
                            commitOp.Start()
                        }
                        
                        if err := core.CommitChanges(dotpilotDir, "Auto-commit before sync"); err != nil {
                                if commitOp != nil {
                                    commitOp.Stop()
                                }
                                utils.Logger.Error().Err(err).Msg("Failed to commit changes")
                                os.Exit(1)
                        }
                        
                        if commitOp != nil {
                            commitOp.Stop()
                        }
                }

                // Pull changes
                if !noPull {
                        utils.Logger.Info().Msg("Pulling changes from remote...")
                        
                        if dryRun {
                                utils.Logger.Info().Msg("[DRY RUN] Would pull changes from remote")
                        } else {
                                // Create progress for pull operation
                                var pullOp *utils.Operation
                                if operationManager != nil {
                                    pullOp = operationManager.AddOperation("pull", "Pulling changes from remote...", utils.Bounce)
                                    pullOp.Start()
                                    pullOp.SimulateProgress(5) // Simulate progress for 5 seconds
                                }
                                
                                if err := core.PullChanges(dotpilotDir); err != nil {
                                        if pullOp != nil {
                                            pullOp.Stop()
                                        }
                                        utils.Logger.Error().Err(err).Msg("Failed to pull changes")
                                        os.Exit(1)
                                }
                                
                                if pullOp != nil {
                                    pullOp.Stop()
                                }

                                // Run post-pull hooks
                                utils.Logger.Info().Msg("Running post-pull hooks...")
                                
                                // Create progress for hooks operation
                                var hooksOp *utils.Operation
                                if operationManager != nil {
                                    hooksOp = operationManager.AddOperation("hooks", "Running post-pull hooks...", utils.Spinner)
                                    hooksOp.Start()
                                }
                                
                                if err := core.RunHooks(dotpilotDir, environment, "postpull.sh"); err != nil {
                                        if hooksOp != nil {
                                            hooksOp.Stop()
                                        }
                                        utils.Logger.Error().Err(err).Msg("Failed to run post-pull hooks")
                                        // Continue anyway
                                }
                                
                                if hooksOp != nil {
                                    hooksOp.Stop()
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
                                
                                // Create progress for conflict resolution (only for non-interactive strategies)
                                var conflictOp *utils.Operation
                                if operationManager != nil && strategy != core.StrategyInteractive {
                                    conflictOp = operationManager.AddOperation("conflicts", 
                                        fmt.Sprintf("Resolving conflicts with %s strategy...", conflictStrategy), 
                                        utils.Dots)
                                    conflictOp.Start()
                                }
                                
                                if err := core.ResolveConflicts(dotpilotDir, strategy); err != nil {
                                        if conflictOp != nil {
                                            conflictOp.Stop()
                                        }
                                        utils.Logger.Error().Err(err).Msg("Failed to resolve conflicts")
                                        os.Exit(1)
                                }
                                
                                if conflictOp != nil {
                                    conflictOp.Stop()
                                }
                        }
                }

                // Apply configurations
                utils.Logger.Info().Msg("Applying configurations...")
                if dryRun {
                        utils.Logger.Info().Msg("[DRY RUN] Would apply configurations")
                } else {
                        // Create progress for applying configurations
                        var configOp *utils.Operation
                        if operationManager != nil {
                            configOp = operationManager.AddOperation("config", "Applying configurations...", utils.Bar)
                            configOp.Start()
                            configOp.SimulateProgress(3) // Simulate progress for 3 seconds
                        }
                        
                        backupEnabled := !noBackup
                        diffPromptEnabled := !noDiffPrompt
                        
                        // Progress indicator is not compatible with diff prompts, so disable it temporarily
                        if diffPromptEnabled && configOp != nil {
                            configOp.Stop()
                            configOp = nil
                        }
                        
                        if err := core.ApplyConfigurationsWithOptions(dotpilotDir, environment, backupEnabled, diffPromptEnabled); err != nil {
                                if configOp != nil {
                                    configOp.Stop()
                                }
                                utils.Logger.Error().Err(err).Msg("Failed to apply configurations")
                                os.Exit(1)
                        }
                        
                        if configOp != nil {
                            configOp.Stop()
                        }
                }

                // Push changes
                if !noPush {
                        utils.Logger.Info().Msg("Pushing changes to remote...")
                        if dryRun {
                                utils.Logger.Info().Msg("[DRY RUN] Would push changes to remote")
                        } else {
                                // Create progress for push operation
                                var pushOp *utils.Operation
                                if operationManager != nil {
                                    pushOp = operationManager.AddOperation("push", "Pushing changes to remote...", utils.Bounce)
                                    pushOp.Start()
                                    pushOp.SimulateProgress(4) // Simulate progress for 4 seconds
                                }
                                
                                if err := core.PushChanges(dotpilotDir); err != nil {
                                        if pushOp != nil {
                                            pushOp.Stop()
                                        }
                                        utils.Logger.Error().Err(err).Msg("Failed to push changes")
                                        os.Exit(1)
                                }
                                
                                if pushOp != nil {
                                    pushOp.Stop()
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
        syncCmd.Flags().BoolVar(&noProgress, "no-progress", false, "Disable animated progress indicators")
        
        // Advanced conflict resolution flags
        syncCmd.Flags().BoolVar(&resolveConflicts, "resolve-conflicts", false, "Detect and resolve conflicts between local and remote files")
        syncCmd.Flags().StringVar(&conflictStrategy, "strategy", "interactive", 
                "Conflict resolution strategy: interactive, keep-local, keep-remote, merge, or backup-both")
        
        // Add completion for strategy flag
        if err := syncCmd.RegisterFlagCompletionFunc("strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
}
