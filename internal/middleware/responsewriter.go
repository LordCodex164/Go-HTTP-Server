package middleware

import (
	"log"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// NewResponseWriter creates a wrapped response writer
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default to 200
	}
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	log.Println(">>>code", code)
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// StatusCode returns the captured status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Write captures the response size
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Size returns the response size
func (rw *ResponseWriter) Size() int {
	return rw.size
}