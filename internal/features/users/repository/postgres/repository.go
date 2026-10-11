package users_postgres_repository

import (
	core_postgers_pool "github.com/timururaimov/golang_todo_app/internal/core/repository/postgres/conn"
)

type UsersRepository struct {
	pool core_postgers_pool.Pool
}

func NewUsersRepository(
	pool core_postgers_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
