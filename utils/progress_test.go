package utils

import (
	"bytes"
	"testing"
	"time"
)

// TestProgressIndicatorTypes verifies that all progress indicator types can be created and used
func TestProgressIndicatorTypes(t *testing.T) {
	// Create a custom writer to capture the output
	var buf bytes.Buffer
	
	// Test all progress indicator types
	styles := []struct {
		name  string
		style ProgressStyle
	}{
		{"Spinner", Spinner},
		{"Bar", Bar},
		{"Bounce", Bounce},
		{"Dots", Dots},
	}
	
	for _, style := range styles {
		t.Run(style.name, func(t *testing.T) {
			// Reset the buffer
			buf.Reset()
			
			// Create a progress indicator with the test style
			indicator := &ProgressIndicator{
				message: "Testing " + style.name,
				style:   style.style,
				output:  &buf,
				done:    make(chan bool),
				active:  false,
			}
			
			// Start the indicator
			indicator.Start()
			
			// For bar type, update progress
			if style.style == Bar {
				indicator.UpdateProgress(50)
			}
			
			// Let it run briefly
			time.Sleep(100 * time.Millisecond)
			
			// Stop the indicator
			indicator.Stop()
			
			// Verify that something was written to the buffer
			if buf.Len() == 0 {
				t.Errorf("%s indicator didn't produce any output", style.name)
			}
			
			t.Logf("%s indicator output: %q", style.name, buf.String())
		})
	}
}

// TestOperationManager verifies that the operation manager can handle multiple operations
func TestOperationManager(t *testing.T) {
	manager := NewOperationManager()
	
	// Add operations
	op1 := manager.AddOperation("test1", "Test Operation 1", Spinner)
	op2 := manager.AddOperation("test2", "Test Operation 2", Bar)
	
	// Start operations
	op1.Start()
	op2.Start()
	
	// Update progress for bar
	op2.SimulateProgress(1)
	
	// Let them run briefly
	time.Sleep(100 * time.Millisecond)
	
	// Stop operations
	op1.Stop()
	op2.Stop()
	
	t.Log("Operation manager test completed successfully")
}