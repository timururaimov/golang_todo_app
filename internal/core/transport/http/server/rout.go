package core_http_server

import "net/http"

type Rout struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NEwRout(
	method string,
	path string,
	handler http.HandlerFunc,
) Rout {
	return Rout{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
