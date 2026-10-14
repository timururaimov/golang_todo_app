package tasks_service

import (
	"context"
	"fmt"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

func (s *TasksService) CreatTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain : %w", err)
	}

	task, err := s.tasksRepository.CreatTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("creat task err : %w", err)
	}

	return task, nil

}
