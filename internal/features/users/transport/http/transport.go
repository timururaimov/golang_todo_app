package users_transport_http

import (
	"context"
	"net/http"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_http_server "github.com/timururaimov/golang_todo_app/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreatUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Rout {
	return []core_http_server.Rout{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreatUser,
		}, {
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		}, {
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		}, {
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		}, {
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
