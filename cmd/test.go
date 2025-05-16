package cmd

import (
        "fmt"
        "time"

        "github.com/dotpilot/utils"
        "github.com/spf13/cobra"
)

var (
        testDuration int
        testNoProgress bool
)

// testCmd represents the test command
var testCmd = &cobra.Command{
        Use:   "test",
        Short: "Test various features of dotpilot",
        Long: `Test command provides a way to test different features of dotpilot 
without affecting your actual dotfiles. Currently supports testing the
animated progress indicators with different styles.`,
        Run: func(cmd *cobra.Command, args []string) {
                if len(args) > 0 && args[0] == "progress" {
                        testProgressIndicators()
                        return
                }

                utils.Logger.Info().Msg("No specific test specified. Available tests: 'progress'")
        },
}

// testProgressCmd represents the test progress subcommand
var testProgressCmd = &cobra.Command{
        Use:   "progress",
        Short: "Test the animated progress indicators",
        Long: `Test the animated progress indicators with different styles 
(spinner, bar, bounce, dots) for a specified duration.`,
        Run: func(cmd *cobra.Command, args []string) {
                testProgressIndicators()
        },
}

func testProgressIndicators() {
        if testNoProgress {
                utils.Logger.Info().Msg("Progress indicators disabled. Use without --no-progress to see animations.")
                return
        }

        duration := time.Duration(testDuration) * time.Second
        utils.Logger.Info().Msgf("Testing progress indicators for %d seconds each", testDuration)
        
        // Create an operation manager to organize multiple indicators
        manager := utils.NewOperationManager()
        
        // Test spinner style
        utils.Logger.Info().Msg("Testing Spinner style...")
        spinnerOp := manager.AddOperation("spinner", "Testing Spinner style...", utils.Spinner)
        spinnerOp.Start()
        time.Sleep(duration)
        spinnerOp.Stop()
        
        // Test bar style
        utils.Logger.Info().Msg("Testing Bar style...")
        barOp := manager.AddOperation("bar", "Testing Bar style...", utils.Bar)
        barOp.Start()
        // Simulate progress for bar
        barOp.SimulateProgress(int(duration.Seconds()))
        barOp.Stop()
        
        // Test bounce style
        utils.Logger.Info().Msg("Testing Bounce style...")
        bounceOp := manager.AddOperation("bounce", "Testing Bounce style...", utils.Bounce)
        bounceOp.Start()
        time.Sleep(duration)
        bounceOp.Stop()
        
        // Test dots style
        utils.Logger.Info().Msg("Testing Dots style...")
        dotsOp := manager.AddOperation("dots", "Testing Dots style...", utils.Dots)
        dotsOp.Start()
        time.Sleep(duration)
        dotsOp.Stop()
        
        // Test multiple concurrent progress indicators
        utils.Logger.Info().Msg("Testing multiple concurrent indicators...")
        op1 := manager.AddOperation("multi1", "Testing concurrent operation 1...", utils.Spinner)
        op2 := manager.AddOperation("multi2", "Testing concurrent operation 2...", utils.Bar)
        op3 := manager.AddOperation("multi3", "Testing concurrent operation 3...", utils.Bounce)
        
        op1.Start()
        op2.Start()
        op2.SimulateProgress(int(duration.Seconds()))
        op3.Start()
        
        time.Sleep(duration)
        
        op1.Stop()
        op2.Stop()
        op3.Stop()
        
        fmt.Println("\nProgress indicator tests completed!")
}

func init() {
        rootCmd.AddCommand(testCmd)
        testCmd.AddCommand(testProgressCmd)
        
        testProgressCmd.Flags().IntVar(&testDuration, "duration", 3, "Duration in seconds to display each progress indicator")
        testProgressCmd.Flags().BoolVar(&testNoProgress, "no-progress", false, "Disable progress indicators")
}