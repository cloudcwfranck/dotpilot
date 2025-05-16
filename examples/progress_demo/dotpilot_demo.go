package main

import (
	"fmt"
	"time"
	
	"github.com/dotpilot/utils"
)

// Progress Indicator Demo
func main() {
	fmt.Println("DotPilot: Animated Progress Indicators Demo")
	fmt.Println("===========================================")
	
	// Demo 1: Simulating Git Operations
	fmt.Println("\n👉 Syncing dotfiles with remote repository...")
	
	// Step 1: Commit changes
	time.Sleep(500 * time.Millisecond)
	manager := utils.NewOperationManager()
	
	spinner := manager.AddOperation("commit", "Auto-committing changes before sync", utils.Spinner)
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.Stop()
	fmt.Println("✓ Changes committed")
	
	// Step 2: Pull changes
	time.Sleep(500 * time.Millisecond)
	bounce := manager.AddOperation("pull", "Pulling changes from remote", utils.Bounce)
	bounce.Start()
	time.Sleep(2 * time.Second)
	bounce.Stop()
	fmt.Println("✓ Changes pulled from remote")
	
	// Step 3: Apply configurations
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\n👉 Applying configurations...")
	bar := manager.AddOperation("apply", "Applying configurations", utils.Bar)
	bar.Start()
	
	// Simulate progress
	bar.SimulateProgress(3) // 3 seconds duration
	time.Sleep(3 * time.Second)
	bar.Stop()
	fmt.Println("✓ Configurations applied")
	
	// Step 4: Push changes
	time.Sleep(500 * time.Millisecond)
	bounce2 := manager.AddOperation("push", "Pushing changes to remote", utils.Bounce)
	bounce2.Start()
	time.Sleep(2 * time.Second)
	bounce2.Stop()
	fmt.Println("✓ Changes pushed to remote")
	
	// Demo 2: Encrypting Secrets
	time.Sleep(1 * time.Second)
	fmt.Println("\n👉 Encrypting sensitive configuration files...")
	
	dots := manager.AddOperation("encrypt", "Encrypting ~/.aws/credentials", utils.Dots)
	dots.Start()
	time.Sleep(3 * time.Second)
	dots.Stop()
	fmt.Println("✓ Credentials encrypted successfully")
	
	// Demo 3: Multiple concurrent operations
	time.Sleep(1 * time.Second)
	fmt.Println("\n👉 Performing multiple concurrent operations...")
	
	op1 := manager.AddOperation("scan", "Scanning for dotfiles", utils.Spinner)
	op2 := manager.AddOperation("analyze", "Analyzing configurations", utils.Bar)
	op3 := manager.AddOperation("check", "Checking remote status", utils.Bounce)
	
	op1.Start()
	time.Sleep(700 * time.Millisecond)
	op2.Start()
	time.Sleep(700 * time.Millisecond)
	op3.Start()
	
	// Simulate progress for the bar
	op2.SimulateProgress(5) // 5 seconds
	
	time.Sleep(5 * time.Second)
	
	op1.Stop()
	fmt.Println("✓ Dotfiles scan complete")
	time.Sleep(300 * time.Millisecond)
	op2.Stop()
	fmt.Println("✓ Configuration analysis complete")
	time.Sleep(300 * time.Millisecond)
	op3.Stop()
	fmt.Println("✓ Remote status checked")
	
	fmt.Println("\n✨ DotPilot operations completed successfully!")
}