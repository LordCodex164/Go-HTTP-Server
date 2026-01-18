package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Logger handles structured logging
type Logger struct {
	logger *log.Logger
}

// New creates a new structured logger
func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	RequestID  string                 `json:"request_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Latency    string                 `json:"latency,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// Info logs an info message
func (l *Logger) Info(message string, extra map[string]interface{}) {
	l.log("INFO", message, "", extra)
}

// Error logs an error message
func (l *Logger) Error(message string, err error, extra map[string]interface{}) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "ERROR",
		Message:   message,
		Error:     errMsg,
		Extra:     extra,
	}
	l.write(entry)
}

// Request logs an HTTP request
func (l *Logger) Request(requestID, method, path string, statusCode int, latency time.Duration) {
	entry := LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Level:      "INFO",
		Message:    "HTTP Request",
		RequestID:  requestID,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		Latency:    latency.String(),
	}
	l.write(entry)
}

func (l *Logger) log(level, message, errMsg string, extra map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Error:     errMsg,
		Extra:     extra,
	}
	l.write(entry)
}

func (l *Logger) write(entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf("Failed to marshal log entry: %v", err)
		return
	}
	l.logger.Println(string(data))
}