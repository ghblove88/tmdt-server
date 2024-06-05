package common

import (
	"time"
)

// TimerTask is a structure that holds the timer's properties
type TimerTask struct {
	ticker       *time.Ticker
	done         chan bool
	initialDelay time.Duration
}

// NewTimer creates a new Timer instance
func NewTimer() *TimerTask {
	return &TimerTask{
		done: make(chan bool),
	}
}

// Start starts the timer that calls the specified function at fixed intervals with an initial delay
func (t *TimerTask) Start(initialDelay, interval time.Duration, task func()) {
	t.initialDelay = initialDelay
	time.AfterFunc(initialDelay, func() {
		t.ticker = time.NewTicker(interval)
		task() // Run the task immediately after the initial delay
		go func() {
			for {
				select {
				case <-t.done:
					return
				case <-t.ticker.C:
					task()
				}
			}
		}()
	})
}

// Stop stops the timer
func (t *TimerTask) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.done <- true
	}
}
