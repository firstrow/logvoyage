// Backlog - holds clients messages that were not delivered to ES.
// Every 10sec backlog tries to resend messages to ES if backlog file isn't empty.
package main

import (
	"log"
	"os"
	"sync"
	"time"
)

const (
	fallbackFileName = "back.log"
	checkFreq        = 10 // Number of seconds to run fallback check
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
		log.Println("Backlog memory is full.")
		saveMessageToFile(m)
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
	file, err := os.OpenFile(getFallbackFile(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Println("Error opening file", err)
	}
	defer file.Close()
	file.WriteString(m)
}

func getFallbackFile() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path + string(os.PathSeparator) + fallbackFileName
}

func initBacklog() {
	ticker := time.NewTicker(checkFreq * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		backlogManager.Resend()
	}
}
