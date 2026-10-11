package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
	core_http_utils "github.com/timururaimov/golang_todo_app/internal/core/transport/http/utils"
)

type GetUsersReasponse []DTOUserResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	reasponeHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetqueryParams(r)
	if err != nil {
		reasponeHandler.ErrorResponse(
			err,
			"faild to get queryPrams",
		)
		return
	}

	domainUsers, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		reasponeHandler.ErrorResponse(
			err,
			"faild to get users",
		)
		return
	}

	response := GetUsersReasponse(usersDTOFromeDomains(domainUsers))
	reasponeHandler.JSONRespone(response, http.StatusOK)

}

func getLimitOffsetqueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParams(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get limit params:%w", err)
	}

	offset, err := core_http_utils.GetIntQueryParams(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get offset params:%w", err)
	}

	return limit, offset, nil

}
