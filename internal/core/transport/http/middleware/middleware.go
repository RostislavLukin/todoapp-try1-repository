package middleware

import (
	"context"
	"net/http"
	"time"
	"todoapp-try1/internal/core/logger"
	"todoapp-try1/internal/core/transport/http/response"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	requestIDHeader = "X-Request-ID"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			rw.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(rw, r)
		})
	}
}

func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			reqL := log.With(
				zap.String("request_id", r.Header.Get(requestIDHeader)),
				zap.String("url", r.URL.String()),
				zap.String("http_method", r.Method),
			)

			ctx := context.WithValue(r.Context(), "log", reqL)

			next.ServeHTTP(rw, r.WithContext(ctx))

		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseHandler := response.NewHTTPResponseHandler(rw, log)

			defer func() {
				p := recover()
				if p != nil {
					responseHandler.PanicResponse(
						p,
						"during HTTP request handle got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseWriter := response.NewResponseWriter(rw)

			before := time.Now()
			log.Debug(
				"<<< incoming HTTP request",
				zap.Time("time", before),
			)

			next.ServeHTTP(responseWriter, r)

			log.Debug(
				">>> done HTTP request",
				zap.Duration("latency", time.Time.Sub(time.Now(), before)),
				zap.Int("status_code", responseWriter.GetStatusCodeOrPanic()),
			)
		})
	}
}
