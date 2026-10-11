package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_middleware "github.com/timururaimov/golang_todo_app/internal/core/transport/http/middleware"

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

func (h *HTTPSrever) RegisterAPIRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prfix := "/api/" + string(router.apiversion)

		h.mux.Handle(
			prfix+"/",
			http.StripPrefix(prfix, router),
		)
	}
}

func (h *HTTPSrever) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Address,
		Handler: mux,
	}
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		h.log.Warn("start HTTP server", zap.String("addr", h.config.Address))

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
		h.log.Warn("shutDown HTTP server...")

		shutDownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeOut,
		)
		defer cancel()
		if err := server.Shutdown(shutDownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		h.log.Warn("HTTP server stoped")

	}
	return nil
}
