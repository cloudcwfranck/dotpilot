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
        sopsSecretName     string
        sopsSecretOverwrite bool
        sopsSecretEdit     bool
        sopsNoProgress    bool // Whether to disable progress indicators
)

// sopsCmd represents the sops command
var sopsCmd = &cobra.Command{
        Use:   "sops",
        Short: "Manage encrypted secrets with SOPS and GPG",
        Long: `Manage encrypted secrets using Mozilla SOPS and GPG for enhanced security.
SOPS provides advanced encryption for configuration files and secrets,
allowing for secure storage of sensitive data in Git.

Requirements:
- GPG must be installed with a key generated
- SOPS must be installed (https://github.com/mozilla/sops)

DotPilot will create a SOPS configuration file and use your GPG key
for encryption and decryption.`,
}

// sopsAddCmd represents the sops add command
var sopsAddCmd = &cobra.Command{
        Use:   "add [file]",
        Short: "Add an encrypted secret using SOPS",
        Long: `Add a file as an encrypted secret to the dotpilot repository using SOPS.
The file will be encrypted with GPG before being stored in the repository.

For example:
  dotpilot sops add ~/.aws/credentials
  dotpilot sops add ~/.ssh/id_rsa --name ssh_key
  dotpilot sops add ~/.npmrc --edit`,
        Args: cobra.ExactArgs(1),
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

                // Expand ~ to home directory
                srcPath := args[0]
                if srcPath[0] == '~' {
                        srcPath = filepath.Join(home, srcPath[1:])
                }

                // Get absolute path
                absPath, err := filepath.Abs(srcPath)
                if err != nil {
                        utils.Logger.Error().Err(err).Msgf("Failed to get absolute path for %s", srcPath)
                        os.Exit(1)
                }

                // Check if file exists
                if _, err := os.Stat(absPath); os.IsNotExist(err) {
                        utils.Logger.Error().Msgf("File does not exist: %s", absPath)
                        os.Exit(1)
                }

                // Determine secret name
                if sopsSecretName == "" {
                        // Use filename as secret name (with directory structure removed)
                        sopsSecretName = filepath.Base(absPath)
                }

                // Create SOPS manager
                sopsManager := core.NewSopsManager(dotpilotDir)
                if err := sopsManager.Initialize(); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to initialize SOPS manager")
                        os.Exit(1)
                }

                // Check if secret already exists
                secrets, err := sopsManager.ListSecrets()
                if err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to list secrets")
                        os.Exit(1)
                }

                secretExists := false
                for _, s := range secrets {
                        if s == sopsSecretName {
                                secretExists = true
                                break
                        }
                }

                if secretExists && !sopsSecretOverwrite {
                        utils.Logger.Error().Msgf("Secret %s already exists. Use --overwrite to replace it.", sopsSecretName)
                        os.Exit(1)
                }

                // Encrypt the file
                utils.Logger.Info().Msgf("Encrypting %s as %s", absPath, sopsSecretName)
                
                // Create progress for encryption operation
                var encryptOp *utils.Operation
                if !sopsNoProgress {
                    encryptOp = utils.NewOperation("encrypt", fmt.Sprintf("Encrypting %s...", filepath.Base(absPath)), utils.Spinner)
                    encryptOp.Start()
                    // For larger files, encryption might take some time
                    encryptOp.SimulateProgress(3) // Simulate progress for 3 seconds
                }
                
                if err := sopsManager.EncryptFile(absPath, sopsSecretName); err != nil {
                        if encryptOp != nil {
                            encryptOp.Stop()
                        }
                        utils.Logger.Error().Err(err).Msg("Failed to encrypt file")
                        os.Exit(1)
                }
                
                if encryptOp != nil {
                    encryptOp.Stop()
                }

                utils.Logger.Info().Msgf("Successfully encrypted %s", sopsSecretName)

                // If edit flag is set, open the secret for editing
                if sopsSecretEdit {
                        utils.Logger.Info().Msg("Opening secret for editing...")
                        if err := sopsManager.EditSecret(sopsSecretName); err != nil {
                                utils.Logger.Error().Err(err).Msg("Failed to edit secret")
                                os.Exit(1)
                        }
                }

                // Commit changes
                utils.Logger.Info().Msg("Committing changes...")
                if err := core.CommitChanges(dotpilotDir, fmt.Sprintf("Added encrypted SOPS secret: %s", sopsSecretName)); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to commit changes")
                        os.Exit(1)
                }

                utils.Logger.Info().Msg("Secret added successfully!")
        },
}

