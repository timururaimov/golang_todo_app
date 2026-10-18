package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/timururaimov/golang_todo_app/docs"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_middleware "github.com/timururaimov/golang_todo_app/internal/core/transport/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type HTTPSrever struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.MiddleWare
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middlWare ...core_http_middleware.MiddleWare,
) *HTTPSrever {
	return &HTTPSrever{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middlWare,
	}
}

func (s *HTTPSrever) RegisterAPIRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prfix := "/api/" + string(router.apiversion)

		s.mux.Handle(
			prfix+"/",
			http.StripPrefix(prfix, router.WithMiddlWare()),
		)
	}
}

func (s *HTTPSrever) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/docs.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)
	s.mux.HandleFunc(
		"/swagger/docs.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "aplication/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPSrever) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Address,
		Handler: mux,
	}
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		s.log.Warn("start HTTP server", zap.String("addr", s.config.Address))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutDown HTTP server...")

		shutDownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeOut,
		)
		defer cancel()
		if err := server.Shutdown(shutDownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		s.log.Warn("HTTP server stoped")

	}
	return nil
}
