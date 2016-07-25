package log

import (
	"encoding/json"
	"errors"
)

// Level of severity.
type Level int

// Log levels.
const (
	UnsetLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = [...]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	FatalLevel: "fatal",
}

// String implements io.Stringer.
func (l Level) String() string {
	return levelNames[l]
}

// MarshalJSON returns the level string.
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// UnmarshalJSON returns the level from the string.
func (l *Level) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	level, err := ParseLevel(s)
	if err != nil {
		return err
	}
	*l = level
	return nil
}

// MarshalText marshals the level to a string.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText returns a level from a string.
func (l *Level) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	level, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*l = level
	return nil
}

// ParseLevel parses level string.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	default:
		return -1, errors.New("invalid level")
	}
}
