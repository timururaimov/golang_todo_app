package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	authorUserId *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeOut())
	defer cancel()

	query := `
		SELECT id,version,title,description,completed,created_at,completed_at,author_user_id
		FROM todoapp.tasks
		 %s
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2;
		`

	args := []any{limit, offset}
	if authorUserId != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, authorUserId)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"select task : %w",
			err,
		)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
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
			return nil, fmt.Errorf(
				"scan task err : %w",
				err,
			)
		}
		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows :%w", err)
	}

	taskdomains := taskDomainsFromModel(taskModels)
	return taskdomains, nil
}
