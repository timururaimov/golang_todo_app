package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

type GetTaskResponse DTOTaskResponse

// GetTask godoc
// @Summary Получить задачу
// @Description Возвращает задачу по её ID.
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} GetTaskResponse "Задача найдена"
// @Failure 400 {object} core_http_respone.ErrorResponse "Некорректный ID задачи"
// @Failure 404 {object} core_http_respone.ErrorResponse "Задача не найдена"
// @Failure 500 {object} core_http_respone.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Faild to get task id ",
		)
		return
	}
	taskDomain, err := h.tasksService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Faild to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONRespone(response, http.StatusOK)

}
