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
	secretDestination string
	secretOverwrite   bool
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage encrypted secrets",
	Long: `Manage encrypted secrets in your dotfiles repository.
Allows you to securely store sensitive configuration files
that will be encrypted before being stored in the Git repository.

DotPilot will use GPG if available, or fall back to AES-256 encryption.`,
}

// addSecretCmd represents the add-secret command
var addSecretCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add an encrypted secret",
	Long: `Add a file as an encrypted secret to the dotpilot repository.
The file will be encrypted before being stored in the repository.

For example:
  dotpilot secrets add ~/.aws/credentials
  dotpilot secrets add ~/.ssh/id_rsa --name ssh_key`,
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
		var secretName string
		if secretDestination != "" {
			secretName = secretDestination
		} else {
			// Use filename as secret name (with directory structure removed)
			secretName = filepath.Base(absPath)
		}

		// Create secret manager
		secretManager := core.NewSecretManager(dotpilotDir)
		if err := secretManager.Initialize(); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to initialize secret manager")
			os.Exit(1)
		}

		// Check if secret already exists
		secrets, err := secretManager.ListSecrets()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to list secrets")
			os.Exit(1)
		}

		secretExists := false
		for _, s := range secrets {
			if s == secretName {
				secretExists = true
				break
			}
		}

		if secretExists && !secretOverwrite {
			utils.Logger.Error().Msgf("Secret %s already exists. Use --overwrite to replace it.", secretName)
			os.Exit(1)
		}

		// Encrypt the file
		utils.Logger.Info().Msgf("Encrypting %s as %s", absPath, secretName)
		if err := secretManager.EncryptFile(absPath, secretName); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to encrypt file")
			os.Exit(1)
		}

		utils.Logger.Info().Msgf("Successfully encrypted %s", secretName)

		// Commit changes
		utils.Logger.Info().Msg("Committing changes...")
		if err := core.CommitChanges(dotpilotDir, fmt.Sprintf("Added encrypted secret: %s", secretName)); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to commit changes")
			os.Exit(1)
		}

		utils.Logger.Info().Msg("Secret added successfully!")
	},
}

// getSecretCmd represents the get-secret command
var getSecretCmd = &cobra.Command{
	Use:   "get [name] [destination]",
	Short: "Get a decrypted secret",
	Long: `Decrypt and retrieve a secret from the dotpilot repository.
The secret will be decrypted and saved to the specified destination.

For example:
  dotpilot secrets get aws_credentials ~/.aws/credentials
  dotpilot secrets get ssh_key ~/.ssh/id_rsa`,
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
		if _, err := os.Stat(destPath); err == nil && !secretOverwrite {
			utils.Logger.Error().Msgf("Destination file already exists: %s. Use --overwrite to replace it.", destPath)
			os.Exit(1)
		}

		// Create secret manager
		secretManager := core.NewSecretManager(dotpilotDir)
		if err := secretManager.Initialize(); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to initialize secret manager")
			os.Exit(1)
		}

		// Decrypt the secret
		utils.Logger.Info().Msgf("Decrypting %s to %s", secretName, destPath)
		if err := secretManager.DecryptFile(secretName, destPath); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to decrypt secret")
			os.Exit(1)
		}

		utils.Logger.Info().Msgf("Successfully decrypted %s to %s", secretName, destPath)
	},
}

// listSecretsCmd represents the list-secrets command
var listSecretsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Long: `List all encrypted secrets stored in the dotpilot repository.

For example:
  dotpilot secrets list`,
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

		// Create secret manager
		secretManager := core.NewSecretManager(dotpilotDir)
		if err := secretManager.Initialize(); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to initialize secret manager")
			os.Exit(1)
		}

		// List secrets
		secrets, err := secretManager.ListSecrets()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to list secrets")
			os.Exit(1)
		}

		if len(secrets) == 0 {
			fmt.Println("No secrets found.")
			return
		}

		fmt.Println("Encrypted secrets:")
		for _, s := range secrets {
			fmt.Printf("- %s\n", s)
		}
	},
}

// removeSecretCmd represents the remove-secret command
var removeSecretCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a secret",
	Long: `Remove an encrypted secret from the dotpilot repository.

For example:
  dotpilot secrets remove aws_credentials`,
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

		// Create secret manager
		secretManager := core.NewSecretManager(dotpilotDir)
		if err := secretManager.Initialize(); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to initialize secret manager")
			os.Exit(1)
		}

		// Remove the secret
		utils.Logger.Info().Msgf("Removing secret %s", secretName)
		if err := secretManager.RemoveSecret(secretName); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to remove secret")
			os.Exit(1)
		}

		// Commit changes
		utils.Logger.Info().Msg("Committing changes...")
		if err := core.CommitChanges(dotpilotDir, fmt.Sprintf("Removed encrypted secret: %s", secretName)); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to commit changes")
			os.Exit(1)
		}

		utils.Logger.Info().Msgf("Successfully removed secret %s", secretName)
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(addSecretCmd)
	secretsCmd.AddCommand(getSecretCmd)
	secretsCmd.AddCommand(listSecretsCmd)
	secretsCmd.AddCommand(removeSecretCmd)

	// Add flags for add-secret command
	addSecretCmd.Flags().StringVar(&secretDestination, "name", "", "Custom name for the secret")
	addSecretCmd.Flags().BoolVar(&secretOverwrite, "overwrite", false, "Overwrite existing secret")

	// Add flags for get-secret command
	getSecretCmd.Flags().BoolVar(&secretOverwrite, "overwrite", false, "Overwrite existing file")
}