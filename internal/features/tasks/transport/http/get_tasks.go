package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

type GetTasksResponse []DTOTaskResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	authorUserId, limit, offset, err := getAuthorUserIdLimitOffsetqueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to get queryParams",
		)
		return
	}

	domainTasks, err := h.tasksService.GetTasks(ctx, authorUserId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(domainTasks))
	responseHandler.JSONRespone(response, http.StatusOK)

}

func getAuthorUserIdLimitOffsetqueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		authorUserIdQueryParamKey = "author_user_id"
		limitQueryParamKey        = "limit"
		offsetQueryParamKey       = "offset"
	)

	authorUserId, err := core_http_request.GetIntQueryParams(r, authorUserIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get author_user_id query param : %w", err)
	}
	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get tasks limit params : %w", err)
	}
	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get tasks offsrt params : %w", err)
	}
	return authorUserId, limit, offset, nil
}
