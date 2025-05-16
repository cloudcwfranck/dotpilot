package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ProgressIndicator represents an animated progress indicator
type ProgressIndicator struct {
	message    string
	style      ProgressStyle
	running    bool
	stopCh     chan struct{}
	wg         sync.WaitGroup
	lastUpdate time.Time
	frames     []string
	frameIndex int
}

// ProgressStyle defines the type of animation to use
type ProgressStyle string

const (
	// Spinner displays a spinning animation
	Spinner ProgressStyle = "spinner"
	// Dots displays a series of advancing dots
	Dots ProgressStyle = "dots"
	// Bar displays a progress bar (requires percent)
	Bar ProgressStyle = "bar"
	// Bounce displays a bouncing animation
	Bounce ProgressStyle = "bounce"
)

// defaultSpinnerFrames contains the animation frames for the spinner
var defaultSpinnerFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

// dotsFrames contains the animation frames for the dots
var dotsFrames = []string{
	"⠿", "⠿⠿", "⠿⠿⠿", "⠿⠿⠿⠿", "⠿⠿⠿⠿⠿",
}

// bounceFrames contains the animation frames for the bounce animation
var bounceFrames = []string{
	"[    =    ]", "[   =     ]", "[  =      ]", "[ =       ]", "[=        ]",
	"[ =       ]", "[  =      ]", "[   =     ]", "[    =    ]", "[     =   ]",
	"[      =  ]", "[       = ]", "[        =]", "[       = ]", "[      =  ]",
	"[     =   ]",
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(message string, style ProgressStyle) *ProgressIndicator {
	p := &ProgressIndicator{
		message:    message,
		style:      style,
		stopCh:     make(chan struct{}),
		lastUpdate: time.Now(),
	}

	// Set frames based on style
	switch style {
	case Spinner:
		p.frames = defaultSpinnerFrames
	case Dots:
		p.frames = dotsFrames
	case Bounce:
		p.frames = bounceFrames
	case Bar:
		// Bar doesn't use frames, but needs a placeholder
		p.frames = []string{""}
	default:
		p.frames = defaultSpinnerFrames
		p.style = Spinner
	}

	return p
}

// Start begins the animation
func (p *ProgressIndicator) Start() {
	if p.running {
		return
	}

	p.running = true
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		
		for {
			select {
			case <-p.stopCh:
				return
			case <-time.After(100 * time.Millisecond):
				p.update()
			}
		}
	}()
}

// Stop ends the animation
func (p *ProgressIndicator) Stop() {
	if !p.running {
		return
	}

	close(p.stopCh)
	p.wg.Wait()
	p.running = false
	
	// Clear the animation line
	fmt.Printf("\r%s\r", strings.Repeat(" ", len(p.message)+20))
}

// SetMessage updates the message shown with the animation
func (p *ProgressIndicator) SetMessage(message string) {
	p.message = message
	p.update()
}

// update renders the current frame
func (p *ProgressIndicator) update() {
	frame := p.frames[p.frameIndex]
	p.frameIndex = (p.frameIndex + 1) % len(p.frames)

	var output string
	switch p.style {
	case Bar:
		output = fmt.Sprintf("\r%s [%s]", p.message, frame)
	case Spinner:
		output = fmt.Sprintf("\r%s %s", frame, p.message)
	case Dots:
		output = fmt.Sprintf("\r%s %s", p.message, frame)
	case Bounce:
		output = fmt.Sprintf("\r%s %s", p.message, frame)
	default:
		output = fmt.Sprintf("\r%s %s", frame, p.message)
	}

	// Print the animation frame
	fmt.Print(output)
}

// UpdateProgress updates the progress for bar-style indicators (0-100)
func (p *ProgressIndicator) UpdateProgress(percent int) {
	if p.style != Bar {
		return
	}

	// Ensure percent is in valid range
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	// Create a progress bar
	width := 20
	completed := width * percent / 100
	remaining := width - completed

	bar := strings.Repeat("=", completed)
	if remaining > 0 {
		bar += ">"
		remaining--
	}
	bar += strings.Repeat(" ", remaining)

	p.frames[0] = fmt.Sprintf("%s %3d%%", bar, percent)
	p.update()
}

// ProgressGroup manages multiple progress indicators
type ProgressGroup struct {
	indicators []*ProgressIndicator
}

// NewProgressGroup creates a new progress group
func NewProgressGroup() *ProgressGroup {
	return &ProgressGroup{
		indicators: make([]*ProgressIndicator, 0),
	}
}

// AddIndicator adds a new progress indicator to the group
func (pg *ProgressGroup) AddIndicator(message string, style ProgressStyle) *ProgressIndicator {
	p := NewProgressIndicator(message, style)
	pg.indicators = append(pg.indicators, p)
	return p
}

// StartAll starts all indicators in the group
func (pg *ProgressGroup) StartAll() {
	for _, indicator := range pg.indicators {
		indicator.Start()
	}
}

// StopAll stops all indicators in the group
func (pg *ProgressGroup) StopAll() {
	for _, indicator := range pg.indicators {
		indicator.Stop()
	}
}