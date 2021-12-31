package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// defines Level type
type Level int8

// init constants for security level
const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

// human-readable error string
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Custom logger struct to hold log info
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// create new Logger instance
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// helper methods to write log entries
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// writes the log entry
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	// anonymous struct to hold log entry
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// include stack tract for ERROR and FATAL entries
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// holds the actual log entry text
	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message" + err.Error())
	}

	// locks mutex to prevent concurrent writing to output
	l.mu.Lock()
	defer l.mu.Unlock()

	// write the entry
	return l.out.Write(append(line, '\n'))
}
