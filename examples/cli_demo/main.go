package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotpilot/utils"
)

func main() {
	// Demo for DotPilot CLI
	fmt.Println("DotPilot CLI Demo")
	fmt.Println("================")
	
	// Create demo directory structure
	home, _ := os.UserHomeDir()
	dotpilotDir := filepath.Join(home, ".dotpilot")
	
	fmt.Println("1. Running 'dotpilot init' command")
	fmt.Println("----------------------------------")
	
	// Simulate init command with progress indicator
	initOp := utils.NewOperationManager()
	op1 := initOp.AddOperation("init", "Initializing dotpilot repository...", utils.Spinner)
	op1.Start()
	
	time.Sleep(2 * time.Second)
	op1.SetState(utils.Success)
	op1.Stop()
	
	fmt.Printf("✓ Created dotpilot directory at %s\n", dotpilotDir)
	fmt.Printf("✓ Initialized Git repository\n")
	fmt.Printf("✓ Created directory structure (common, envs, machine)\n")
	
	// Simulate track command
	fmt.Println("\n2. Running 'dotpilot track ~/.bashrc' command")
	fmt.Println("-------------------------------------------")
	
	trackOp := utils.NewOperationManager()
	op2 := trackOp.AddOperation("track", "Tracking ~/.bashrc...", utils.Spinner)
	op2.Start()
	
	time.Sleep(1 * time.Second)
	op2.SetState(utils.Success)
	op2.Stop()
	
	fmt.Printf("✓ Copied ~/.bashrc to %s/common/.bashrc\n", dotpilotDir)
	fmt.Printf("✓ Created symlink from %s/common/.bashrc to ~/.bashrc\n", dotpilotDir)
	fmt.Printf("✓ Added .bashrc to tracking list\n")
	
	// Simulate sync command
	fmt.Println("\n3. Running 'dotpilot sync' command")
	fmt.Println("--------------------------------")
	
	syncOp := utils.NewOperationManager()
	op3 := syncOp.AddOperation("commit", "Auto-committing changes before sync...", utils.Spinner)
	op3.Start()
	
	time.Sleep(1 * time.Second)
	op3.SetState(utils.Success)
	op3.Stop()
	
	op4 := syncOp.AddOperation("pull", "Pulling changes from remote...", utils.Bounce)
	op4.Start()
	
	time.Sleep(2 * time.Second)
	op4.SetState(utils.Success)
	op4.Stop()
	
	op5 := syncOp.AddOperation("apply", "Applying configurations...", utils.Bar)
	op5.Start()
	
	// Simulate progress for bar
	for i := 0; i <= 100; i += 5 {
		op5.UpdateProgress(i, 100)
		time.Sleep(50 * time.Millisecond)
	}
	
	op5.SetState(utils.Success)
	op5.Stop()
	
	op6 := syncOp.AddOperation("push", "Pushing changes to remote...", utils.Bounce)
	op6.Start()
	
	time.Sleep(1 * time.Second)
	op6.SetState(utils.Success)
	op6.Stop()
	
	fmt.Printf("✓ Changes committed\n")
	fmt.Printf("✓ Changes pulled from remote\n")
	fmt.Printf("✓ Configurations applied\n")
	fmt.Printf("✓ Changes pushed to remote\n")
	
	// Simulate bootstrap command
	fmt.Println("\n4. Running 'dotpilot bootstrap' command")
	fmt.Println("--------------------------------------")
	
	bootstrapOp := utils.NewOperationManager()
	
	op7 := bootstrapOp.AddOperation("common", "Applying common dotfiles...", utils.Bar)
	op7.Start()
	
	for i := 0; i <= 100; i += 10 {
		op7.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op7.SetState(utils.Success)
	op7.Stop()
	
	op8 := bootstrapOp.AddOperation("env", "Applying environment-specific dotfiles...", utils.Bar)
	op8.Start()
	
	for i := 0; i <= 100; i += 10 {
		op8.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op8.SetState(utils.Success)
	op8.Stop()
	
	op9 := bootstrapOp.AddOperation("machine", "Applying machine-specific dotfiles...", utils.Bar)
	op9.Start()
	
	hostname, _ := os.Hostname()
	fmt.Printf("   Detected hostname: %s\n", hostname)
	
	for i := 0; i <= 100; i += 10 {
		op9.UpdateProgress(i, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	op9.SetState(utils.Success)
	op9.Stop()
	
	op10 := bootstrapOp.AddOperation("scripts", "Running setup scripts...", utils.Pulse)
	op10.Start()
	
	time.Sleep(2 * time.Second)
	op10.SetState(utils.Success)
	op10.Stop()
	
	fmt.Printf("✓ Applied common dotfiles\n")
	fmt.Printf("✓ Applied environment-specific dotfiles\n")
	fmt.Printf("✓ Applied machine-specific dotfiles\n")
	fmt.Printf("✓ Ran setup scripts\n")
	
	fmt.Println("\n✨ DotPilot CLI Demo Completed!")
}