package logger

import (
	apexlog "github.com/apex/log"
	"github.com/apex/log/handlers/discard"
)

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value interface{}) *apexlog.Entry {
	return apexlog.WithField(key, value)
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) *apexlog.Entry {
	return apexlog.WithError(err)
}

// Debug level message.
func Debug(msg string) {
	apexlog.Debug(msg)
}

// Info level message.
func Info(msg string) {
	apexlog.Info(msg)
}

// Warn level message.
func Warn(msg string) {
	apexlog.Warn(msg)
}

// Error level message.
func Error(msg string) {
	apexlog.Error(msg)
}

// Fatal level message, followed by an exit.
func Fatal(msg string) {
	apexlog.Fatal(msg)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	apexlog.Debugf(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	apexlog.Infof(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	apexlog.Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	apexlog.Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	apexlog.Fatalf(msg, v...)
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func Trace(msg string) *apexlog.Entry {
	return apexlog.Trace(msg)
}

// Discard all log entries, for benchmarking purposes
func Discard() {
	apexlog.SetHandler(discard.Default)
}
