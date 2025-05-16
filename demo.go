package main

import (
	"fmt"
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

// Progress Indicator Demo
func main() {
	fmt.Println("DotPilot: Animated Progress Indicators Demo")
	fmt.Println("===========================================")
	
	// Demo 1: Simulating Git Operations
	fmt.Println("\n👉 Syncing dotfiles with remote repository...")
	
	// Step 1: Commit changes
	time.Sleep(500 * time.Millisecond)
	spinner := NewSpinnerIndicator("Auto-committing changes before sync")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.Stop()
	fmt.Println("✓ Changes committed")
	
	// Step 2: Pull changes
	time.Sleep(500 * time.Millisecond)
	bounce := NewBounceIndicator("Pulling changes from remote")
	bounce.Start()
	time.Sleep(2 * time.Second)
	bounce.Stop()
	fmt.Println("✓ Changes pulled from remote")
	
	// Step 3: Apply configurations
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\n👉 Applying configurations...")
	bar := NewBarIndicator("Applying configurations")
	bar.Start()
	
	// Simulate progress
	for i := 0; i <= 100; i += 5 {
		bar.SetProgress(i)
		time.Sleep(150 * time.Millisecond)
	}
	bar.Stop()
	fmt.Println("✓ Configurations applied")
	
	// Step 4: Push changes
	time.Sleep(500 * time.Millisecond)
	bounce2 := NewBounceIndicator("Pushing changes to remote")
	bounce2.Start()
	time.Sleep(2 * time.Second)
	bounce2.Stop()
	fmt.Println("✓ Changes pushed to remote")
	
	// Demo 2: Encrypting Secrets
	time.Sleep(1 * time.Second)
	fmt.Println("\n👉 Encrypting sensitive configuration files...")
	
	dots := NewDotsIndicator("Encrypting ~/.aws/credentials")
	dots.Start()
	time.Sleep(3 * time.Second)
	dots.Stop()
	fmt.Println("✓ Credentials encrypted successfully")
	
	// Demo 3: Multiple concurrent operations
	time.Sleep(1 * time.Second)
	fmt.Println("\n👉 Performing multiple concurrent operations...")
	
	op1 := NewSpinnerIndicator("Scanning for dotfiles")
	op2 := NewBarIndicator("Analyzing configurations")
	op3 := NewBounceIndicator("Checking remote status")
	
	op1.Start()
	time.Sleep(700 * time.Millisecond)
	op2.Start()
	time.Sleep(700 * time.Millisecond)
	op3.Start()
	
	// Simulate progress for the bar
	go func() {
		for i := 0; i <= 100; i += 4 {
			op2.SetProgress(i)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	
	time.Sleep(5 * time.Second)
	
	op1.Stop()
	fmt.Println("✓ Dotfiles scan complete")
	time.Sleep(300 * time.Millisecond)
	op2.Stop()
	fmt.Println("✓ Configuration analysis complete")
	time.Sleep(300 * time.Millisecond)
	op3.Stop()
	fmt.Println("✓ Remote status checked")
	
	fmt.Println("\n✨ DotPilot operations completed successfully!")
}

// -- Spinner Indicator Implementation --

type SpinnerIndicator struct {
	message string
	stop    chan bool
	done    chan bool
}

func NewSpinnerIndicator(message string) *SpinnerIndicator {
	return &SpinnerIndicator{
		message: message,
		stop:    make(chan bool),
		done:    make(chan bool),
	}
}

func (s *SpinnerIndicator) Start() {
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Printf("\r%s\r", "                                                                ")
				s.done <- true
				return
			default:
				frame := frames[i%len(frames)]
				fmt.Printf("\r%s %s", frame, s.message)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (s *SpinnerIndicator) Stop() {
	s.stop <- true
	<-s.done
}

// -- Bar Indicator Implementation --

type BarIndicator struct {
	message  string
	progress int
	stop     chan bool
	update   chan int
	done     chan bool
}

func NewBarIndicator(message string) *BarIndicator {
	return &BarIndicator{
		message:  message,
		progress: 0,
		stop:     make(chan bool),
		update:   make(chan int),
		done:     make(chan bool),
	}
}

func (b *BarIndicator) SetProgress(progress int) {
	b.update <- progress
}

func (b *BarIndicator) Start() {
	go func() {
		for {
			select {
			case <-b.stop:
				fmt.Printf("\r%s\r", "                                                                ")
				b.done <- true
				return
			case progress := <-b.update:
				b.progress = progress
				if b.progress > 100 {
					b.progress = 100
				}
				b.renderBar()
			default:
				b.renderBar()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (b *BarIndicator) renderBar() {
	width := 20
	filled := width * b.progress / 100
	padding := width - filled
	
	bar := "["
	for i := 0; i < filled; i++ {
		bar += "="
	}
	for i := 0; i < padding; i++ {
		bar += " "
	}
	bar += "]"
	
	fmt.Printf("\r%s %s %d%%", bar, b.message, b.progress)
}

func (b *BarIndicator) Stop() {
	b.stop <- true
	<-b.done
}

// -- Bounce Indicator Implementation --

type BounceIndicator struct {
	message string
	stop    chan bool
	done    chan bool
}

func NewBounceIndicator(message string) *BounceIndicator {
	return &BounceIndicator{
		message: message,
		stop:    make(chan bool),
		done:    make(chan bool),
	}
}

func (b *BounceIndicator) Start() {
	go func() {
		width := 20
		pos := 0
		direction := 1
		
		for {
			select {
			case <-b.stop:
				fmt.Printf("\r%s\r", "                                                                ")
				b.done <- true
				return
			default:
				bar := "["
				for i := 0; i < width; i++ {
					if i == pos {
						bar += "⚫"
					} else {
						bar += " "
					}
				}
				bar += "]"
				
				fmt.Printf("\r%s %s", bar, b.message)
				
				if pos == width-1 {
					direction = -1
				} else if pos == 0 {
					direction = 1
				}
				pos += direction
				
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (b *BounceIndicator) Stop() {
	b.stop <- true
	<-b.done
}

// -- Dots Indicator Implementation --

type DotsIndicator struct {
	message string
	stop    chan bool
	done    chan bool
}

func NewDotsIndicator(message string) *DotsIndicator {
	return &DotsIndicator{
		message: message,
		stop:    make(chan bool),
		done:    make(chan bool),
	}
}

func (d *DotsIndicator) Start() {
	go func() {
		count := 0
		max := 5
		
		for {
			select {
			case <-d.stop:
				fmt.Printf("\r%s\r", "                                                                ")
				d.done <- true
				return
			default:
				dots := ""
				for i := 0; i < count; i++ {
					dots += "."
				}
				padding := ""
				for i := 0; i < max-count; i++ {
					padding += " "
				}
				fmt.Printf("\r%s%s%s", d.message, dots, padding)
				
				count = (count + 1) % (max + 1)
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
}

func (d *DotsIndicator) Stop() {
	d.stop <- true
	<-d.done
}