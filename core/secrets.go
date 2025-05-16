package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dotpilot/utils"
	"golang.org/x/crypto/pbkdf2"
)

// SecretManager handles encrypted secrets
type SecretManager struct {
	dotpilotDir string
	keyFile     string
	secretsDir  string
	useGPG      bool
}

// NewSecretManager creates a new secret manager
func NewSecretManager(dotpilotDir string) *SecretManager {
	return &SecretManager{
		dotpilotDir: dotpilotDir,
		keyFile:     filepath.Join(dotpilotDir, ".secret_key"),
		secretsDir:  filepath.Join(dotpilotDir, "secrets"),
		useGPG:      isGPGAvailable(),
	}
}

// isGPGAvailable checks if GPG is available on the system
func isGPGAvailable() bool {
	_, err := exec.LookPath("gpg")
	return err == nil
}

// Initialize sets up the secrets directory and encryption keys
func (sm *SecretManager) Initialize() error {
	// Create secrets directory if it doesn't exist
	if err := os.MkdirAll(sm.secretsDir, 0700); err != nil {
		return err
	}

	// If using GPG, no need to create a key file
	if sm.useGPG {
		utils.Logger.Info().Msg("Using GPG for secrets encryption")
		return nil
	}

	// Check if key file exists
	if _, err := os.Stat(sm.keyFile); os.IsNotExist(err) {
		// Generate a new key
		key := make([]byte, 32) // 256-bit key
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return err
		}

		// Save the key to file with restricted permissions
		encodedKey := base64.StdEncoding.EncodeToString(key)
		if err := ioutil.WriteFile(sm.keyFile, []byte(encodedKey), 0600); err != nil {
			return err
		}

		utils.Logger.Info().Msg("Generated new encryption key")
	}

	utils.Logger.Info().Msg("Secret manager initialized")
	return nil
}

// EncryptFile encrypts a file and stores it in the secrets directory
func (sm *SecretManager) EncryptFile(srcPath, name string) error {
	// Create destination path
	destPath := filepath.Join(sm.secretsDir, name)

	// Read the source file
	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Use GPG if available
	if sm.useGPG {
		return sm.encryptWithGPG(data, destPath)
	}

	// Use AES otherwise
	return sm.encryptWithAES(data, destPath)
}

// DecryptFile decrypts a file from the secrets directory
func (sm *SecretManager) DecryptFile(name, destPath string) error {
	// Get the source path
	srcPath := filepath.Join(sm.secretsDir, name)

	// Check if the file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("secret file %s does not exist", name)
	}

	// Use GPG if available
	if sm.useGPG {
		return sm.decryptWithGPG(srcPath, destPath)
	}

	// Use AES otherwise
	return sm.decryptWithAES(srcPath, destPath)
}

// ListSecrets returns a list of all secret files
func (sm *SecretManager) ListSecrets() ([]string, error) {
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
func (sm *SecretManager) RemoveSecret(name string) error {
	// Get the file path
	path := filepath.Join(sm.secretsDir, name)

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("secret file %s does not exist", name)
	}

	// Remove the file
	return os.Remove(path)
}

// encryptWithGPG encrypts data using GPG
func (sm *SecretManager) encryptWithGPG(data []byte, destPath string) error {
	// Get GPG recipient (default to user's GPG ID)
	recipient, err := getGPGRecipient()
	if err != nil {
		return err
	}

	// Create a temp file for input
	tmpfile, err := ioutil.TempFile("", "dotpilot-gpg-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	// Write data to temp file
	if _, err := tmpfile.Write(data); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	// Use GPG to encrypt
	cmd := exec.Command("gpg", "--encrypt", "--recipient", recipient, "--output", destPath, tmpfile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gpg encryption failed: %s - %s", err, string(output))
	}

	utils.Logger.Info().Msgf("Encrypted file with GPG to %s", destPath)
	return nil
}

// decryptWithGPG decrypts a file using GPG
func (sm *SecretManager) decryptWithGPG(srcPath, destPath string) error {
	// Use GPG to decrypt
	cmd := exec.Command("gpg", "--decrypt", "--output", destPath, srcPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gpg decryption failed: %s - %s", err, string(output))
	}

	// Set correct permissions on the output file
	if err := os.Chmod(destPath, 0600); err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Decrypted file with GPG to %s", destPath)
	return nil
}

// getGPGRecipient gets the default GPG key ID
func getGPGRecipient() (string, error) {
	// Run gpg --list-keys to get the default key
	cmd := exec.Command("gpg", "--list-keys")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse the output to find a key ID
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if strings.Contains(line, "pub") && i+1 < len(lines) {
			// The key ID is usually in the format XXXXXXXX in the next line
			nextLine := lines[i+1]
			fields := strings.Fields(nextLine)
			if len(fields) > 0 {
				// Usually the last field is the email or name, the key ID is before that
				if len(fields) > 1 {
					return fields[len(fields)-1], nil
				}
				return fields[0], nil
			}
		}
	}

	return "", errors.New("unable to find GPG key, please specify recipient manually")
}

// encryptWithAES encrypts data using AES-256-GCM
func (sm *SecretManager) encryptWithAES(data []byte, destPath string) error {
	// Get the encryption key
	key, err := sm.getEncryptionKey()
	if err != nil {
		return err
	}

	// Generate a random salt for PBKDF2
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	// Derive a key using PBKDF2
	derivedKey := pbkdf2.Key(key, salt, 4096, 32, sha256.New)

	// Create a new AES cipher block
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Combine salt + nonce + ciphertext for storage
	combined := make([]byte, len(salt)+len(nonce)+len(ciphertext))
	copy(combined, salt)
	copy(combined[len(salt):], nonce)
	copy(combined[len(salt)+len(nonce):], ciphertext)

	// Encode as base64
	encoded := base64.StdEncoding.EncodeToString(combined)

	// Write to file
	if err := ioutil.WriteFile(destPath, []byte(encoded), 0600); err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Encrypted file with AES to %s", destPath)
	return nil
}

// decryptWithAES decrypts a file using AES-256-GCM
func (sm *SecretManager) decryptWithAES(srcPath, destPath string) error {
	// Get the encryption key
	key, err := sm.getEncryptionKey()
	if err != nil {
		return err
	}

	// Read the encrypted data
	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}

	// Extract salt, nonce, and ciphertext
	if len(decoded) < 32 {
		return errors.New("invalid encrypted data format")
	}

	salt := decoded[:16]
	nonce := decoded[16:32]
	ciphertext := decoded[32:]

	// Derive the key using PBKDF2
	derivedKey := pbkdf2.Key(key, salt, 4096, 32, sha256.New)

	// Create a new AES cipher block
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	// Write to file
	if err := ioutil.WriteFile(destPath, plaintext, 0600); err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Decrypted file with AES to %s", destPath)
	return nil
}

// getEncryptionKey reads the encryption key from the key file
func (sm *SecretManager) getEncryptionKey() ([]byte, error) {
	// Read the key file
	data, err := ioutil.ReadFile(sm.keyFile)
	if err != nil {
		return nil, err
	}

	// Decode from base64
	key, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	return key, nil
}