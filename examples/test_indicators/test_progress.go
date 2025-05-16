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

// Progress Indicator test
func main() {
	fmt.Println("Testing animated progress indicators...")
	
	// Test spinner
	fmt.Println("\nTesting Spinner style:")
	spinner := NewSpinnerIndicator("Processing data")
	spinner.Start()
	time.Sleep(3 * time.Second)
	spinner.Stop()
	
	// Test progress bar
	fmt.Println("\nTesting Progress Bar style:")
	bar := NewBarIndicator("Downloading files")
	bar.Start()
	// Simulate progress
	for i := 0; i <= 100; i += 5 {
		bar.SetProgress(i)
		time.Sleep(100 * time.Millisecond)
	}
	bar.Stop()
	
	// Test bouncing indicator
	fmt.Println("\nTesting Bounce style:")
	bounce := NewBounceIndicator("Synchronizing")
	bounce.Start()
	time.Sleep(3 * time.Second)
	bounce.Stop()
	
	// Test dots indicator
	fmt.Println("\nTesting Dots style:")
	dots := NewDotsIndicator("Loading")
	dots.Start()
	time.Sleep(3 * time.Second)
	dots.Stop()
	
	// Test multiple indicators
	fmt.Println("\nTesting multiple concurrent indicators:")
	
	spinner2 := NewSpinnerIndicator("Operation 1")
	bar2 := NewBarIndicator("Operation 2")
	bounce2 := NewBounceIndicator("Operation 3")
	
	spinner2.Start()
	time.Sleep(500 * time.Millisecond)
	bar2.Start()
	time.Sleep(500 * time.Millisecond)
	bounce2.Start()
	
	// Simulate progress for the bar
	go func() {
		for i := 0; i <= 100; i += 2 {
			bar2.SetProgress(i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	time.Sleep(5 * time.Second)
	
	spinner2.Stop()
	bar2.Stop()
	bounce2.Stop()
	
	fmt.Println("\nProgress indicator tests completed!")
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
				fmt.Printf("\r%s\r", "                                        ")
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
				fmt.Printf("\r%s\r", "                                        ")
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
				fmt.Printf("\r%s\r", "                                        ")
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
				fmt.Printf("\r%s\r", "                                        ")
				d.done <- true
				return
			default:
				dots := ""
				for i := 0; i < count; i++ {
					dots += "."
				}
				fmt.Printf("\r%s%s%s", d.message, dots, "     ")
				
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