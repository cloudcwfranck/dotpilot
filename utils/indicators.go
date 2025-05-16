package utils

import (
        "fmt"
        "io"
        "os"
        "strings"
        "sync"
        "time"
)

// ProgressStyle defines the visual style for an animated progress indicator
type ProgressStyle int

const (
        // Spinner is a rotating spinner animation
        Spinner ProgressStyle = iota
        // Bar is a progress bar that fills up
        Bar
        // Bounce is a bouncing animation
        Bounce
        // Dots is a series of animated dots
        Dots
        // Pulse is a pulsing animation that changes intensity
        Pulse
        // Rainbow is a color-cycling animation
        Rainbow
)

// ProgressIndicator represents an animated progress indicator
type ProgressIndicator struct {
        message     string
        style       ProgressStyle
        output      io.Writer
        done        chan bool
        stopOnce    sync.Once
        active      bool
        progressPct int // Only used for Bar style
        state       ProgressState // Current state (Normal, Success, Warning, Error, Info)
        mutex       sync.Mutex
}

// NewProgressIndicator creates a new progress indicator with the specified style
func NewProgressIndicator(message string, style ProgressStyle) *ProgressIndicator {
        return &ProgressIndicator{
                message: message,
                style:   style,
                output:  os.Stdout,
                done:    make(chan bool),
                active:  false,
                state:   Normal,
        }
}

// Start begins the progress animation
func (p *ProgressIndicator) Start() {
        p.mutex.Lock()
        if p.active {
                p.mutex.Unlock()
                return
        }
        p.active = true
        p.mutex.Unlock()

        go func() {
                switch p.style {
                case Spinner:
                        p.runSpinner()
                case Bar:
                        p.runBar()
                case Bounce:
                        p.runBounce()
                case Dots:
                        p.runDots()
                case Pulse:
                        p.runPulse()
                case Rainbow:
                        p.runRainbow()
                }
        }()
}

// Stop ends the progress animation
func (p *ProgressIndicator) Stop() {
        p.stopOnce.Do(func() {
                p.mutex.Lock()
                if !p.active {
                        p.mutex.Unlock()
                        return
                }
                p.active = false
                p.mutex.Unlock()
                p.done <- true
                // Clear the line after stopping
                fmt.Fprintf(p.output, "\r%s\r", strings.Repeat(" ", 80))
        })
}

// UpdateProgress updates the progress percentage (mainly for Bar style)
// This method can be called with either UpdateProgress(percent) or UpdateProgress(current, total)
func (p *ProgressIndicator) UpdateProgress(args ...int) {
        p.mutex.Lock()
        defer p.mutex.Unlock()
        
        var percent int
        
        if len(args) == 1 {
                // Called with just percentage
                percent = args[0]
        } else if len(args) >= 2 {
                // Called with current and total
                current := args[0]
                total := args[1]
                
                if total > 0 {
                        percent = (current * 100) / total
                }
        }
        
        if percent < 0 {
                percent = 0
        } else if percent > 100 {
                percent = 100
        }
        
        p.progressPct = percent
}

// SetMessage updates the message displayed with the progress indicator
func (p *ProgressIndicator) SetMessage(message string) {
        p.mutex.Lock()
        defer p.mutex.Unlock()
        p.message = message
}

// SetState updates the state of the progress indicator (Normal, Success, Warning, Error, Info)
func (p *ProgressIndicator) SetState(state ProgressState) {
        p.mutex.Lock()
        defer p.mutex.Unlock()
        p.state = state
}

// runSpinner displays a spinning animation
func (p *ProgressIndicator) runSpinner() {
        frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
        interval := 100 * time.Millisecond
        i := 0
        
        for {
                select {
                case <-p.done:
                        return
                default:
                        p.mutex.Lock()
                        if !p.active {
                                p.mutex.Unlock()
                                return
                        }
                        
                        frame := frames[i%len(frames)]
                        color := GetColorForState(p.state)
                        fmt.Fprintf(p.output, "\r%s%s%s %s", color, frame, Reset, p.message)
                        p.mutex.Unlock()
                        
                        time.Sleep(interval)
                        i++
                }
        }
}

