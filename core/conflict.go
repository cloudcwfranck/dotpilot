package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotpilot/utils"
)

// ConflictResolutionStrategy defines how conflicts should be resolved
type ConflictResolutionStrategy string

const (
	// StrategyInteractive prompts the user for each conflict
	StrategyInteractive ConflictResolutionStrategy = "interactive"
	// StrategyKeepLocal keeps the local version
	StrategyKeepLocal ConflictResolutionStrategy = "keep-local"
	// StrategyKeepRemote keeps the remote version
	StrategyKeepRemote ConflictResolutionStrategy = "keep-remote"
	// StrategyMerge attempts to merge changes
	StrategyMerge ConflictResolutionStrategy = "merge"
	// StrategyBackupBoth keeps both versions
	StrategyBackupBoth ConflictResolutionStrategy = "backup-both"
)

// ConflictFile represents a file with potential conflicts
type ConflictFile struct {
	LocalPath  string
	RemotePath string
	Target     string
	Diff       string
}

// ResolveConflicts identifies and resolves conflicts between local and remote files
func ResolveConflicts(dotpilotDir string, strategy ConflictResolutionStrategy) error {
	// Get the current list of conflicts
	conflicts, err := detectConflicts(dotpilotDir)
	if err != nil {
		return err
	}

	if len(conflicts) == 0 {
		utils.Logger.Info().Msg("No conflicts detected")
		return nil
	}

	utils.Logger.Info().Msgf("Detected %d conflicts", len(conflicts))

	// Process each conflict according to the strategy
	for _, conflict := range conflicts {
		utils.Logger.Info().Msgf("Resolving conflict for %s", conflict.Target)
		
		if err := resolveConflict(conflict, strategy); err != nil {
			utils.Logger.Error().Err(err).Msgf("Failed to resolve conflict for %s", conflict.Target)
			continue
		}
	}

	return nil
}

// detectConflicts identifies files with potential conflicts
func detectConflicts(dotpilotDir string) ([]ConflictFile, error) {
	var conflicts []ConflictFile

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Get current environment
	cfg := GetConfig()
	environment := cfg.CurrentEnvironment
	if environment == "" {
		environment = "default"
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Collect files that might have conflicts
	// We'll check files from all three layers:
	// 1. Common
	// 2. Environment-specific
	// 3. Machine-specific
	var allPaths []string

	// 1. Common files
	commonDir := filepath.Join(dotpilotDir, "common")
	commonFiles, err := collectFiles(commonDir)
	if err != nil {
		return nil, err
	}
	allPaths = append(allPaths, commonFiles...)

	// 2. Environment-specific files
	if environment != "" {
		envDir := filepath.Join(dotpilotDir, "envs", environment)
		envFiles, err := collectFiles(envDir)
		if err != nil {
			return nil, err
		}
		allPaths = append(allPaths, envFiles...)
	}

	// 3. Machine-specific files
	machineDir := filepath.Join(dotpilotDir, "machine", hostname)
	machineFiles, err := collectFiles(machineDir)
	if err != nil {
		return nil, err
	}
	allPaths = append(allPaths, machineFiles...)

	// Check each file for conflicts
	for _, path := range allPaths {
		// Get relative path from dotpilotDir
		relPath, err := filepath.Rel(dotpilotDir, path)
		if err != nil {
			utils.Logger.Error().Err(err).Msgf("Failed to get relative path for %s", path)
			continue
		}

		// Skip special files and directories
		if strings.HasPrefix(relPath, ".git") || relPath == "README.md" {
			continue
		}

		// Construct the target path in the home directory
		// We need to determine which part of the path structure this is in
		var targetPath string

		// Extract the type of file (common, env, machine)
		parts := strings.Split(relPath, string(os.PathSeparator))
		if len(parts) < 2 {
			continue
		}

		// Determine the base directory based on the file type
		switch parts[0] {
		case "common":
			targetPath = filepath.Join(home, filepath.Join(parts[2:]...))
		case "envs":
			if len(parts) < 3 {
				continue
			}
			targetPath = filepath.Join(home, filepath.Join(parts[3:]...))
		case "machine":
			if len(parts) < 3 {
				continue
			}
			targetPath = filepath.Join(home, filepath.Join(parts[3:]...))
		default:
			continue
		}

		// Check if the target exists and is not a symlink to our path
		targetInfo, err := os.Lstat(targetPath)
		if err != nil {
			// Target doesn't exist, no conflict
			continue
		}

		isSymlink := targetInfo.Mode()&os.ModeSymlink != 0
		if isSymlink {
			// Check if symlink points to our dotpilot path
			linkTarget, err := os.Readlink(targetPath)
			if err == nil && linkTarget == path {
				// No conflict, symlink points to our file
				continue
			}
		}

		// At this point, we have a potential conflict
		// Get the diff for the user to see
		diff, err := FileDiff(targetPath, path)
		if err != nil {
			utils.Logger.Warn().Err(err).Msgf("Failed to get diff for %s", targetPath)
			diff = "Unable to generate diff"
		}

		conflicts = append(conflicts, ConflictFile{
			LocalPath:  targetPath,
			RemotePath: path,
			Target:     targetPath,
			Diff:       diff,
		})
	}

	return conflicts, nil
}

// collectFiles recursively collects all files in a directory
func collectFiles(dir string) ([]string, error) {
	var files []string

	// Check if directory exists
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return files, nil
	}

	// Walk through the directory
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

// resolveConflict resolves a single conflict based on the strategy
func resolveConflict(conflict ConflictFile, strategy ConflictResolutionStrategy) error {
	switch strategy {
	case StrategyInteractive:
		return resolveInteractive(conflict)
	case StrategyKeepLocal:
		return resolveKeepLocal(conflict)
	case StrategyKeepRemote:
		return resolveKeepRemote(conflict)
	case StrategyMerge:
		return resolveMerge(conflict)
	case StrategyBackupBoth:
		return resolveBackupBoth(conflict)
	default:
		return fmt.Errorf("unknown conflict resolution strategy: %s", strategy)
	}
}

// resolveInteractive prompts the user to resolve the conflict
func resolveInteractive(conflict ConflictFile) error {
	fmt.Printf("\nConflict detected for %s\n", conflict.Target)
	fmt.Printf("Diff:\n%s\n", conflict.Diff)
	fmt.Println("\nHow would you like to resolve this conflict?")
	fmt.Println("1) Keep local version")
	fmt.Println("2) Keep remote version")
	fmt.Println("3) Merge changes (requires merge tool)")
	fmt.Println("4) View diff in external tool")
	fmt.Println("5) Edit file manually")
	fmt.Println("6) Keep both versions (create backup)")
	fmt.Println("7) Skip this conflict")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nEnter your choice (1-7): ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			return resolveKeepLocal(conflict)
		case "2":
			return resolveKeepRemote(conflict)
		case "3":
			return resolveMerge(conflict)
		case "4":
			if err := viewDiffExternal(conflict); err != nil {
				utils.Logger.Error().Err(err).Msg("Failed to view diff in external tool")
			}
			// After viewing, ask again
			continue
		case "5":
			if err := editFileManually(conflict); err != nil {
				utils.Logger.Error().Err(err).Msg("Failed to edit file manually")
			}
			// After editing, ask again
			continue
		case "6":
			return resolveBackupBoth(conflict)
		case "7":
			utils.Logger.Info().Msgf("Skipping conflict for %s", conflict.Target)
			return nil
		default:
			fmt.Println("Invalid choice, please try again")
		}
	}
}

