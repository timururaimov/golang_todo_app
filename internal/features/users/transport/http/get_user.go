package users_transport_http

import (
	"net/http"

	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

type GetUserReasponse DTOUserResponse

// GetUser godoc
// @Summary Получить пользователя
// @Description Возвращает пользователя по его ID.
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя" example(5)
// @Success 200 {object} GetUserReasponse "Пользователь найден"
// @Failure 400 {object} core_http_respone.ErrorResponse "Некорректный ID пользователя"
// @Failure 404 {object} core_http_respone.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} core_http_respone.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)

	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"field to get userID ",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to get user",
		)
		return
	}
	response := GetUserReasponse(userDTOFromDomain(user))

	responseHandler.JSONRespone(response, http.StatusOK)

}