// runBar displays a progress bar
func (p *ProgressIndicator) runBar() {
        barWidth := 20
        interval := 100 * time.Millisecond
        
        for {
                select {
                case <-p.done:
                        return
                default:
                        p.mutex.Lock()
                        if !p.active {
                                p.mutex.Unlock()
                                return
                        }
                        
                        progress := p.progressPct
                        filled := barWidth * progress / 100
                        unfilled := barWidth - filled
                        
                        color := GetColorForState(p.state)
                        bar := "[" + color + strings.Repeat("=", filled) + Reset + strings.Repeat(" ", unfilled) + "]"
                        
                        // Add colored percentage based on state
                        percentStr := fmt.Sprintf("%s%d%%%s", color, progress, Reset)
                        
                        fmt.Fprintf(p.output, "\r%s %s %s", bar, p.message, percentStr)
                        p.mutex.Unlock()
                        
                        time.Sleep(interval)
                }
        }
}

// runBounce displays a bouncing animation
func (p *ProgressIndicator) runBounce() {
        width := 20
        pos := 0
        direction := 1
        interval := 100 * time.Millisecond
        
        for {
                select {
                case <-p.done:
                        return
                default:
                        p.mutex.Lock()
                        if !p.active {
                                p.mutex.Unlock()
                                return
                        }
                        
                        color := GetColorForState(p.state)
                        line := strings.Repeat(" ", width)
                        runes := []rune(line)
                        
                        // Replace the position with a colored ball
                        runes[pos] = '⚫'
                        line = string(runes)
                        
                        fmt.Fprintf(p.output, "\r[%s%s%s] %s", color, line, Reset, p.message)
                        p.mutex.Unlock()
                        
                        if pos == width-1 {
                                direction = -1
                        } else if pos == 0 {
                                direction = 1
                        }
                        pos += direction
                        
                        time.Sleep(interval)
                }
        }
}

// runDots displays animated dots
func (p *ProgressIndicator) runDots() {
        max := 5
        i := 0
        interval := 300 * time.Millisecond
        
        for {
                select {
                case <-p.done:
                        return
                default:
                        p.mutex.Lock()
                        if !p.active {
                                p.mutex.Unlock()
                                return
                        }
                        
                        color := GetColorForState(p.state)
                        dots := strings.Repeat(".", i)
                        
                        // Colorize the dots
                        coloredDots := color + dots + Reset
                        
                        fmt.Fprintf(p.output, "\r%s%s%s", p.message, coloredDots, strings.Repeat(" ", max-i))
                        p.mutex.Unlock()
                        
                        i = (i + 1) % (max + 1)
                        
                        time.Sleep(interval)
                }
        }
}

// runPulse displays a pulsing animation that changes intensity
func (p *ProgressIndicator) runPulse() {
        symbols := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▂"}
        interval := 100 * time.Millisecond
        i := 0
        
        for {
                select {
                case <-p.done:
                        return
                default:
                        p.mutex.Lock()
                        if !p.active {
                                p.mutex.Unlock()
                                return
                        }
                        
                        color := GetColorForState(p.state)
                        symbol := symbols[i%len(symbols)]
                        
                        fmt.Fprintf(p.output, "\r%s%s%s %s", color, symbol, Reset, p.message)
                        p.mutex.Unlock()
                        
                        time.Sleep(interval)
                        i++
                }
        }
}

// runRainbow displays a rainbow animation with cycling colors
func (p *ProgressIndicator) runRainbow() {
        colors := []string{Red, Yellow, Green, Cyan, Blue, Purple}
        symbol := "◆"
        interval := 100 * time.Millisecond
        i := 0
        
        for {
                select {
                case <-p.done:
                        return
                default:
                        p.mutex.Lock()
                        if !p.active {
                                p.mutex.Unlock()
                                return
                        }
                        
                        // Cycle through colors regardless of state
                        color := colors[i%len(colors)]
                        
                        fmt.Fprintf(p.output, "\r%s%s%s %s", color, symbol, Reset, p.message)
                        p.mutex.Unlock()
                        
                        time.Sleep(interval)
                        i++
                }
        }
}