package statistics_postgres_repoeitory

import core_postgers_pool "github.com/timururaimov/golang_todo_app/internal/core/repository/postgres/conn"

type StatisticsRepository struct {
	pool core_postgers_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgers_pool.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
