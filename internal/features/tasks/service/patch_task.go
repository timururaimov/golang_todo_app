package tasks_service

import (
	"context"
	"fmt"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	taskId int,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskId)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task :%w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("patch task :%w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, taskId, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task %w", err)
	}
	return patchedTask, nil
}
