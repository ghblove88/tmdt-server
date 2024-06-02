package common

import (
	"sync"
	"time"
)

type StringQueue struct {
	queue     chan string
	popCancel chan struct{} // Channel to signal cancellation popCancel
	WriteLock sync.Mutex
}

func NewStringQueue(bufferSize int) *StringQueue {
	return &StringQueue{
		queue:     make(chan string, bufferSize),
		popCancel: make(chan struct{}),
		WriteLock: sync.Mutex{},
	}
}

func (sq *StringQueue) Push(item string) {
	sq.WriteLock.Lock()
	defer sq.WriteLock.Unlock()
	sq.queue <- item
}

func (sq *StringQueue) Pop() string {
	select {
	case item := <-sq.queue:
		return item
	case <-sq.popCancel:
		return ""
	}
}

func (sq *StringQueue) PopOrWait(timeout time.Duration) string {
	// 取消 pop 的等待
	close(sq.popCancel)
	defer func() { sq.popCancel = make(chan struct{}) }()

	select {
	case item := <-sq.queue:
		return item
	case <-time.After(timeout):
		return ""
	}
}

func (sq *StringQueue) Clear() {
	for len(sq.queue) > 0 {
		<-sq.queue // Drain the channel
	}
}
