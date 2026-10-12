package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
	"go.uber.org/zap"
)

const requestIdHEader = "X-Request-ID"

func RequestId() MiddleWare {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIdHEader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIdHEader, requestID)
			w.Header().Set(requestIdHEader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHEader)

			l := log.With(
				zap.String("request_id", requestId),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromeContext(ctx)
			rw := core_http_respone.NewResponseWriter(w)
			before := time.Now()
			log.Debug(
				">>>> incoming HTTP Request",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<<<Done HTTP Request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("letency", time.Since(before)),
			)
		})
	}
}

func Panic() MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromeContext(ctx)

			responseHandler := core_http_respone.NewHTTPResponseHandler(log, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)

				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
