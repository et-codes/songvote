// Package httplogger is a middleware that adds simple logging messages for
// HTTP requests received by the server.
package httplogger

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type HTTPLogger struct {
	handler http.Handler
}

// New() returns an http.Handler wrapped inside the logger middleware.
//
// Example use case:
//
//	mux := http.NewServeMux()
//	mux.Handle("/", http.HandleFunc(myHandler))
//	...
//	loggingMux := httplogger.New(mux)
//	log.Fatal(http.ListenAndServe(":5000", loggingMux))
func New(handlerToWrap http.Handler) *HTTPLogger {
	return &HTTPLogger{handlerToWrap}
}

func (l *HTTPLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body string

	// Check if the body contains anything.
	if r.ContentLength > 0 {
		// Read body contents.
		buf, _ := io.ReadAll(r.Body)
		body = string(buf)

		// Put body contents into a reader and add it back to the request.
		reader := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader
	}

	// Write log message.
	log.Printf("%s %s %s", r.Method, r.URL.Path, body)

	l.handler.ServeHTTP(w, r)
}
