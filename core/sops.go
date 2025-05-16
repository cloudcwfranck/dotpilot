package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dotpilot/utils"
)

// SopsManager handles encrypted secrets using Mozilla SOPS
type SopsManager struct {
	dotpilotDir string
	secretsDir  string
	hasSops     bool
	hasGPG      bool
	fingerprint string
}

// NewSopsManager creates a new SOPS secret manager
func NewSopsManager(dotpilotDir string) *SopsManager {
	sm := &SopsManager{
		dotpilotDir: dotpilotDir,
		secretsDir:  filepath.Join(dotpilotDir, "sops-secrets"),
	}

	// Check if SOPS is available
	_, err := exec.LookPath("sops")
	sm.hasSops = err == nil

	// Check if GPG is available
	_, err = exec.LookPath("gpg")
	sm.hasGPG = err == nil

	return sm
}

// Initialize sets up the SOPS secrets directory and configuration
func (sm *SopsManager) Initialize() error {
	// Create secrets directory if it doesn't exist
	if err := os.MkdirAll(sm.secretsDir, 0700); err != nil {
		return err
	}

	// Check if we have the required tools
	if !sm.hasSops {
		return fmt.Errorf("sops is not installed, please install it to use secure secrets encryption")
	}

	if !sm.hasGPG {
		return fmt.Errorf("gpg is not installed, please install it to use secure secrets encryption")
	}

	// Get or create GPG key for encryption
	fingerprint, err := sm.getGPGFingerprint()
	if err != nil {
		return err
	}
	sm.fingerprint = fingerprint

	// Create or update SOPS configuration file
	err = sm.createSopsConfig()
	if err != nil {
		return err
	}

	utils.Logger.Info().Msg("SOPS Secret manager initialized")
	return nil
}

// getGPGFingerprint gets or generates a GPG key for SOPS encryption
func (sm *SopsManager) getGPGFingerprint() (string, error) {
	// Try to find an existing GPG key
	cmd := exec.Command("gpg", "--list-secret-keys", "--with-colons")
	output, err := cmd.Output()
	if err == nil {
		fingerprint := sm.parseGPGFingerprint(string(output))
		if fingerprint != "" {
			utils.Logger.Debug().Msgf("Using existing GPG key: %s", fingerprint)
			return fingerprint, nil
		}
	}

	// If no key found or error, ask user to create one
	utils.Logger.Info().Msg("No suitable GPG key found. You need to create a GPG key for encrypting secrets.")
	utils.Logger.Info().Msg("Run the following command to create a key: gpg --full-generate-key")
	return "", fmt.Errorf("no GPG key available, please create one and try again")
}

// parseGPGFingerprint extracts a fingerprint from GPG output
func (sm *SopsManager) parseGPGFingerprint(output string) string {
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "sec:") && i+1 < len(lines) {
			// The next line usually contains the fingerprint
			fprLine := lines[i+1]
			parts := strings.Split(fprLine, ":")
			if len(parts) > 9 {
				return parts[9]
			}
		}
	}
	return ""
}

// createSopsConfig creates or updates the SOPS configuration file
func (sm *SopsManager) createSopsConfig() error {
	configPath := filepath.Join(sm.dotpilotDir, ".sops.yaml")
	
	// Create SOPS config content
	config := fmt.Sprintf(`---
creation_rules:
  - path_regex: %s/.*
    pgp: %s
`, sm.secretsDir, sm.fingerprint)

	// Write the configuration file
	err := ioutil.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		return err
	}

	utils.Logger.Debug().Msgf("Created SOPS config at %s", configPath)
	return nil
}

// EncryptFile encrypts a file using SOPS and stores it in the secrets directory
func (sm *SopsManager) EncryptFile(srcPath, name string) error {
	// Create destination path
	destPath := filepath.Join(sm.secretsDir, name)

	// Use SOPS to encrypt the file
	cmd := exec.Command("sops", "--encrypt", "--input-type", "json", "--output-type", "json", srcPath)
	encryptedData, err := cmd.Output()
	if err != nil {
		errOutput := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			errOutput = string(exitErr.Stderr)
		}
		return fmt.Errorf("encryption failed: %v - %s", err, errOutput)
	}

	// Write encrypted data to file
	err = ioutil.WriteFile(destPath, encryptedData, 0600)
	if err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Encrypted file with SOPS to %s", destPath)
	return nil
}

