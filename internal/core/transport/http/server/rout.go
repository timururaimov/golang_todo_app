package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/timururaimov/golang_todo_app/internal/core/transport/http/middleware"
)

type Rout struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	MiddllWare []core_http_middleware.MiddleWare
}

func (r *Rout) WithMiddlWare() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.MiddllWare...,
	)
}

// func NEwRout(
// 	method string,
// 	path string,
// 	handler http.HandlerFunc,
// ) Rout {
// 	return Rout{
// 		Method:  method,
// 		Path:    path,
// 		Handler: handler,
// 	}
// }