// resolveKeepLocal keeps the local version and updates the remote file
func resolveKeepLocal(conflict ConflictFile) error {
	utils.Logger.Info().Msgf("Keeping local version for %s", conflict.Target)

	// Copy the local file to remote
	if err := copyFile(conflict.LocalPath, conflict.RemotePath, 0644); err != nil {
		return err
	}

	// Update the symlink
	if err := updateSymlink(conflict.RemotePath, conflict.LocalPath); err != nil {
		return err
	}

	return nil
}

// resolveKeepRemote keeps the remote version and updates the local file
func resolveKeepRemote(conflict ConflictFile) error {
	utils.Logger.Info().Msgf("Keeping remote version for %s", conflict.Target)

	// Backup the local file
	backupPath, err := BackupFile(conflict.LocalPath)
	if err != nil {
		return err
	}
	if backupPath != "" {
		utils.Logger.Info().Msgf("Backed up local file to %s", backupPath)
	}

	// Create symlink to remote file
	if err := updateSymlink(conflict.RemotePath, conflict.LocalPath); err != nil {
		return err
	}

	return nil
}

// resolveMerge attempts to merge changes using an external merge tool
func resolveMerge(conflict ConflictFile) error {
	utils.Logger.Info().Msgf("Attempting to merge changes for %s", conflict.Target)

	// Check if we have common merge tools installed
	mergeTools := []string{"meld", "kdiff3", "vimdiff", "code -d"}
	selectedTool := ""

	for _, tool := range mergeTools {
		// Extract the command (part before any space)
		cmd := strings.Split(tool, " ")[0]
		_, err := exec.LookPath(cmd)
		if err == nil {
			selectedTool = tool
			break
		}
	}

	if selectedTool == "" {
		return fmt.Errorf("no merge tool found, please install a merge tool (meld, kdiff3, vimdiff)")
	}

	// Create a temporary file for the merged result
	mergedFile, err := os.CreateTemp("", "dotpilot-merge-*")
	if err != nil {
		return err
	}
	mergedPath := mergedFile.Name()
	mergedFile.Close()

	// Copy remote file to merged file as a starting point
	if err := copyFile(conflict.RemotePath, mergedPath, 0644); err != nil {
		os.Remove(mergedPath)
		return err
	}

	// Build the merge command
	var cmdParts []string
	if selectedTool == "vimdiff" {
		cmdParts = []string{selectedTool, conflict.LocalPath, mergedPath, conflict.RemotePath}
	} else {
		// General format for most merge tools
		cmdParts = strings.Split(selectedTool, " ")
		cmdParts = append(cmdParts, conflict.LocalPath, mergedPath, conflict.RemotePath)
	}

	// Execute the merge tool
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utils.Logger.Info().Msgf("Launching merge tool: %s", strings.Join(cmdParts, " "))
	if err := cmd.Run(); err != nil {
		os.Remove(mergedPath)
		return err
	}

	// After the merge tool completes, copy the merged result to both local and remote
	if err := copyFile(mergedPath, conflict.LocalPath, 0644); err != nil {
		os.Remove(mergedPath)
		return err
	}

	if err := copyFile(mergedPath, conflict.RemotePath, 0644); err != nil {
		os.Remove(mergedPath)
		return err
	}

	// Clean up
	os.Remove(mergedPath)

	// Update the symlink
	if err := updateSymlink(conflict.RemotePath, conflict.LocalPath); err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Successfully merged changes for %s", conflict.Target)
	return nil
}

