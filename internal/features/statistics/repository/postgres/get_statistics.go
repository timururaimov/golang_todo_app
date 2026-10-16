package statistics_postgres_repoeitory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	ctx, cencel := context.WithTimeout(ctx, r.pool.OpTimeOut())
	defer cencel()
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id,version,title,description,completed,created_at,completed_at,author_user_id
	FROM todoapp.tasks
	`)
	args := []any{}
	conditions := []string{}

	if userId != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userId)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}
	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}
	queryBuilder.WriteString(" ORDER BY id ASC ")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks :%w", err)
	}
	defer rows.Close()
	var tasksModels []TaskModel

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
		tasksModels = append(tasksModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows :%w", err)
	}

	taskDomains := taskDomainsFromModel(tasksModels)
	return taskDomains, nil

}
