package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
	core_http_types "github.com/timururaimov/golang_todo_app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Сделат Домашку"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Сделат Домащку до 6:00"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

type PatchTasResponse DTOTaskResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title cany be nul")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("Title msut be between 1 and 100")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("Description must be betweeen 1 and 1000 symvols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed cant be NULL")
		}
	}
	return nil
}

// PatchTask godoc
// @Summary Изменить задачу
// @Description Частично изменяет заголовок, описание и/или статус задачи.
// @Description #### Логика обновления полей (three-state logic)
// @Description 1. **Поле не передано** `description`— значение в БД не меняется.
// @Description 2. **Передано значение**`"description":"Утром пробежка в 6:00"` — поле получает новое значение.
// @Description 3. **Передан null**  `"description": null`— очищает поле в БД (Set to nul)
// @Description `title` не может быть null. `completed` должен быть true или false, но не null.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param request body PatchTaskRequest true "Поля задачи для изменения"
// @Success 200 {object} PatchTasResponse "Задача изменена"
// @Failure 400 {object} core_http_respone.ErrorResponse "Некорректный ID или данные запроса"
// @Failure 404 {object} core_http_respone.ErrorResponse "Задача не найдена"
// @Failure 409 {object} core_http_respone.ErrorResponse "Конфликт
// @Failure 500 {object} core_http_respone.ErrorResponse "Внутренняя ошибка сервера"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Faild to get TaskId",
		)
		return
	}
	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"Failed to Decode and Validatr",
		)
		return
	}
	taspPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskId, taspPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Failed to patch task",
		)
		return
	}

	response := PatchTasResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONRespone(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
