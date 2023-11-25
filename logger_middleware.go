package songvote

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type LoggerMiddleware struct {
	handler http.Handler
}

func NewLoggerMiddleware(handlerToWrap http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{handlerToWrap}
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var body string
	if r.ContentLength > 0 {
		buf, _ := io.ReadAll(r.Body)
		body = string(buf)
		reader := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader
	}

	log.Printf("%s %s %s", r.Method, r.URL.Path, body)
	l.handler.ServeHTTP(w, r)
}
