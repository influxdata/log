package log

import (
	"fmt"
	stdlog "log"
	"os"
	"time"
)

// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}

// Logger represents a logger with configurable Level and Handler.
type Logger struct {
	Level   Level
	Handler Handler
	parent  *Logger
	fields  Fields
}

// New returns a new logger at the default InfoLevel with the handler.
func New(h Handler) *Logger {
	return &Logger{
		Level:   InfoLevel,
		Handler: h,
	}
}

// WithLevel returns a new logger with `level` set.
func (l *Logger) WithLevel(level Level) *Logger {
	return &Logger{
		Level:  level,
		parent: l,
		fields: l.fields,
	}
}

// WithFields returns a new logger with `fields` set.
func (l *Logger) WithFields(fields Fields) *Logger {
	f := Fields{}
	for k, v := range l.fields {
		f[k] = v
	}
	for k, v := range fields {
		f[k] = v
	}
	return &Logger{
		parent: l,
		fields: f,
	}
}

// WithField returns a new logger with the `key` and `value` set.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.WithFields(Fields{key: value})
}

// WithError returns a new logger with the "error" set to `err`.
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err.Error())
}

// Log logs a message at the specified level.
func (l *Logger) Log(level Level, msg string) {
	if level < l.level() {
		return
	}

	handler := l.handler()
	if handler == nil {
		return
	}

	e := Entry{
		Fields:    l.fields,
		Level:     level,
		Message:   msg,
		Timestamp: time.Now(),
	}

	if err := handler.HandleLog(&e); err != nil {
		stdlog.Printf("error logging: %s", err)
	}
}

// Debug level message.
func (l *Logger) Debug(msg string) {
	l.Log(DebugLevel, msg)
}

// Info level message.
func (l *Logger) Info(msg string) {
	l.Log(InfoLevel, msg)
}

// Warn level message.
func (l *Logger) Warn(msg string) {
	l.Log(WarnLevel, msg)
}

// Error level message.
func (l *Logger) Error(msg string) {
	l.Log(ErrorLevel, msg)
}

// Fatal level message, followed by an exit.
func (l *Logger) Fatal(msg string) {
	l.Log(FatalLevel, msg)
	os.Exit(1)
}

// Logf logs a formatted message at the specified level.
func (l *Logger) Logf(level Level, msg string, v ...interface{}) {
	l.Log(level, fmt.Sprintf(msg, v...))
}

// Debugf level formatted message.
func (l *Logger) Debugf(msg string, v ...interface{}) {
	l.Logf(DebugLevel, msg, v...)
}

// Infof level formatted message.
func (l *Logger) Infof(msg string, v ...interface{}) {
	l.Logf(InfoLevel, msg, v...)
}

// Warnf level formatted message.
func (l *Logger) Warnf(msg string, v ...interface{}) {
	l.Logf(WarnLevel, msg, v...)
}

// Errorf level formatted message.
func (l *Logger) Errorf(msg string, v ...interface{}) {
	l.Logf(ErrorLevel, msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func (l *Logger) Fatalf(msg string, v ...interface{}) {
	l.Logf(FatalLevel, msg, v...)
	os.Exit(1)
}

// level returns the level of the Logger. If the level is
// not set, the parent's log level is used recursively.
func (l *Logger) level() Level {
	cur := l
	for cur != nil {
		if cur.Level != UnsetLevel {
			return cur.Level
		}
		cur = cur.parent
	}
	return InfoLevel
}

// handler returns the handler of the Logger. If the handler is
// not set, the parent's handler is used recursively.
func (l *Logger) handler() Handler {
	cur := l
	for cur != nil {
		if cur.Handler != nil {
			return cur.Handler
		}
		cur = cur.parent
	}
	return nil
}
