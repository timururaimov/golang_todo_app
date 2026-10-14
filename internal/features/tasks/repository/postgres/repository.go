package tasks_postgres_repository

import core_postgers_pool "github.com/timururaimov/golang_todo_app/internal/core/repository/postgres/conn"

type TasksRepository struct {
	pool core_postgers_pool.Pool
}

func NewTasksRepository(
	pool core_postgers_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
