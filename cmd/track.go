package cmd

import (
        "os"
        "path/filepath"

        "github.com/dotpilot/core"
        "github.com/dotpilot/utils"
        "github.com/spf13/cobra"
)

var (
        destPath      string
        overwrite     bool
        environmentOp string
)

// trackCmd represents the track command
var trackCmd = &cobra.Command{
        Use:   "track [file or directory]",
        Short: "Track a file or directory in dotpilot",
        Long: `Track a file or directory to be managed by dotpilot.
This will copy the file or directory to the dotpilot repository and create a symlink
in the original location.

For example:
  dotpilot track ~/.zshrc
  dotpilot track ~/.config/nvim --env dev`,
        Args: cobra.MinimumNArgs(1),
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

                // Track each file or directory
                for _, src := range args {
                        // Expand ~ to home directory
                        if src[0] == '~' {
                                src = filepath.Join(home, src[1:])
                        }

                        // Get absolute path
                        absPath, err := filepath.Abs(src)
                        if err != nil {
                                utils.Logger.Error().Err(err).Msgf("Failed to get absolute path for %s", src)
                                continue
                        }

                        // Check if file or directory exists
                        if _, err := os.Stat(absPath); os.IsNotExist(err) {
                                utils.Logger.Error().Msgf("File or directory does not exist: %s", absPath)
                                continue
                        }

                        // Determine destination path within dotpilot
                        var destination string
                        if destPath != "" {
                                destination = destPath
                        } else {
                                // Make path relative to home if it's under home
                                relPath := absPath
                                if filepath.HasPrefix(absPath, home) {
                                        relPath, _ = filepath.Rel(home, absPath)
                                }

                                // Determine environment path
                                var envDir string
                                switch environmentOp {
                                case "common":
                                        envDir = "common"
                                case "machine":
                                        hostname, err := os.Hostname()
                                        if err != nil {
                                                utils.Logger.Error().Err(err).Msg("Failed to get hostname")
                                                hostname = "unknown"
                                        }
                                        envDir = filepath.Join("machine", hostname)
                                default:
                                        if environmentOp != "" {
                                                envDir = filepath.Join("envs", environmentOp)
                                        } else {
                                                // Use current environment from config
                                                cfg := core.GetConfig()
                                                if cfg.CurrentEnvironment != "" {
                                                        envDir = filepath.Join("envs", cfg.CurrentEnvironment)
                                                } else {
                                                        envDir = "common"
                                                }
                                        }
                                }

                                destination = filepath.Join(dotpilotDir, envDir, relPath)
                        }

                        // Track the file
                        if err := core.TrackFile(absPath, destination, dotpilotDir, overwrite); err != nil {
                                utils.Logger.Error().Err(err).Msgf("Failed to track %s", absPath)
                                continue
                        }

                        utils.Logger.Info().Msgf("Successfully tracked %s", absPath)
                }

                // Commit changes
                utils.Logger.Info().Msg("Committing changes...")
                if err := core.CommitChanges(dotpilotDir, "Added tracked files via dotpilot"); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to commit changes")
                        os.Exit(1)
                }

                utils.Logger.Info().Msg("Files tracked successfully!")
        },
}

func init() {
        trackCmd.Flags().StringVar(&destPath, "dest", "", "Custom destination path in the dotpilot repo")
        trackCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")
        trackCmd.Flags().StringVar(&environmentOp, "env", "", "Environment to track in (common, machine, or specific environment name)")

        // Add file path completion for track command arguments
        if err := trackCmd.RegisterFlagCompletionFunc("env", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                // Get available environments
                envs := []string{"common", "machine"}
                
                // Add environment-specific directories
                home, err := os.UserHomeDir()
                if err == nil {
                        dotpilotDir := filepath.Join(home, ".dotpilot")
                        envsDir := filepath.Join(dotpilotDir, "envs")
                        if info, err := os.Stat(envsDir); err == nil && info.IsDir() {
                                if dirs, err := os.ReadDir(envsDir); err == nil {
                                        for _, dir := range dirs {
                                                if dir.IsDir() && !strings.HasPrefix(dir.Name(), ".") {
                                                        envs = append(envs, dir.Name())
                                                }
                                        }
                                }
                        }
                }
                
                return envs, cobra.ShellCompDirectiveNoFileComp
        }); err != nil {
                utils.Logger.Debug().Err(err).Msg("Failed to register environment flag completion")
        }

        // Enable filepath completion for arguments
        trackCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                return nil, cobra.ShellCompDirectiveDefault
        }
}
