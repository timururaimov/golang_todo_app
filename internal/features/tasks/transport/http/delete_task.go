package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

// DeleteTask godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по её ID.
// @Tags tasks
// @Param id path int true "ID задачи"
// @Success 204 "Задача удалена"
// @Failure 400 {object} core_http_respone.ErrorResponse "Некорректный ID задачи"
// @Failure 404 {object} core_http_respone.ErrorResponse "Задача не найдена"
// @Failure 500 {object} core_http_respone.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Faild to get taskId",
		)
		return
	}
	if err := h.tasksService.DeleteTask(ctx, taskId); err != nil {
		responseHandler.ErrorResponse(
			err,
			"Faild to delete task",
		)
		return
	}
	responseHandler.NoContetnResponse()
}
