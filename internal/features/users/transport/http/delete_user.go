package users_transport_http

import (
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID.
// @Tags users
// @Param id path int true "ID пользователя" example(5)
// @Success 204 "Пользователь удалён"
// @Failure 400 {object} core_http_respone.ErrorResponse "Некорректный ID пользователя"
// @Failure 404 {object} core_http_respone.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} core_http_respone.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to get usreId value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to dekete user",
		)
		return
	}

	responseHandler.NoContetnResponse()

}
