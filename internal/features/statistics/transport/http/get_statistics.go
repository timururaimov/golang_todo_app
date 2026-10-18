package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

type GetStatisticsResponse struct {
	TaskCreted               int      `json:"task_creted"`
	TaskComleted             int      `json:"task_comleted"`
	TaskCompletedRate        *float64 `json:"task_completed_rate"`
	TaskAverageCompletedTime *string  `json:"task_average_completed_time"`
}

// GetStatistics godoc
// @Summary Получить статистику задач
// @Description Возвращает количество созданных и выполненных задач, процент выполнения и среднее время выполнения.
// @Description Можно отфильтровать статистику по автору и периоду создания задач. Для периода укажите обе даты: from включительно, to не включительно.
// @Tags statistics
// @Produce json
// @Param author_user_id query int false "ID пользователя-автора задач"
// @Param from query string false "Начало периода включительно, формат YYYY-MM-DD" example(2026-01-01)
// @Param to query string false "Конец периода не включительно, формат YYYY-MM-DD" example(2026-02-01)
// @Success 200 {object} GetStatisticsResponse "Статистика задач"
// @Failure 400 {object} core_http_respone.ErrorResponse "Некорректные параметры запроса или диапазон дат"
// @Failure 500 {object} core_http_respone.ErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)

	userId, from, to, err := getUserIdFromQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to get query params userId/from/to",
		)
		return
	}

	statistics, err := h.statisticsSrevice.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to get statistics",
		)
		return
	}

	response := toDTOFromDomain(statistics)
	responseHandler.JSONRespone(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TaskAverageCompletionTime != nil {
		duration := statistics.TaskAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TaskCreted:               statistics.TaskCreted,
		TaskComleted:             statistics.TaskComleted,
		TaskCompletedRate:        statistics.TaskCompletedRate,
		TaskAverageCompletedTime: avgTime,
	}
}

func getUserIdFromQueryParam(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIdQueryParamKey = "author_user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userId, err := core_http_request.GetIntQueryParams(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user_id query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get from query param %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get to query param %w", err)
	}
	return userId, from, to, nil
}
