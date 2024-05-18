// Package logger provides support for initializing the log system.
package logger

import (
	"context"

	"encore.dev/rlog"
)

// Logger represents a logger for logging information.
type Logger struct {
	handler rlog.Ctx
	events  Events
}

// New constructs a new log for application use.
func New(serviceName string) *Logger {
	return new(serviceName, Events{})
}

// NewWithEvents constructs a new log for application use with events.
func NewWithEvents(serviceName string, events Events) *Logger {
	return new(serviceName, events)
}

// Debug logs at LevelDebug with the given context.
func (log *Logger) Debug(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelDebug, 3, msg, args...)
}

// Debugc logs the information at the specified call stack position.
func (log *Logger) Debugc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelDebug, caller, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (log *Logger) Info(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelInfo, 3, msg, args...)
}

// Infoc logs the information at the specified call stack position.
func (log *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelInfo, caller, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func (log *Logger) Warn(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelWarn, 3, msg, args...)
}

// Warnc logs the information at the specified call stack position.
func (log *Logger) Warnc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelWarn, caller, msg, args...)
}

// Error logs at LevelError with the given context.
func (log *Logger) Error(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelError, 3, msg, args...)
}

// Errorc logs the information at the specified call stack position.
func (log *Logger) Errorc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelError, caller, msg, args...)
}

func (log *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	switch level {
	case LevelDebug:
		log.handler.Debug(msg, args...)
	case LevelInfo:
		log.handler.Info(msg, args...)
	case LevelWarn:
		log.handler.Warn(msg, args...)
	case LevelError:
		log.handler.Error(msg, args...)
	}
}

func new(serviceName string, events Events) *Logger {
	return &Logger{
		handler: rlog.With("service", serviceName),
		events:  events,
	}
}
