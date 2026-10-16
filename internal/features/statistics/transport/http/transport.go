package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_http_server "github.com/timururaimov/golang_todo_app/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsSrevice StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsSrevice: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Rout {
	return []core_http_server.Rout{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
