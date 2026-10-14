package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	taskId int,
) error {
	if err := s.tasksRepository.DeleteTask(ctx, taskId); err != nil {
		return fmt.Errorf("Fail to delete task : %w", err)
	}
	return nil
}
