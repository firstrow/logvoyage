// Backlog - holds clients messages that were not delivered to ES.
// Every 10sec backlog tries to resend messages to ES if backlog file isn't empty.
package main

import (
	"sync"
	"time"
)

const (
	backFileName = "back.log"
)

var (
	backFilePath   = ""
	backlogManager = &backlog{}
)

type backlog struct {
	sync.RWMutex
	lines []string
}

// Add new message to queue
func (b *backlog) AddMessage(m string) {
	b.Lock()
	b.lines = append(b.lines, m)
	b.Unlock()
}

func (b *backlog) Resend() {
	b.Lock()
	processing := b.lines
	b.lines = []string{}
	b.Unlock()
	for _, msg := range processing {
		processMessage(msg)
	}
}

func toBacklog(message string) {
	backlogManager.AddMessage(message)
}

func initBacklog() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		backlogManager.Resend()
	}
}
