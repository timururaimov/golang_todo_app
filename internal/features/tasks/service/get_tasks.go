package tasks_service

import (
	"context"
	"fmt"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	authorUserId *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be not negative : %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be not negative : %w",
			core_errors.ErrInvalidArgument,
		)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, authorUserId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository:%w", err)
	}
	return tasks, nil

}
