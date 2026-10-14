package tasks_transport_http

import (
	"net/http"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

type CreatTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserId int     `json:"author_user_id" validate:"required"`
}

type CreatTaskResponse DTOTaskResponse

func (h *TasksHTTPHandler) CreatTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	var request CreatTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to decode and validate",
		)
		return
	}
	taskDomain := domainFromeDTO(request)
	taskDomain, err := h.tasksService.CreatTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to creat new Task",
		)
		return
	}
	response := CreatTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONRespone(response, http.StatusCreated)

}

func domainFromeDTO(dto CreatTaskRequest) domain.Task {
	return domain.NewTaskUnInizilized(
		dto.Title,
		dto.Description,
		dto.AuthorUserId,
	)

}
