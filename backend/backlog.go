// Backlog - holds clients messages that were not delivered to ES.
// Every 10sec backlog tries to resend messages to ES if backlog file isn't empty.
package main

import (
	"log"
	"sync"
	"time"
)

const (
	backFileName = "back.log"
)

var (
	backFilePath   = ""
	backlogManager = &backlog{}
	numStoreMsg    = 10000 // Hold in mem only N number of messages, otherwise write to file
)

type backlog struct {
	sync.RWMutex
	lines []string
	count int64
}

// Add new message to queue
func (b *backlog) AddMessage(m string) {
	b.Lock()
	if len(b.lines) <= numStoreMsg {
		b.lines = append(b.lines, m)
		b.count++
	} else {
		log.Println("Backlof is full.")
	}
	b.Unlock()
}

// Tries to resend messages
func (b *backlog) Resend() {
	b.Lock()
	processing := b.lines
	b.lines = []string{}
	b.count = 0
	b.Unlock()
	for _, msg := range processing {
		processMessage(msg)
	}
}

func toBacklog(m string) {
	backlogManager.AddMessage(m)
}

func saveMessageToFile(m string) {

}

func initBacklog() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		backlogManager.Resend()
	}
}
