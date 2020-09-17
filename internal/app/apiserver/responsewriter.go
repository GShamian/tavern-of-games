package apiserver

import "net/http"

// responseWriter object. Stores http ResponseWriter
// and code variable.
type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader func. Wraps http WriteHeader.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
