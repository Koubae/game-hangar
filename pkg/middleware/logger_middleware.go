package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	if lrw.wroteHeader {
		return
	}

	lrw.status = statusCode
	lrw.wroteHeader = true
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if !lrw.wroteHeader {
		lrw.WriteHeader(http.StatusOK)
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}

func AccessLogger(logger common.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{
				ResponseWriter: w,
			}

			next.ServeHTTP(lrw, r)

			if lrw.status == 0 {
				lrw.status = http.StatusOK
			}

			logger.Info(
				"http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", lrw.status),
				zap.Int("bytes", lrw.bytes),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Duration("duration", time.Since(start)),
			)
		},
	)
}

func RecoveryMiddleware(logger common.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			lrw, ok := w.(*loggingResponseWriter)
			if !ok {
				lrw = &loggingResponseWriter{ResponseWriter: w}
				w = lrw
			}

			defer func() {
				if rec := recover(); rec != nil {
					lrw.status = http.StatusInternalServerError
					http.Error(
						lrw,
						"Unexpected Server error",
						http.StatusInternalServerError,
					)

					logger.DPanic(
						fmt.Sprintf(
							"unhandled exception occurred: panic recovered in http handler, panic: %s",
							panicToString(rec),
						),
						zap.Any("panic", rec),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("query", r.URL.RawQuery),
						zap.String("remote_addr", r.RemoteAddr),
						zap.String("user_agent", r.UserAgent()),
						zap.String("stack", string(debug.Stack())),
					)
				}
			}()

			next.ServeHTTP(w, r)
		},
	)
}

func panicToString(r any) string {
	switch v := r.(type) {
	case nil:
		return ""
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
