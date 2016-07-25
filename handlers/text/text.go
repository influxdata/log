// Package text implements a development-friendly textual handler.
package text

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/influxdata/log"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

const timeFormat = "2006/01/02 15:04:05 -0700"

// ColorMode determines whether colors are used.
type ColorMode int

const (
	// AutoColor automatically determines if a color is used from
	// checking the file descriptor.
	//
	// Color is only used when the underlying io.Writer is an *os.File.
	AutoColor ColorMode = iota

	// ForceColor forces colors to be used. The io.Writer can be any type.
	ForceColor

	// NoColor forces colors to not be used.
	NoColor
)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var Colors = [...]int{
	log.DebugLevel: gray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// field used for sorting.
type field struct {
	Name  string
	Value interface{}
}

// byName sorts projects by field name.
type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Handler implementation.
type Handler struct {
	mu         sync.Mutex
	isTerminal bool
	ColorMode  ColorMode
	Writer     io.Writer
}

// New creates a new handler with io.Writer.
func New(w io.Writer) *Handler {
	isTerminal := false
	if f, ok := w.(*os.File); ok {
		isTerminal = terminal.IsTerminal(int(f.Fd()))
	}
	return &Handler{
		Writer:     w,
		isTerminal: isTerminal,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	useColor := h.ColorMode == ForceColor || (h.ColorMode == AutoColor && h.isTerminal)

	var color int
	if useColor {
		color = Colors[e.Level]
	}
	level := Strings[e.Level]

	fields := make([]field, 0, len(e.Fields))
	for k, v := range e.Fields {
		if k == "service" {
			continue
		}
		fields = append(fields, field{k, v})
	}
	sort.Sort(byName(fields))

	var buf bytes.Buffer
	buf.WriteString(e.Timestamp.Format(timeFormat))
	if useColor {
		fmt.Fprintf(&buf, " |\033[%dm%5s\033[0m| ", color, level)
	} else {
		fmt.Fprintf(&buf, " |%5s| ", level)
	}
	buf.WriteString(e.Message)

	for _, f := range fields {
		if useColor {
			color := green
			if f.Name == "error" {
				color = red
			}
			fmt.Fprintf(&buf, " \033[%dm%6s\033[0m=%v", color, f.Name, f.Value)
		} else {
			fmt.Fprintf(&buf, " %6s=%v", f.Name, f.Value)
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintln(h.Writer, buf.String())
	return nil
}
