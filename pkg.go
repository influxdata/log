package log

// Log is singleton commonly used as the root logger.
var Log = &Logger{
	Level: InfoLevel,
}

// SetHandler sets the handler. This is not thread-safe.
func SetHandler(h Handler) {
	Log.Handler = h
}

// SetLevel sets the log level. This is not thread-safe.
func SetLevel(level Level) {
	Log.Level = level
}

// WithLevel returns a new logger with `level` set.
func WithLevel(level Level) *Logger {
	return Log.WithLevel(level)
}

// WithFields returns a new logger with `fields` set.
func WithFields(fields Fields) *Logger {
	return Log.WithFields(fields)
}

// WithField returns a new logger with the `key` and `value` set.
func WithField(key string, value interface{}) *Logger {
	return Log.WithField(key, value)
}

// WithError returns a new logger with the "error" set to `err`.
func WithError(err error) *Logger {
	return Log.WithError(err)
}

// Debug level message.
func Debug(msg string) {
	Log.Debug(msg)
}

// Info level message.
func Info(msg string) {
	Log.Info(msg)
}

// Warn level message.
func Warn(msg string) {
	Log.Warn(msg)
}

// Error level message.
func Error(msg string) {
	Log.Error(msg)
}

// Fatal level message, followed by an exit.
func Fatal(msg string) {
	Log.Fatal(msg)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	Log.Debugf(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	Log.Infof(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	Log.Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	Log.Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	Log.Fatalf(msg, v...)
}
