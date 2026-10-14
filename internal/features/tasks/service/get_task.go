package tasks_service

import (
	"context"
	"fmt"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	taskId int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskId)
	if err != nil {
		return domain.Task{}, fmt.Errorf("faild to get task from repoeitory : %w", err)
	}

	return task, nil

}
