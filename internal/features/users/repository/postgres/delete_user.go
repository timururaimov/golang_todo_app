package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeOut())
	defer cancel()

	query := `
	DELETE from todoapp.users
	WHERE id=$1
	`

	cmTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query : %w", err)
	}
	//Проверка на повеление запроса на таблицу
	if cmTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id : %d : %w", id, core_errors.ErrNotFound)
	}

	return nil

}
