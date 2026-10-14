package tasks_service

import (
	"context"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreatTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
	GetTasks(
		ctx context.Context,
		authorUserId *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)
	GetTask(
		ctx context.Context,
		taskId int,
	) (domain.Task, error)
	DeleteTask(
		ctx context.Context,
		taskId int,
	) error
	PatchTask(
		ctx context.Context,
		taskId int,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
