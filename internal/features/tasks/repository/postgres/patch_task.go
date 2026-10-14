package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	taskId int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeOut())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET 
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4,
		version=version + 1
	WHERE id=$5 AND version=$6

	RETURNING
		id,
		version,
		title,
		description,
		completed,
		created_at,
		completed_at,
		author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.Id,
		task.Version,
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
				"task with id='%d', cocurently accsesed : %w",
				taskId,
				core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
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
