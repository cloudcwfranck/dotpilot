package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotpilot/utils"
)

// This demo showcases a real-world scenario of using DotPilot
// to manage dotfiles across multiple machines
func main() {
	// Get home directory for the demo
	home, _ := os.UserHomeDir()
	
	fmt.Println("DotPilot Real-World Usage Demo")
	fmt.Println("==============================")
	fmt.Println("This demo shows how to use DotPilot to manage dotfiles across multiple machines.")
	
	// 1. Setting up DotPilot on a new machine
	fmt.Println("\n1. Setting up DotPilot on a new developer machine")
	fmt.Println("------------------------------------------------")
	
	// Creating a demo directory structure
	demoDir := filepath.Join(home, ".dotpilot_demo")
	os.MkdirAll(demoDir, 0755)
	
	// Initialize Operation Manager for progress indicators
	initOp := utils.NewOperationManager()
	
	// Initialize DotPilot with a remote repository
	fmt.Println("$ dotpilot init --remote https://github.com/username/dotfiles.git --env dev")
	op1 := initOp.AddOperation("init", "Initializing dotpilot repository...", utils.Spinner)
	op1.Start()
	time.Sleep(2 * time.Second)
	op1.SetState(utils.Success)
	op1.Stop()
	
	fmt.Printf("✓ Created dotpilot directory at %s\n", filepath.Join(home, ".dotpilot"))
	fmt.Printf("✓ Cloned repository from https://github.com/username/dotfiles.git\n")
	fmt.Printf("✓ Created directory structure (common, envs/dev, machine/dev-laptop)\n")
	
	// 2. Track essential dotfiles
	fmt.Println("\n2. Tracking essential configuration files")
	fmt.Println("---------------------------------------")
	trackOp := utils.NewOperationManager()
	
	// Track .zshrc
	fmt.Println("$ dotpilot track ~/.zshrc --env common")
	op2 := trackOp.AddOperation("track-zshrc", "Tracking .zshrc...", utils.Spinner)
	op2.Start()
	time.Sleep(1 * time.Second)
	op2.SetState(utils.Success)
	op2.Stop()
	fmt.Printf("✓ Added .zshrc to common environment\n")
	
	// Track .gitconfig
	fmt.Println("$ dotpilot track ~/.gitconfig --env dev")
	op3 := trackOp.AddOperation("track-gitconfig", "Tracking .gitconfig...", utils.Spinner)
	op3.Start()
	time.Sleep(800 * time.Millisecond)
	op3.SetState(utils.Success)
	op3.Stop()
	fmt.Printf("✓ Added .gitconfig to dev environment\n")
	
	// Track .vimrc
	fmt.Println("$ dotpilot track ~/.vimrc --env common")
	op4 := trackOp.AddOperation("track-vimrc", "Tracking .vimrc...", utils.Spinner)
	op4.Start()
	time.Sleep(900 * time.Millisecond)
	op4.SetState(utils.Success)
	op4.Stop()
	fmt.Printf("✓ Added .vimrc to common environment\n")
	
	// Track VS Code settings
	fmt.Println("$ dotpilot track ~/.config/Code/User/settings.json --env dev")
	op5 := trackOp.AddOperation("track-vscode", "Tracking VS Code settings...", utils.Spinner)
	op5.Start()
	time.Sleep(1200 * time.Millisecond)
	op5.SetState(utils.Success)
	op5.Stop()
	fmt.Printf("✓ Added VS Code settings to dev environment\n")
	
	// 3. Sync with remote repository
	fmt.Println("\n3. Syncing changes with remote repository")
	fmt.Println("----------------------------------------")
	fmt.Println("$ dotpilot sync")
	
	syncOp := utils.NewOperationManager()
	
	// Commit changes
	op6 := syncOp.AddOperation("commit", "Auto-committing changes before sync...", utils.Spinner)
	op6.Start()
	time.Sleep(1500 * time.Millisecond)
	op6.SetState(utils.Success)
	op6.Stop()
	
	// Pull from remote
	op7 := syncOp.AddOperation("pull", "Pulling changes from remote...", utils.Bounce)
	op7.Start()
	time.Sleep(2000 * time.Millisecond)
	op7.SetState(utils.Success)
	op7.Stop()
	
	// Apply configurations
	op8 := syncOp.AddOperation("apply", "Applying configurations...", utils.Bar)
	op8.Start()
	
	// Simulate progress
	for i := 0; i <= 100; i += 10 {
		op8.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op8.SetState(utils.Success)
	op8.Stop()
	
	// Push to remote
	op9 := syncOp.AddOperation("push", "Pushing changes to remote...", utils.Bounce)
	op9.Start()
	time.Sleep(1800 * time.Millisecond)
	op9.SetState(utils.Success)
	op9.Stop()
	
	fmt.Printf("✓ All changes committed\n")
	fmt.Printf("✓ Changes pulled from remote\n")
	fmt.Printf("✓ Dotfiles synchronized\n")
	fmt.Printf("✓ Changes pushed to remote repository\n")
	
	// 4. Setting up a new machine using bootstrap
	fmt.Println("\n4. Setting up dotfiles on a new machine")
	fmt.Println("--------------------------------------")
	fmt.Println("On a new machine after installing DotPilot:")
	fmt.Println("$ dotpilot init --remote https://github.com/username/dotfiles.git --env dev")
	fmt.Println("$ dotpilot bootstrap")
	
	bootstrapOp := utils.NewOperationManager()
	
	// Apply common configurations
	op10 := bootstrapOp.AddOperation("common", "Applying common dotfiles...", utils.Bar)
	op10.Start()
	
	for i := 0; i <= 100; i += 10 {
		op10.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op10.SetState(utils.Success)
	op10.Stop()
	
	// Apply environment-specific configurations
	op11 := bootstrapOp.AddOperation("env", "Applying dev environment dotfiles...", utils.Bar)
	op11.Start()
	
	for i := 0; i <= 100; i += 10 {
		op11.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op11.SetState(utils.Success)
	op11.Stop()
	
	// Apply machine-specific configurations
	op12 := bootstrapOp.AddOperation("machine", "Applying machine-specific dotfiles...", utils.Bar)
	op12.Start()
	
	hostname, _ := os.Hostname()
	fmt.Printf("   Detected hostname: %s\n", hostname)
	
	for i := 0; i <= 100; i += 10 {
		op12.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op12.SetState(utils.Success)
	op12.Stop()
	
	// Run setup scripts
	op13 := bootstrapOp.AddOperation("scripts", "Running setup scripts...", utils.Pulse)
	op13.Start()
	time.Sleep(2 * time.Second)
	op13.SetState(utils.Success)
	op13.Stop()
	
	fmt.Printf("✓ Applied common dotfiles (.zshrc, .vimrc)\n")
	fmt.Printf("✓ Applied dev environment dotfiles (.gitconfig, VS Code settings)\n")
	fmt.Printf("✓ Applied machine-specific dotfiles\n")
	fmt.Printf("✓ Ran setup scripts (package installations)\n")
	
	// 5. Managing dotfiles with sensitive information
	fmt.Println("\n5. Managing dotfiles with sensitive information")
	fmt.Println("-------------------------------------------")
	fmt.Println("$ dotpilot secrets add ~/.aws/credentials --name aws_keys")
	
	secretsOp := utils.NewOperationManager()
	op14 := secretsOp.AddOperation("encrypt", "Encrypting AWS credentials...", utils.Dots)
	op14.Start()
	time.Sleep(2 * time.Second)
	op14.SetState(utils.Success)
	op14.Stop()
	
	fmt.Printf("✓ Encrypted AWS credentials securely\n")
	fmt.Printf("✓ Added encrypted file to dotpilot repository\n")
	fmt.Printf("✓ Created symlink to original location\n")
	
	// Check the status
	fmt.Println("\n6. Checking dotfiles status")
	fmt.Println("-------------------------")
	fmt.Println("$ dotpilot status")
	
	statusOp := utils.NewOperationManager()
	op15 := statusOp.AddOperation("status", "Checking dotpilot status...", utils.Spinner)
	op15.Start()
	time.Sleep(1 * time.Second)
	op15.SetState(utils.Success)
	op15.Stop()
	
	fmt.Printf("Repository: https://github.com/username/dotfiles.git\n")
	fmt.Printf("Current environment: dev\n")
	fmt.Printf("Machine: %s\n", hostname)
	fmt.Printf("\nTracked files:\n")
	fmt.Printf("  common/\n")
	fmt.Printf("    ├── .zshrc\n")
	fmt.Printf("    └── .vimrc\n")
	fmt.Printf("  envs/dev/\n")
	fmt.Printf("    ├── .gitconfig\n")
	fmt.Printf("    └── .config/Code/User/settings.json\n")
	fmt.Printf("  secrets/\n")
	fmt.Printf("    └── aws_keys.enc\n")
	fmt.Printf("\nGit status: Clean\n")
	
	fmt.Println("\n✨ DotPilot Real-World Demo Completed!")
	fmt.Println("Note: Clean up the demo directory...")
	
	// Clean up demo directory if it was created
	os.RemoveAll(demoDir)
}