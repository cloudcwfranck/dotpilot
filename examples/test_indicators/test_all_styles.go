package main

import (
        "fmt"
        "time"

        "github.com/dotpilot/utils"
)

// Test all progress indicator styles including the new ones
func main() {
        fmt.Println("Testing all progress indicator styles...")
        
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
                bar.UpdateProgress(i) // This works for the ProgressIndicator directly
                if i > 75 {
                        bar.SetState(utils.Success)
                } else if i > 50 {
                        bar.SetState(utils.Info)
                } else if i > 25 {
                        bar.SetState(utils.Warning)
                }
                time.Sleep(100 * time.Millisecond)
        }
        bar.Stop()
        
        // Test bouncing indicator
        fmt.Println("\nTesting Bounce style:")
        bounce := utils.NewProgressIndicator("Synchronizing", utils.Bounce)
        bounce.Start()
        time.Sleep(1 * time.Second)
        bounce.SetState(utils.Info)
        time.Sleep(1 * time.Second)
        bounce.SetState(utils.Success)
        time.Sleep(1 * time.Second)
        bounce.Stop()
        
        // Test dots indicator
        fmt.Println("\nTesting Dots style:")
        dots := utils.NewProgressIndicator("Loading", utils.Dots)
        dots.Start()
        time.Sleep(1 * time.Second)
        dots.SetState(utils.Warning)
        time.Sleep(1 * time.Second)
        dots.SetState(utils.Error)
        time.Sleep(1 * time.Second)
        dots.Stop()
        
        // Test pulse indicator (new)
        fmt.Println("\nTesting Pulse style:")
        pulse := utils.NewProgressIndicator("Encrypting files", utils.Pulse)
        pulse.Start()
        time.Sleep(1 * time.Second)
        pulse.SetState(utils.Info)
        time.Sleep(1 * time.Second)
        pulse.SetState(utils.Success)
        time.Sleep(1 * time.Second)
        pulse.Stop()
        
        // Test rainbow indicator (new)
        fmt.Println("\nTesting Rainbow style:")
        rainbow := utils.NewProgressIndicator("Processing critical data", utils.Rainbow)
        rainbow.Start()
        time.Sleep(3 * time.Second)
        rainbow.Stop()
        
        // Test multiple indicators with all styles
        fmt.Println("\nTesting multiple concurrent indicators with all styles:")
        
        manager := utils.NewOperationManager()
        
        op1 := manager.AddOperation("spinner", "Operation 1 (Spinner)", utils.Spinner)
        op2 := manager.AddOperation("bar", "Operation 2 (Bar)", utils.Bar)
        op3 := manager.AddOperation("bounce", "Operation 3 (Bounce)", utils.Bounce)
        op4 := manager.AddOperation("dots", "Operation 4 (Dots)", utils.Dots)
        op5 := manager.AddOperation("pulse", "Operation 5 (Pulse)", utils.Pulse)
        op6 := manager.AddOperation("rainbow", "Operation 6 (Rainbow)", utils.Rainbow)
        
        // Start all operations with different delays
        op1.Start()
        time.Sleep(300 * time.Millisecond)
        op2.Start()
        time.Sleep(300 * time.Millisecond)
        op3.Start()
        time.Sleep(300 * time.Millisecond)
        op4.Start()
        time.Sleep(300 * time.Millisecond)
        op5.Start()
        time.Sleep(300 * time.Millisecond)
        op6.Start()
        
        // Simulate progress for the bar
        go func() {
                for i := 0; i <= 100; i += 2 {
                        op2.UpdateProgress(i, 100)
                        time.Sleep(100 * time.Millisecond)
                }
        }()
        
        // Set different states
        time.Sleep(2 * time.Second)
        op1.SetState(utils.Success)
        op3.SetState(utils.Warning)
        op4.SetState(utils.Error)
        op5.SetState(utils.Info)
        
        time.Sleep(3 * time.Second)
        
        // Stop all operations
        op1.Stop()
        op2.Stop()
        op3.Stop()
        op4.Stop()
        op5.Stop()
        op6.Stop()
        
        fmt.Println("\nProgress indicator style tests completed!")
}