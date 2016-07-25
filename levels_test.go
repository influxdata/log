package log_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/influxdata/log"
)

func TestLevel_String(t *testing.T) {
	for i, tt := range []struct {
		level log.Level
		want  string
	}{
		{
			level: log.DebugLevel,
			want:  "debug",
		},
		{
			level: log.InfoLevel,
			want:  "info",
		},
		{
			level: log.WarnLevel,
			want:  "warn",
		},
		{
			level: log.ErrorLevel,
			want:  "error",
		},
		{
			level: log.FatalLevel,
			want:  "fatal",
		},
	} {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("%d. unexpected string: got=%s want=%s", i, got, tt.want)
		}
	}
}

func TestLevel_MarshalJSON(t *testing.T) {
	v := struct {
		Level log.Level `json:"level"`
	}{Level: log.DebugLevel}

	got, err := json.Marshal(&v)
	if err != nil {
		t.Error(err)
	} else if want := []byte(`{"level":"debug"}`); !bytes.Equal(got, want) {
		t.Errorf("unexpected output: %s != %s", string(got), string(want))
	}
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	var v struct {
		Level log.Level `json:"level"`
	}

	if err := json.Unmarshal([]byte(`{"level":"debug"}`), &v); err != nil {
		t.Error(err)
	} else if v.Level != log.DebugLevel {
		t.Errorf("unexpected level: got=%s want=%s", v.Level, log.DebugLevel)
	}
}

func TestLevel_MarshalText(t *testing.T) {
	level := log.WarnLevel
	if got, err := level.MarshalText(); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if want := []byte("warn"); !bytes.Equal(got, want) {
		t.Errorf("unexpected text: got=%s want=%s", string(got), string(want))
	}
}

func TestLevel_UnmarshalText(t *testing.T) {
	var got log.Level
	if err := got.UnmarshalText([]byte("warn")); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if want := log.WarnLevel; got != want {
		t.Errorf("unexpected level: got=%s want=%s", got, want)
	}
}

func TestParseLevel(t *testing.T) {
	for i, tt := range []struct {
		s    string
		want log.Level
	}{
		{
			s:    "debug",
			want: log.DebugLevel,
		},
	} {
		level, err := log.ParseLevel(tt.s)
		if err != nil {
			t.Error(err)
		} else if level != tt.want {
			t.Errorf("%d. unexpected level: got=%s want=%s", i, level, tt.want)
		}
	}
}
