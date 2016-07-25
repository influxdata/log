package log_test

import (
	"reflect"
	"testing"

	"github.com/influxdata/log"
)

func TestLogger(t *testing.T) {
	h := &Handler{}
	logger := log.New(h)

	logger.Info("a")
	if h.LastEntry == nil {
		t.Error("unexpected entry: nil")
	} else {
		if h.LastEntry.Level != log.InfoLevel {
			t.Errorf("unexpected level: %v != %v", h.LastEntry.Level, log.InfoLevel)
		}
		if h.LastEntry.Message != "a" {
			t.Errorf("unexpected message: %v != %v", h.LastEntry.Message, "a")
		}
		if len(h.LastEntry.Fields) != 0 {
			t.Errorf("unexpected fields: %v != %v", h.LastEntry.Fields, nil)
		}
	}
	h.LastEntry = nil

	logger.Debug("b")
	if h.LastEntry != nil {
		t.Errorf("unexpected entry: %v", h.LastEntry)
	}
	h.LastEntry = nil

	logger.Warn("c")
	if h.LastEntry == nil {
		t.Error("unexpected entry: nil")
	} else {
		if h.LastEntry.Level != log.WarnLevel {
			t.Errorf("unexpected level: %v != %v", h.LastEntry.Level, log.WarnLevel)
		}
		if h.LastEntry.Message != "c" {
			t.Errorf("unexpected message: %v != %v", h.LastEntry.Message, "c")
		}
		if len(h.LastEntry.Fields) != 0 {
			t.Errorf("unexpected fields: %v != %v", h.LastEntry.Fields, nil)
		}
	}
}

func TestLogger_WithLevel(t *testing.T) {
	h := &Handler{}
	logger := log.New(h)
	l := logger.WithLevel(log.WarnLevel)

	l.Info("a")
	if h.LastEntry != nil {
		t.Errorf("unexpected entry: %v", h.LastEntry)
	}
	h.LastEntry = nil

	l.Warn("b")
	if h.LastEntry == nil {
		t.Error("unexpected entry: nil")
	} else {
		if h.LastEntry.Level != log.WarnLevel {
			t.Errorf("unexpected level: %v != %v", h.LastEntry.Level, log.WarnLevel)
		}
		if h.LastEntry.Message != "b" {
			t.Errorf("unexpected message: %v != %v", h.LastEntry.Message, "b")
		}
		if len(h.LastEntry.Fields) != 0 {
			t.Errorf("unexpected fields: %v != %v", h.LastEntry.Fields, nil)
		}
	}
}

func TestLogger_WithFields(t *testing.T) {
	h := &Handler{}
	logger := log.New(h)
	logger.Level = log.DebugLevel

	l := logger.WithFields(log.Fields{
		"foo": "bar",
	})

	l.Debug("a")
	if h.LastEntry == nil {
		t.Error("unexpected entry: nil")
	} else {
		if h.LastEntry.Level != log.DebugLevel {
			t.Errorf("unexpected level: %v != %v", h.LastEntry.Level, log.DebugLevel)
		}
		if h.LastEntry.Message != "a" {
			t.Errorf("unexpected message: %v != %v", h.LastEntry.Message, "a")
		}
		want := log.Fields{"foo": "bar"}
		if !reflect.DeepEqual(h.LastEntry.Fields, want) {
			t.Errorf("unexpected fields: %v != %v", h.LastEntry.Fields, want)
		}
	}
}

type Handler struct {
	Entries   []*log.Entry
	LastEntry *log.Entry
}

func (h *Handler) HandleLog(e *log.Entry) error {
	h.LastEntry = e
	h.Entries = append(h.Entries, e)
	return nil
}
