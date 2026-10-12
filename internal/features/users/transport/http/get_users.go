package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
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
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get limit params:%w", err)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get offset params:%w", err)
	}

	return limit, offset, nil

}
