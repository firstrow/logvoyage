// Help to track per-second user statistic.
// Collects logs count and size per second.
package main

import (
	"log"
	"sync"
	"time"
)

type perSecondStorage struct {
	sync.Mutex
	Logs  map[string]int
	Sizes map[string]int
}

var prs = perSecondStorage{Logs: make(map[string]int), Sizes: make(map[string]int)}

func initTimers() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		log.Println(prs)

		prs.Lock()
		prs.Logs = make(map[string]int)
		prs.Sizes = make(map[string]int)
		prs.Unlock()
	}
}

// indexName - user api key
func increaseCounter(indexName string, message string) {
	if _, ok := prs.Logs[indexName]; ok {
		prs.Logs[indexName] += 1
		prs.Sizes[indexName] += len(message)
	} else {
		prs.Logs[indexName] = 1
		prs.Sizes[indexName] = len(message)
	}
}
