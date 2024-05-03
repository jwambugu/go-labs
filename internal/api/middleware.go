package api

import (
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type LoggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *LoggingResponseWriter) Write(bb []byte) (int, error) {
	wb, err := lrw.w.Write(bb)
	return wb, err
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.w.WriteHeader(statusCode)
	lrw.statusCode = statusCode
}

func (a *Api) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			start     = time.Now()
			loggingRW = &LoggingResponseWriter{w: w}
		)

		next.ServeHTTP(loggingRW, r)

		duration := time.Since(start)

		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
				remoteAddr = "unknown"
			} else {
				remoteAddr = ip
			}
		}

		log.Printf(
			"%s | %s [%d] %s %s %s",
			loggingRW.Header().Get("X-Labs-Request-ID"),
			remoteAddr,
			loggingRW.statusCode,
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}

func (a *Api) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.Must(uuid.NewV6()).String()
		w.Header().Set("X-Labs-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}

func (a *Api) chainMiddleware(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}

		return next
	}
}
