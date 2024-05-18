package logger

import (
	"context"
)

// Level represents a logging level.
type Level int

// A set of possible logging levels.
const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

// EventFn is a function to be executed when configured against a log level.
type EventFn func(ctx context.Context, msg string)

// Events contains an assignment of an event function to a log level.
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}
