package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	taskId int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeOut())
	defer cancel()

	query := `
	SELECT id,version,title,description,completed,created_at,completed_at,author_user_id
	FROM todoapp.tasks
	WHERE id=$1;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		taskId,
	)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.Id,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"Task with id = '%d' : '%w'",
				taskId,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("Scan error: %w", err)

	}
	taskDomain := domain.NewTask(
		taskModel.Id,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserId,
	)

	return taskDomain, nil

}
