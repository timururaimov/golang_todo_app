package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/timururaimov/golang_todo_app/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type ApiVersionRouter struct {
	*http.ServeMux
	apiversion ApiVersion
	middleware []core_http_middleware.MiddleWare
}

func NewAPIVersionRouter(
	apiVersion ApiVersion,
	middleware ...core_http_middleware.MiddleWare,
) *ApiVersionRouter {
	return &ApiVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiversion: apiVersion,
		middleware: middleware,
	}
}

func (r *ApiVersionRouter) ReqisterRoutes(routes ...Rout) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddlWare())
	}
}

func (r *ApiVersionRouter) WithMiddlWare() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
