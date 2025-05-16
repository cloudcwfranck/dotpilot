# DotPilot Progress Indicators

This document provides an overview of the progress indicator system in DotPilot, explaining how to use different indicator styles and states to provide visual feedback for long-running operations.

## Overview

DotPilot features a flexible progress indicator system that provides real-time visual feedback for operations like Git synchronization, file operations, and configuration management. These indicators help users understand what's happening during potentially long-running tasks.

## Indicator Styles

DotPilot supports the following progress indicator styles:

### 1. Spinner (`Spinner`)

A rotating spinner animation best used for operations with unknown completion time.

```
⠋ Committing changes...
⠙ Committing changes...
⠹ Committing changes...
```

**Best for:** Operations where progress percentage can't be determined, such as Git operations or network tasks.

### 2. Progress Bar (`Bar`)

A horizontal bar that fills up based on completion percentage.

```
[==========          ] Syncing files 50%
[====================] Syncing files 100%
```

**Best for:** Operations with measurable progress like file transfers or batch processing.

### 3. Bounce (`Bounce`)

A bouncing ball animation that moves back and forth.

```
[⚫                   ] Pulling changes from remote
[   ⚫                ] Pulling changes from remote
[                   ⚫] Pulling changes from remote
```

**Best for:** Network operations or tasks with indeterminate length.

### 4. Dots (`Dots`)

A simple animation that adds dots progressively.

```
Loading
Loading.
Loading..
Loading...
```

**Best for:** Simple waiting indicators or less visually intensive feedback.

### 5. Pulse (`Pulse`)

A pulsing animation that changes intensity.

```
▁ Encrypting files
▄ Encrypting files
█ Encrypting files
▄ Encrypting files
```

**Best for:** Security operations or status monitoring.

### 6. Rainbow (`Rainbow`)

A color-cycling animation for visually distinctive progress indication.

```
◆ Processing (cycles through colors)
```

**Best for:** Drawing attention to critical operations or creating visual distinction between multiple concurrent tasks.

## Indicator States

Progress indicators can be in different states that affect their color:

1. **Normal** - Default state (default terminal color)
2. **Success** - Indicates successful operation (green)
3. **Warning** - Indicates a warning condition (yellow)
4. **Error** - Indicates an error condition (red)
5. **Info** - Indicates informational status (cyan)

## Using Progress Indicators in DotPilot Commands

### Basic Usage

```go
// Create a new spinner indicator
spinner := utils.NewProgressIndicator("Loading configurations", utils.Spinner)

// Start the animation
spinner.Start()

// Do some work...
// ...

// Update the state if needed
spinner.SetState(utils.Success)

// Stop the animation when done
spinner.Stop()
```

### For Measurable Progress

```go
// Create a progress bar
bar := utils.NewProgressIndicator("Downloading dotfiles", utils.Bar)
bar.Start()

// Update progress as work completes
for i := 0; i <= 100; i += 10 {
    bar.UpdateProgress(i)
    // Do 10% of the work...
    time.Sleep(500 * time.Millisecond)
}

bar.SetState(utils.Success)
bar.Stop()
```

### Multiple Concurrent Indicators

DotPilot supports running multiple progress indicators concurrently:

```go
// Create operation manager
manager := utils.NewOperationManager()

// Add multiple operations
spinner := manager.AddOperation("fetch", "Fetching remote changes", utils.Spinner)
bar := manager.AddOperation("apply", "Applying configurations", utils.Bar)

// Start operations
spinner.Start()
// ... do some work
bar.Start()
// ... do more work

// Update progress as needed
bar.UpdateProgress(50)

// Stop when done
spinner.Stop()
bar.Stop()
```

## CLI Arguments

DotPilot commands support a global `--no-progress` flag that disables animated progress indicators, which is useful in CI/CD environments or when redirecting output to a file.

```bash
dotpilot sync --no-progress
```

## Customization

Progress indicators are fully customizable. You can:

1. Change the message during operation
2. Update the state to reflect current status
3. Choose the most appropriate style for each operation
4. Use an operation manager for coordinating multiple indicators

## Best Practices

1. Use Spinner for operations with unknown duration
2. Use Bar for operations where progress can be measured
3. Use Bounce for network operations
4. Use Dots for simple waiting indicators
5. Use Pulse for security operations
6. Use Rainbow for critical or attention-requiring operations
7. Update indicator state to reflect operation status (success, warning, error)
8. Always stop indicators when operations complete
9. Use the `--no-progress` flag in scripts and CI/CD environments