// sopsGetCmd represents the sops get command
var sopsGetCmd = &cobra.Command{
        Use:   "get [name] [destination]",
        Short: "Get a decrypted secret",
        Long: `Decrypt and retrieve a secret from the dotpilot repository.
The secret will be decrypted and saved to the specified destination.

For example:
  dotpilot sops get aws_credentials ~/.aws/credentials
  dotpilot sops get ssh_key ~/.ssh/id_rsa`,
        Args: cobra.ExactArgs(2),
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

                // Get secret name and destination
                secretName := args[0]
                destPath := args[1]

                // Expand ~ to home directory in destination
                if destPath[0] == '~' {
                        destPath = filepath.Join(home, destPath[1:])
                }

                // Get absolute path for destination
                destPath, err = filepath.Abs(destPath)
                if err != nil {
                        utils.Logger.Error().Err(err).Msgf("Failed to get absolute path for %s", destPath)
                        os.Exit(1)
                }

                // Create parent directories if needed
                parentDir := filepath.Dir(destPath)
                if err := os.MkdirAll(parentDir, 0755); err != nil {
                        utils.Logger.Error().Err(err).Msgf("Failed to create directory %s", parentDir)
                        os.Exit(1)
                }

                // Check if destination file exists
                if _, err := os.Stat(destPath); err == nil && !sopsSecretOverwrite {
                        utils.Logger.Error().Msgf("Destination file already exists: %s. Use --overwrite to replace it.", destPath)
                        os.Exit(1)
                }

                // Create SOPS manager
                sopsManager := core.NewSopsManager(dotpilotDir)
                if err := sopsManager.Initialize(); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to initialize SOPS manager")
                        os.Exit(1)
                }

                // Decrypt the secret
                utils.Logger.Info().Msgf("Decrypting %s to %s", secretName, destPath)
                
                // Create progress for decryption operation
                var decryptOp *utils.Operation
                if !sopsNoProgress {
                    decryptOp = utils.NewOperation("decrypt", fmt.Sprintf("Decrypting %s...", secretName), utils.Dots)
                    decryptOp.Start()
                    decryptOp.SimulateProgress(2) // Simulate progress for 2 seconds
                }
                
                if err := sopsManager.DecryptFile(secretName, destPath); err != nil {
                        if decryptOp != nil {
                            decryptOp.Stop()
                        }
                        utils.Logger.Error().Err(err).Msg("Failed to decrypt secret")
                        os.Exit(1)
                }
                
                if decryptOp != nil {
                    decryptOp.Stop()
                }

                utils.Logger.Info().Msgf("Successfully decrypted %s to %s", secretName, destPath)
        },
}

// sopsListCmd represents the sops list command
var sopsListCmd = &cobra.Command{
        Use:   "list",
        Short: "List all SOPS secrets",
        Long: `List all encrypted secrets stored in the dotpilot repository.

For example:
  dotpilot sops list`,
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

                // Create SOPS manager
                sopsManager := core.NewSopsManager(dotpilotDir)
                if err := sopsManager.Initialize(); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to initialize SOPS manager")
                        os.Exit(1)
                }

                // List secrets
                secrets, err := sopsManager.ListSecrets()
                if err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to list secrets")
                        os.Exit(1)
                }

                if len(secrets) == 0 {
                        fmt.Println("No SOPS secrets found.")
                        return
                }

                fmt.Println("SOPS encrypted secrets:")
                for _, s := range secrets {
                        fmt.Printf("- %s\n", s)
                }
        },
}

// sopsRemoveCmd represents the sops remove command
var sopsRemoveCmd = &cobra.Command{
        Use:   "remove [name]",
        Short: "Remove a SOPS secret",
        Long: `Remove an encrypted SOPS secret from the dotpilot repository.

For example:
  dotpilot sops remove aws_credentials`,
        Args: cobra.ExactArgs(1),
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

                // Get secret name
                secretName := args[0]

                // Create SOPS manager
                sopsManager := core.NewSopsManager(dotpilotDir)
                if err := sopsManager.Initialize(); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to initialize SOPS manager")
                        os.Exit(1)
                }

                // Remove the secret
                utils.Logger.Info().Msgf("Removing secret %s", secretName)
                if err := sopsManager.RemoveSecret(secretName); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to remove secret")
                        os.Exit(1)
                }

                // Commit changes
                utils.Logger.Info().Msg("Committing changes...")
                if err := core.CommitChanges(dotpilotDir, fmt.Sprintf("Removed encrypted SOPS secret: %s", secretName)); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to commit changes")
                        os.Exit(1)
                }

                utils.Logger.Info().Msgf("Successfully removed secret %s", secretName)
        },
}