// EncryptData encrypts data directly using SOPS
func (sm *SopsManager) EncryptData(data []byte, name string) error {
	// Create destination path
	destPath := filepath.Join(sm.secretsDir, name)

	// Create a temporary file for SOPS
	tmpFile, err := ioutil.TempFile("", "sops-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Wrap data in JSON if it's not already JSON
	var jsonData []byte
	if json.Valid(data) {
		jsonData = data
	} else {
		// Create a simple JSON object with the data
		wrapper := map[string]string{
			"data": string(data),
		}
		jsonData, err = json.Marshal(wrapper)
		if err != nil {
			return err
		}
	}

	// Write to temp file
	if _, err := tmpFile.Write(jsonData); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	// Use SOPS to encrypt the file
	cmd := exec.Command("sops", "--encrypt", "--input-type", "json", "--output-type", "json", tmpFile.Name())
	encryptedData, err := cmd.Output()
	if err != nil {
		errOutput := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			errOutput = string(exitErr.Stderr)
		}
		return fmt.Errorf("encryption failed: %v - %s", err, errOutput)
	}

	// Write encrypted data to file
	err = ioutil.WriteFile(destPath, encryptedData, 0600)
	if err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Encrypted data with SOPS to %s", destPath)
	return nil
}

// DecryptFile decrypts a file from the secrets directory
func (sm *SopsManager) DecryptFile(name, destPath string) error {
	// Get the source path
	srcPath := filepath.Join(sm.secretsDir, name)

	// Check if the file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("secret file %s does not exist", name)
	}

	// Use SOPS to decrypt the file
	cmd := exec.Command("sops", "--decrypt", srcPath)
	decryptedData, err := cmd.Output()
	if err != nil {
		errOutput := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			errOutput = string(exitErr.Stderr)
		}
		return fmt.Errorf("decryption failed: %v - %s", err, errOutput)
	}

	// Check if the data is wrapped
	var jsonData map[string]interface{}
	if err := json.Unmarshal(decryptedData, &jsonData); err == nil {
		// If there's only a single "data" field, extract it
		if len(jsonData) == 1 {
			if dataStr, ok := jsonData["data"].(string); ok {
				decryptedData = []byte(dataStr)
			}
		}
	}

	// Write to destination file
	err = ioutil.WriteFile(destPath, decryptedData, 0600)
	if err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Decrypted file with SOPS to %s", destPath)
	return nil
}

// DecryptData decrypts data directly from a file
func (sm *SopsManager) DecryptData(name string) ([]byte, error) {
	// Get the source path
	srcPath := filepath.Join(sm.secretsDir, name)

	// Check if the file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("secret file %s does not exist", name)
	}

	// Use SOPS to decrypt the file
	cmd := exec.Command("sops", "--decrypt", srcPath)
	decryptedData, err := cmd.Output()
	if err != nil {
		errOutput := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			errOutput = string(exitErr.Stderr)
		}
		return nil, fmt.Errorf("decryption failed: %v - %s", err, errOutput)
	}

	// Check if the data is wrapped
	var jsonData map[string]interface{}
	if err := json.Unmarshal(decryptedData, &jsonData); err == nil {
		// If there's only a single "data" field, extract it
		if len(jsonData) == 1 {
			if dataStr, ok := jsonData["data"].(string); ok {
				decryptedData = []byte(dataStr)
			}
		}
	}

	return decryptedData, nil
}

// ListSecrets returns a list of all secret files
func (sm *SopsManager) ListSecrets() ([]string, error) {
	var secrets []string

	// Check if the secrets directory exists
	if _, err := os.Stat(sm.secretsDir); os.IsNotExist(err) {
		return secrets, nil
	}

	// Read the directory
	files, err := ioutil.ReadDir(sm.secretsDir)
	if err != nil {
		return nil, err
	}

	// Add each file to the list
	for _, f := range files {
		if !f.IsDir() {
			secrets = append(secrets, f.Name())
		}
	}

	return secrets, nil
}

// RemoveSecret removes a secret file
func (sm *SopsManager) RemoveSecret(name string) error {
	// Get the file path
	path := filepath.Join(sm.secretsDir, name)

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("secret file %s does not exist", name)
	}

	// Remove the file
	return os.Remove(path)
}

// EditSecret opens a secret in an editor for direct editing
func (sm *SopsManager) EditSecret(name string) error {
	// Get the file path
	path := filepath.Join(sm.secretsDir, name)

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("secret file %s does not exist", name)
	}

	// Use SOPS to edit the file
	cmd := exec.Command("sops", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utils.Logger.Info().Msgf("Opening secret %s for editing", name)
	return cmd.Run()
}