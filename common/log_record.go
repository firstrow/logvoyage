package common

import (
	"time"
)

type LogRecord struct {
	Datetime time.Time `json:"datetime"`
	Message  string    `json:"message"`
}
