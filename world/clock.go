// world/clock.go
package world

import (
	"fmt"
	"sync"
	"time"
)

// Clock manages timing for the world
type Clock struct {
	interval time.Duration
	memory   *Memory2D
	running  bool
	stopChan chan struct{}
	mu       sync.Mutex
}

// NewClock creates a new clock with the specified interval in milliseconds
func NewClock(intervalMs int) *Clock {
	return &Clock{
		interval: time.Duration(intervalMs) * time.Millisecond,
		stopChan: make(chan struct{}),
	}
}

// Start begins the clock with the given memory
func (c *Clock) Start(memory *Memory2D) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("clock is already running")
	}

	c.memory = memory
	c.running = true
	c.stopChan = make(chan struct{})

	// Start the clock in a separate goroutine
	go c.run()

	return nil
}

// Stop halts the clock
func (c *Clock) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return fmt.Errorf("clock is not running")
	}

	c.running = false
	close(c.stopChan)
	return nil
}

// run is the main clock loop
func (c *Clock) run() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.memory != nil {
				// Process all heads
				if errors := c.memory.Bang(); len(errors) > 0 {
					// Log errors but continue running
					for _, err := range errors {
						fmt.Printf("Error during bang: %v\n", err)
					}
				}
			}
		case <-c.stopChan:
			return
		}
	}
}

// IsRunning returns whether the clock is currently running
func (c *Clock) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

// SetInterval updates the clock's interval
func (c *Clock) SetInterval(intervalMs int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("cannot change interval while clock is running")
	}

	c.interval = time.Duration(intervalMs) * time.Millisecond
	return nil
}