// resolveBackupBoth keeps both versions with the remote in dotpilot and the local as-is
func resolveBackupBoth(conflict ConflictFile) error {
	utils.Logger.Info().Msgf("Keeping both versions for %s", conflict.Target)

	// Generate a unique backup name for the remote file
	backupName := fmt.Sprintf("%s.local.%s", filepath.Base(conflict.RemotePath), time.Now().Format("20060102150405"))
	backupDir := filepath.Dir(conflict.RemotePath)
	backupPath := filepath.Join(backupDir, backupName)

	// Copy the local file to the backup location in dotpilot
	if err := copyFile(conflict.LocalPath, backupPath, 0644); err != nil {
		return err
	}

	utils.Logger.Info().Msgf("Created backup of local file at %s", backupPath)
	utils.Logger.Info().Msgf("Original remote file remains at %s", conflict.RemotePath)
	utils.Logger.Info().Msgf("Local file remains at %s", conflict.LocalPath)

	return nil
}

// viewDiffExternal shows the diff in an external diff tool
func viewDiffExternal(conflict ConflictFile) error {
	// Check for available diff tools
	diffTools := []string{"meld", "kdiff3", "vimdiff", "code -d", "diff -u"}
	selectedTool := ""

	for _, tool := range diffTools {
		// Extract the command (part before any space)
		cmd := strings.Split(tool, " ")[0]
		_, err := exec.LookPath(cmd)
		if err == nil {
			selectedTool = tool
			break
		}
	}

	if selectedTool == "" {
		// Fallback to printing the diff
		fmt.Printf("Diff between %s and %s:\n%s\n", conflict.LocalPath, conflict.RemotePath, conflict.Diff)
		return nil
	}

	// Build the diff command
	cmdParts := strings.Split(selectedTool, " ")
	cmdParts = append(cmdParts, conflict.LocalPath, conflict.RemotePath)

	// Execute the diff tool
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utils.Logger.Info().Msgf("Launching diff tool: %s", strings.Join(cmdParts, " "))
	return cmd.Run()
}

// editFileManually opens the file in an editor for manual editing
func editFileManually(conflict ConflictFile) error {
	// Try to find an editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Try common editors
		editors := []string{"nano", "vim", "vi", "emacs", "code"}
		for _, ed := range editors {
			_, err := exec.LookPath(ed)
			if err == nil {
				editor = ed
				break
			}
		}
	}

	if editor == "" {
		return fmt.Errorf("no editor found, please set the EDITOR environment variable")
	}

	// Create a temporary file with the content
	tmpFile, err := os.CreateTemp("", "dotpilot-edit-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Copy the remote file as a starting point
	if err := copyFile(conflict.RemotePath, tmpPath, 0644); err != nil {
		os.Remove(tmpPath)
		return err
	}

	// Open the editor
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utils.Logger.Info().Msgf("Opening %s in %s", tmpPath, editor)
	if err := cmd.Run(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	// After editing, ask if the user wants to use this version
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Use this edited version? (y/n): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" || response == "yes" {
		// Copy the edited file to both local and remote
		if err := copyFile(tmpPath, conflict.LocalPath, 0644); err != nil {
			os.Remove(tmpPath)
			return err
		}

		if err := copyFile(tmpPath, conflict.RemotePath, 0644); err != nil {
			os.Remove(tmpPath)
			return err
		}

		// Update the symlink
		if err := updateSymlink(conflict.RemotePath, conflict.LocalPath); err != nil {
			os.Remove(tmpPath)
			return err
		}

		utils.Logger.Info().Msgf("Applied edited version to %s", conflict.Target)
	} else {
		utils.Logger.Info().Msg("Edited version discarded")
	}

	// Clean up
	os.Remove(tmpPath)
	return nil
}

// updateSymlink creates or updates a symlink
func updateSymlink(source, target string) error {
	// Remove the target if it exists
	_, err := os.Lstat(target)
	if err == nil {
		if err := os.Remove(target); err != nil {
			return err
		}
	}

	// Create symlink
	return os.Symlink(source, target)
}