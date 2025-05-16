package main

import (
        "fmt"
        "os"
        "time"

        "github.com/dotpilot/utils"
)

func main() {
        // Create a demo operation manager
        manager := utils.NewOperationManager()

        fmt.Println("DotPilot: Bootstrap Command Demo")
        fmt.Println("=================================")
        
        // 1. Initialize a mock dotfiles repository
        fmt.Println("👉 Initializing dotfiles repository...")
        initOp := manager.AddOperation("init", "Initializing repository...", utils.Spinner)
        initOp.Start()
        
        time.Sleep(2 * time.Second)
        initOp.SetState(utils.Success)
        initOp.Stop()
        
        // 2. Apply common configs
        fmt.Println("👉 Applying common configurations...")
        commonOp := manager.AddOperation("common", "Applying common dotfiles...", utils.Bar)
        commonOp.Start()
        
        for i := 0; i <= 100; i += 5 {
                commonOp.UpdateProgress(i, 100)
                time.Sleep(100 * time.Millisecond)
        }
        
        commonOp.SetState(utils.Success)
        commonOp.Stop()
        
        // 3. Apply environment-specific configs
        fmt.Println("👉 Applying environment-specific configurations...")
        envOp := manager.AddOperation("env", "Applying environment dotfiles...", utils.Bar)
        envOp.Start()
        
        for i := 0; i <= 100; i += 10 {
                envOp.UpdateProgress(i, 100)
                time.Sleep(150 * time.Millisecond)
        }
        
        envOp.SetState(utils.Success)
        envOp.Stop()
        
        // 4. Apply machine-specific configs
        fmt.Println("👉 Applying machine-specific configurations...")
        machineOp := manager.AddOperation("machine", "Applying machine dotfiles...", utils.Bar)
        machineOp.Start()
        
        hostname, _ := os.Hostname()
        fmt.Printf("   Detected hostname: %s\n", hostname)
        
        for i := 0; i <= 100; i += 15 {
                machineOp.UpdateProgress(i, 100)
                time.Sleep(200 * time.Millisecond)
        }
        
        machineOp.SetState(utils.Success)
        machineOp.Stop()
        
        // 5. Run installation scripts
        fmt.Println("👉 Running installation scripts...")
        scriptOp := manager.AddOperation("scripts", "Running setup scripts...", utils.Pulse)
        scriptOp.Start()
        
        time.Sleep(3 * time.Second)
        
        // Simulate some warnings during script execution
        scriptOp.SetState(utils.Warning)
        fmt.Println("   ⚠️ Some packages could not be installed")
        
        time.Sleep(1 * time.Second)
        scriptOp.SetState(utils.Success)
        scriptOp.Stop()
        
        // 6. Multiple concurrent operations
        fmt.Println("👉 Performing final verification...")
        
        // Add several concurrent operations with different styles
        op1 := manager.AddOperation("verify", "Verifying symlinks...", utils.Spinner)
        op2 := manager.AddOperation("check", "Checking file permissions...", utils.Dots)
        op3 := manager.AddOperation("analyze", "Analyzing configuration...", utils.Rainbow)
        
        op1.Start()
        op2.Start()
        op3.Start()
        
        // Let them run for a few seconds with different results
        time.Sleep(2 * time.Second)
        op1.SetState(utils.Success)
        op1.Stop()
        
        time.Sleep(1 * time.Second)
        op2.SetState(utils.Success)
        op2.Stop()
        
        time.Sleep(1 * time.Second)
        op3.SetState(utils.Success)
        op3.Stop()
        
        fmt.Println("✨ Bootstrap process completed successfully!")
}