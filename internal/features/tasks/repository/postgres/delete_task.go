package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(
	ctx context.Context,
	taskId int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeOut())
	defer cancel()

	query := `
		DELETE FROM todoapp.tasks
		WHERE id=$1;
	`
	cmTag, err := r.pool.Exec(ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("exec query : %w", err)
	}
	if cmTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id : %d : %w", taskId, core_errors.ErrNotFound)
	}
	return nil

}
