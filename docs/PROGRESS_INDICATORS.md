# DotPilot Progress Indicators

This document provides detailed information about the animated progress indicators in DotPilot.

## Overview

DotPilot includes animated progress indicators that provide visual feedback during long-running operations like Git operations, file synchronization, encryption/decryption, and configuration application.

![DotPilot Progress Demo](demo.gif)

## Indicator Types

DotPilot implements four styles of animated progress indicators:

1. **Spinner** (`utils.Spinner`): A rotating animation that spins continuously
   - Best for: Operations with unknown completion time
   - Example: Commit operations, hooks execution

2. **Bar** (`utils.Bar`): A horizontal bar that fills from left to right
   - Best for: Operations with measurable progress
   - Example: Configuration application, file sync operations

3. **Bounce** (`utils.Bounce`): A dot that bounces back and forth 
   - Best for: Network operations
   - Example: Git pull/push operations

4. **Dots** (`utils.Dots`): Text followed by animated dots
   - Best for: Encryption/decryption operations
   - Example: SOPS file processing

## Usage in Code

### Basic Usage

```go
import "github.com/dotpilot/utils"

// Create a simple spinner
spinner := utils.NewOperation("commit", "Committing changes...", utils.Spinner)
spinner.Start()

// Perform the operation
// ... your code here ...

// Stop the spinner when done
spinner.Stop()
```

### Progress Bar with Updates

```go
// Create a progress bar
bar := utils.NewOperation("sync", "Syncing files...", utils.Bar)
bar.Start()

// Update progress as tasks complete
totalFiles := 10
for i := 1; i <= totalFiles; i++ {
    // ... process file ...
    
    // Update progress percentage
    percent := (i * 100) / totalFiles
    bar.Progress.UpdateProgress(percent)
}

// Stop the bar when done
bar.Stop()
```

### Multiple Concurrent Indicators

```go
// Create an operation manager
manager := utils.NewOperationManager()

// Add and start multiple operations
pullOp := manager.AddOperation("pull", "Pulling changes...", utils.Bounce)
configOp := manager.AddOperation("config", "Processing configs...", utils.Bar)

pullOp.Start()
configOp.Start()

// Simulate progress for the bar
configOp.SimulateProgress(5) // Simulate completion over 5 seconds

// ... your code here ...

// Stop operations when done
pullOp.Stop()
configOp.Stop()
```

## Disabling Progress Indicators

DotPilot commands provide a `--no-progress` flag to disable animated indicators when needed:

```bash
# Run without progress indicators (useful for scripts or CI)
dotpilot sync --no-progress
```

## Testing Progress Indicators

```bash
# Test all indicator types
go run test_progress.go

# Run a demo with simulated operations
go run demo.go

# Run the automated tests
go test ./utils -run TestProgressIndicatorTypes
```

## Implementation Details

The progress indicators are implemented in the `utils` package:
- `indicators.go`: Core progress indicator implementation
- `operations.go`: Operations management with progress tracking

The system is designed to be:
- Thread-safe for concurrent operations
- Cross-platform compatible
- Zero-dependency (uses only standard library)
- Customizable for different terminal environments