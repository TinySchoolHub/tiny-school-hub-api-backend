package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// Helper to create a logger with buffer for testing
func newTestLogger(buf *bytes.Buffer, level zerolog.Level) *Logger {
	// Set the global level for this test
	zerolog.SetGlobalLevel(level)

	// Create logger with buffer
	zlog := zerolog.New(buf).With().Timestamp().Logger()

	return &Logger{
		logger: zlog,
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
	}{
		{"json format info level", "info", "json"},
		{"console format debug level", "debug", "console"},
		{"json format error level", "error", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.level, tt.format)
			if logger == nil {
				t.Error("New() returned nil")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"warning", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"fatal", zerolog.FatalLevel},
		{"unknown", zerolog.InfoLevel},
		{"", zerolog.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseLevel(tt.input)
			if got != tt.want {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.DebugLevel)

	logger.Debug("debug message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Debug() output = %q, want to contain 'debug message'", output)
	}
	if !strings.Contains(output, `"level":"debug"`) {
		t.Errorf("Debug() output = %q, want to contain level:debug", output)
	}
}

func TestLogger_Debugf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.DebugLevel)

	logger.Debugf("debug %s %d", "test", 123)

	output := buf.String()
	if !strings.Contains(output, "debug test 123") {
		t.Errorf("Debugf() output = %q, want to contain 'debug test 123'", output)
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.InfoLevel)

	logger.Info("info message")

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Errorf("Info() output = %q, want to contain 'info message'", output)
	}
	if !strings.Contains(output, `"level":"info"`) {
		t.Errorf("Info() output = %q, want to contain level:info", output)
	}
}

func TestLogger_Infof(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.InfoLevel)

	logger.Infof("info %s %d", "test", 456)

	output := buf.String()
	if !strings.Contains(output, "info test 456") {
		t.Errorf("Infof() output = %q, want to contain 'info test 456'", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.WarnLevel)

	logger.Warn("warning message")

	output := buf.String()
	if !strings.Contains(output, "warning message") {
		t.Errorf("Warn() output = %q, want to contain 'warning message'", output)
	}
	if !strings.Contains(output, `"level":"warn"`) {
		t.Errorf("Warn() output = %q, want to contain level:warn", output)
	}
}

func TestLogger_Warnf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.WarnLevel)

	logger.Warnf("warning %s", "test")

	output := buf.String()
	if !strings.Contains(output, "warning test") {
		t.Errorf("Warnf() output = %q, want to contain 'warning test'", output)
	}
}

func TestLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.ErrorLevel)

	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("Error() output = %q, want to contain 'error message'", output)
	}
	if !strings.Contains(output, `"level":"error"`) {
		t.Errorf("Error() output = %q, want to contain level:error", output)
	}
}

func TestLogger_Errorf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.ErrorLevel)

	logger.Errorf("error %d", 500)

	output := buf.String()
	if !strings.Contains(output, "error 500") {
		t.Errorf("Errorf() output = %q, want to contain 'error 500'", output)
	}
}

func TestLogger_WithField(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.InfoLevel)

	logger.WithField("user_id", "123").Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("WithField() output = %q, want to contain 'test message'", output)
	}
	if !strings.Contains(output, `"user_id":"123"`) {
		t.Errorf("WithField() output = %q, want to contain user_id:123", output)
	}
}

func TestLogger_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.InfoLevel)

	fields := map[string]interface{}{
		"user_id":    "123",
		"request_id": "abc-456",
		"count":      42,
	}

	logger.WithFields(fields).Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("WithFields() output = %q, want to contain 'test message'", output)
	}
	if !strings.Contains(output, `"user_id":"123"`) {
		t.Errorf("WithFields() output = %q, want to contain user_id", output)
	}
	if !strings.Contains(output, `"request_id":"abc-456"`) {
		t.Errorf("WithFields() output = %q, want to contain request_id", output)
	}
	if !strings.Contains(output, `"count":42`) {
		t.Errorf("WithFields() output = %q, want to contain count", output)
	}
}

func TestLogger_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.ErrorLevel)

	testErr := errors.New("test error")
	logger.WithError(testErr).Error("operation failed")

	output := buf.String()
	if !strings.Contains(output, "operation failed") {
		t.Errorf("WithError() output = %q, want to contain 'operation failed'", output)
	}
	if !strings.Contains(output, "test error") {
		t.Errorf("WithError() output = %q, want to contain 'test error'", output)
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	tests := []struct {
		name         string
		logLevel     zerolog.Level
		logFunc      func(*Logger)
		shouldAppear bool
	}{
		{
			name:     "debug logged when level is debug",
			logLevel: zerolog.DebugLevel,
			logFunc: func(l *Logger) {
				l.Debug("debug message")
			},
			shouldAppear: true,
		},
		{
			name:     "debug not logged when level is info",
			logLevel: zerolog.InfoLevel,
			logFunc: func(l *Logger) {
				l.Debug("debug message")
			},
			shouldAppear: false,
		},
		{
			name:     "info logged when level is info",
			logLevel: zerolog.InfoLevel,
			logFunc: func(l *Logger) {
				l.Info("info message")
			},
			shouldAppear: true,
		},
		{
			name:     "info not logged when level is error",
			logLevel: zerolog.ErrorLevel,
			logFunc: func(l *Logger) {
				l.Info("info message")
			},
			shouldAppear: false,
		},
		{
			name:     "error logged when level is info",
			logLevel: zerolog.InfoLevel,
			logFunc: func(l *Logger) {
				l.Error("error message")
			},
			shouldAppear: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := newTestLogger(buf, tt.logLevel)

			tt.logFunc(logger)

			output := buf.String()
			isEmpty := output == "" || output == "\n"

			if tt.shouldAppear && isEmpty {
				t.Errorf("Expected log output but got none")
			}
			if !tt.shouldAppear && !isEmpty {
				t.Errorf("Expected no log output but got: %q", output)
			}
		})
	}
}

func TestLogger_JSONOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.InfoLevel)

	logger.WithFields(map[string]interface{}{
		"user_id": "123",
		"action":  "login",
	}).Info("user logged in")

	output := buf.String()

	// Verify it's valid JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonMap); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// Check for expected fields
	if jsonMap["level"] != "info" {
		t.Errorf("JSON level = %v, want info", jsonMap["level"])
	}
	if jsonMap["message"] != "user logged in" {
		t.Errorf("JSON message = %v, want 'user logged in'", jsonMap["message"])
	}
	if jsonMap["user_id"] != "123" {
		t.Errorf("JSON user_id = %v, want 123", jsonMap["user_id"])
	}
}

func TestGlobal(t *testing.T) {
	logger := Global()
	if logger == nil {
		t.Error("Global() returned nil")
	}
}

func TestLogger_ChainedCalls(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf, zerolog.InfoLevel)

	// Test chaining WithField and WithError
	testErr := errors.New("chain error")
	logger.
		WithField("step", "1").
		WithField("component", "test").
		WithError(testErr).
		Info("chained log")

	output := buf.String()
	if !strings.Contains(output, "chained log") {
		t.Errorf("Chained output missing message: %q", output)
	}
	if !strings.Contains(output, `"step":"1"`) {
		t.Errorf("Chained output missing step field: %q", output)
	}
	if !strings.Contains(output, `"component":"test"`) {
		t.Errorf("Chained output missing component field: %q", output)
	}
	if !strings.Contains(output, "chain error") {
		t.Errorf("Chained output missing error: %q", output)
	}
}
