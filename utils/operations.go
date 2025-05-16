package utils

import (
        "time"
)

// Operation represents a long-running operation
type Operation struct {
        Name        string
        Description string
        Progress    *ProgressIndicator
        Total       int
        Current     int
        Done        bool
}

// NewOperation creates a new operation with progress tracking
func NewOperation(name, description string, style ProgressStyle) *Operation {
        op := &Operation{
                Name:        name,
                Description: description,
                Progress:    NewProgressIndicator(description, style),
                Total:       100,
                Current:     0,
                Done:        false,
        }
        return op
}

// Start begins the operation and progress tracking
func (op *Operation) Start() {
        op.Progress.Start()
}

// Stop ends the operation and progress tracking
func (op *Operation) Stop() {
        op.Progress.Stop()
        op.Done = true
}

// UpdateProgress updates the operation's progress
func (op *Operation) UpdateProgress(current, total int) {
        op.Current = current
        op.Total = total
        
        var percent int
        if total > 0 {
                percent = (current * 100) / total
        } else {
                percent = 0
        }
        
        op.Progress.UpdateProgress(percent)
}

// SetMessage updates the operation's description
func (op *Operation) SetMessage(message string) {
        op.Description = message
        op.Progress.SetMessage(message)
}

// SetState updates the state of the operation's progress indicator
func (op *Operation) SetState(state ProgressState) {
        op.Progress.SetState(state)
}

// SimulateProgress simulates progress for operations that don't report actual progress
func (op *Operation) SimulateProgress(seconds int) {
        go func() {
                startTime := time.Now()
                duration := time.Duration(seconds) * time.Second
                
                for time.Since(startTime) < duration && !op.Done {
                        elapsed := time.Since(startTime)
                        percent := int((elapsed.Seconds() / duration.Seconds()) * 100)
                        if percent > 100 {
                                percent = 100
                        }
                        
                        if op.Progress.style == Bar {
                                op.Progress.UpdateProgress(percent)
                        }
                        
                        time.Sleep(100 * time.Millisecond)
                }
                
                if !op.Done {
                        op.Progress.UpdateProgress(100)
                }
        }()
}

// OperationManager manages multiple operations
type OperationManager struct {
        Operations []*Operation
}

// NewOperationManager creates a new operation manager
func NewOperationManager() *OperationManager {
        return &OperationManager{
                Operations: make([]*Operation, 0),
        }
}

// AddOperation adds a new operation to the manager
func (om *OperationManager) AddOperation(name, description string, style ProgressStyle) *Operation {
        op := NewOperation(name, description, style)
        om.Operations = append(om.Operations, op)
        return op
}

// StartAll starts all operations in the manager
func (om *OperationManager) StartAll() {
        for _, op := range om.Operations {
                op.Start()
        }
}

// StopAll stops all operations in the manager
func (om *OperationManager) StopAll() {
        for _, op := range om.Operations {
                op.Stop()
        }
}

// FindOperation finds an operation by name
func (om *OperationManager) FindOperation(name string) *Operation {
        for _, op := range om.Operations {
                if op.Name == name {
                        return op
                }
        }
        return nil
}