// sopsEditCmd represents the sops edit command
var sopsEditCmd = &cobra.Command{
        Use:   "edit [name]",
        Short: "Edit a SOPS secret",
        Long: `Edit an encrypted SOPS secret directly.
This will open the secret in an editor for secure editing.

For example:
  dotpilot sops edit aws_credentials`,
        Args: cobra.ExactArgs(1),
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

                // Get secret name
                secretName := args[0]

                // Create SOPS manager
                sopsManager := core.NewSopsManager(dotpilotDir)
                if err := sopsManager.Initialize(); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to initialize SOPS manager")
                        os.Exit(1)
                }

                // Edit the secret
                utils.Logger.Info().Msgf("Editing secret %s", secretName)
                if err := sopsManager.EditSecret(secretName); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to edit secret")
                        os.Exit(1)
                }

                // Commit changes
                utils.Logger.Info().Msg("Committing changes...")
                if err := core.CommitChanges(dotpilotDir, fmt.Sprintf("Edited encrypted SOPS secret: %s", secretName)); err != nil {
                        utils.Logger.Error().Err(err).Msg("Failed to commit changes")
                        os.Exit(1)
                }

                utils.Logger.Info().Msgf("Successfully edited secret %s", secretName)
        },
}

func init() {
        rootCmd.AddCommand(sopsCmd)
        sopsCmd.AddCommand(sopsAddCmd)
        sopsCmd.AddCommand(sopsGetCmd)
        sopsCmd.AddCommand(sopsListCmd)
        sopsCmd.AddCommand(sopsRemoveCmd)
        sopsCmd.AddCommand(sopsEditCmd)

        // Add flags for add command
        sopsAddCmd.Flags().StringVar(&sopsSecretName, "name", "", "Custom name for the secret")
        sopsAddCmd.Flags().BoolVar(&sopsSecretOverwrite, "overwrite", false, "Overwrite existing secret")
        sopsAddCmd.Flags().BoolVar(&sopsSecretEdit, "edit", false, "Open the secret for editing after adding")
        sopsAddCmd.Flags().BoolVar(&sopsNoProgress, "no-progress", false, "Disable animated progress indicators")

        // Add flags for get command
        sopsGetCmd.Flags().BoolVar(&sopsSecretOverwrite, "overwrite", false, "Overwrite existing file")
        sopsGetCmd.Flags().BoolVar(&sopsNoProgress, "no-progress", false, "Disable animated progress indicators")

        // Add completion for file paths and secret names
        sopsAddCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                return nil, cobra.ShellCompDirectiveDefault
        }

        // Add completion for get and remove commands (complete with available secrets)
        sopsSecretCompleter := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
                // If we already have an argument, return file completion for the destination
                if len(args) > 0 && cmd == sopsGetCmd {
                        return nil, cobra.ShellCompDirectiveDefault
                }

                // Get available secrets
                home, err := os.UserHomeDir()
                if err != nil {
                        return nil, cobra.ShellCompDirectiveNoFileComp
                }

                dotpilotDir := filepath.Join(home, ".dotpilot")
                if _, err := os.Stat(dotpilotDir); os.IsNotExist(err) {
                        return nil, cobra.ShellCompDirectiveNoFileComp
                }

                secretManager := core.NewSopsManager(dotpilotDir)
                if err := secretManager.Initialize(); err != nil {
                        return nil, cobra.ShellCompDirectiveNoFileComp
                }

                secrets, err := secretManager.ListSecrets()
                if err != nil {
                        return nil, cobra.ShellCompDirectiveNoFileComp
                }

                return secrets, cobra.ShellCompDirectiveNoFileComp
        }

        sopsGetCmd.ValidArgsFunction = sopsSecretCompleter
        sopsRemoveCmd.ValidArgsFunction = sopsSecretCompleter
        sopsEditCmd.ValidArgsFunction = sopsSecretCompleter
}