package main

import (
	"log"
	"net/http"
	"time"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *StatusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srw := &StatusResponseWriter{ResponseWriter: w, status: http.StatusOK}
			defer func(ts time.Time) {
				log.Printf("[%s]: %s %s %d %s", name, r.Method, r.URL.Path, srw.status, time.Since(ts))
			}(time.Now())
			next.ServeHTTP(srw, r)
		})
	}
}
