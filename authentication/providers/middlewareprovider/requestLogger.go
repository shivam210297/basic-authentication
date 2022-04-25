package middlewareprovider

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

const (
	HTTPRequestMethod     = "http.request.method"
	HTTPpRequestUserAgent = "http.request.user_agent"
	HTTPRequestScheme     = "http.request.scheme"
	HTTPRequestRemoteAddr = "http.request.remote_addr"
	HTTPRequestURI        = "http.request.uri"
	HTTPRequestDuration   = "http.request_duration"
	HTTPResponseStatus    = "http.response.status"
	HTTPResponseSize      = "http.response.size"
	Stack                 = "stack"
	Error                 = "error"
)

// StructuredLogger is a simple but powerful middleware logger backed by logrus.
type StructuredLogger struct{}

func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{}
}

// NewLogEntry creates a new middleware.logEntry based on the contents of the
// http.Request.
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fieldMap := map[string]interface{}{
		HTTPRequestScheme:     scheme,
		HTTPRequestMethod:     r.Method,
		HTTPRequestRemoteAddr: r.RemoteAddr,
		HTTPpRequestUserAgent: r.UserAgent(),
		HTTPRequestURI:        fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
	}

	return &StructuredLoggerEntry{
		Fields: fieldMap,
	}
}

// StructuredLoggerEntry is a single log entry in a StructuredLogger.
type StructuredLoggerEntry struct {
	Fields logrus.Fields
}

// Write is run at the end of a requests processing.
func (e *StructuredLoggerEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	fieldMap := e.Fields

	fieldMap[HTTPResponseSize] = bytes
	fieldMap[HTTPResponseStatus] = status
	fieldMap[HTTPRequestDuration] = float64(elapsed.Nanoseconds() / 1000000)

	logrus.WithFields(fieldMap).Infof("complete URI [%s]", fieldMap[HTTPRequestURI])
}

// Panic is run when a request panics during its processing.
func (e *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	fieldMap := e.Fields
	fieldMap[Stack] = string(stack)
	fieldMap[Error] = v
	logrus.WithFields(fieldMap).Error("request panicked!")
}
