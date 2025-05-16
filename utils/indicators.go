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
func (p *ProgressIndicator) UpdateProgress(percent int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
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
			fmt.Fprintf(p.output, "\r%s %s", frame, p.message)
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
			
			bar := "[" + strings.Repeat("=", filled) + strings.Repeat(" ", unfilled) + "]"
			fmt.Fprintf(p.output, "\r%s %s %d%%", bar, p.message, progress)
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
			
			line := strings.Repeat(" ", width)
			runes := []rune(line)
			runes[pos] = '⚫'
			line = string(runes)
			
			fmt.Fprintf(p.output, "\r[%s] %s", line, p.message)
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
			
			dots := strings.Repeat(".", i)
			fmt.Fprintf(p.output, "\r%s%s%s", p.message, dots, strings.Repeat(" ", max-i))
			p.mutex.Unlock()
			
			i = (i + 1) % (max + 1)
			
			time.Sleep(interval)
		}
	}
}