package core_http_middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
	"go.uber.org/zap"
)

const requestIdHEader = "X-Request-ID"

const corsAllowedOriginsEnv = "CORS_ALLOWED_ORIGINS"

var corsAllowedMethods = map[string]struct{}{
	http.MethodGet: {}, http.MethodPost: {}, http.MethodPatch: {},
	http.MethodDelete: {}, http.MethodOptions: {},
}

var corsAllowedHeaders = map[string]struct{}{
	"content-type": {}, "authorization": {},
}

func CORS() MiddleWare {
	allowedOrigins := parseAllowedOrigins(os.Getenv(corsAllowedOriginsEnv))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			_, originAllowed := allowedOrigins[origin]
			if originAllowed {
				w.Header().Add("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}
			if r.Method == http.MethodOptions {
				if !originAllowed || !validPreflight(r) {
					http.Error(w, "CORS preflight request denied", http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func parseAllowedOrigins(value string) map[string]struct{} {
	origins := make(map[string]struct{})
	for _, origin := range strings.Split(value, ",") {
		origin = strings.TrimSpace(origin)
		if origin == "" || origin == "*" {
			continue
		}
		origins[origin] = struct{}{}
	}
	return origins
}

func validPreflight(r *http.Request) bool {
	method := r.Header.Get("Access-Control-Request-Method")
	if method == "" {
		return false
	}
	if _, ok := corsAllowedMethods[method]; !ok {
		return false
	}
	for _, header := range strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",") {
		header = strings.TrimSpace(strings.ToLower(header))
		if header == "" {
			continue
		}
		if _, ok := corsAllowedHeaders[header]; !ok {
			return false
		}
	}
	return true
}

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
