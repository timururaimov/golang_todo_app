package statistics_postgres_repoeitory

import (
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

type TaskModel struct {
	Id           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserId int
}

func taskDomainsFromModel(tasks []TaskModel) []domain.Task {
	taskDomains := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		taskDomains[i] = domain.NewTask(
			task.Id,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.CreatedAt,
			task.CompletedAt,
			task.AuthorUserId,
		)
	}
	return taskDomains

}
