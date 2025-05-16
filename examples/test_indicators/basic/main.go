package main

import (
	"fmt"
	"time"

	"github.com/dotpilot/utils"
)

// Test basic progress indicator functionality
func main() {
	fmt.Println("Testing animated progress indicators...")
	
	// Test spinner
	fmt.Println("\nTesting Spinner style:")
	spinner := utils.NewProgressIndicator("Processing data", utils.Spinner)
	spinner.Start()
	time.Sleep(3 * time.Second)
	spinner.SetState(utils.Success)
	time.Sleep(1 * time.Second)
	spinner.Stop()
	
	// Test progress bar
	fmt.Println("\nTesting Progress Bar style:")
	bar := utils.NewProgressIndicator("Downloading files", utils.Bar)
	bar.Start()
	// Simulate progress
	for i := 0; i <= 100; i += 5 {
		bar.UpdateProgress(i)
		time.Sleep(100 * time.Millisecond)
	}
	bar.Stop()
	
	// Test bouncing indicator
	fmt.Println("\nTesting Bounce style:")
	bounce := utils.NewProgressIndicator("Synchronizing", utils.Bounce)
	bounce.Start()
	time.Sleep(3 * time.Second)
	bounce.Stop()
	
	// Test dots indicator
	fmt.Println("\nTesting Dots style:")
	dots := utils.NewProgressIndicator("Loading", utils.Dots)
	dots.Start()
	time.Sleep(3 * time.Second)
	dots.Stop()
	
	// Test multiple concurrent indicators
	fmt.Println("\nTesting multiple concurrent indicators:")
	
	manager := utils.NewOperationManager()
	
	op1 := manager.AddOperation("spinner", "Operation 1", utils.Spinner)
	op2 := manager.AddOperation("bar", "Operation 2", utils.Bar)
	op3 := manager.AddOperation("bounce", "Operation 3", utils.Bounce)
	
	// Start all operations
	op1.Start()
	op2.Start()
	op3.Start()
	
	// Simulate progress for the bar
	go func() {
		for i := 0; i <= 100; i += 2 {
			op2.UpdateProgress(i, 100)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	// Let the indicators run
	time.Sleep(5 * time.Second)
	
	// Stop all operations
	op1.Stop()
	op2.Stop()
	op3.Stop()
	
	fmt.Println("\nProgress indicator tests completed!")
}