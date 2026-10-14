package tasks_transport_http

import (
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

type DTOTaskResponse struct {
	Id           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserId int        `json:"author_user_id"`
}

func taskDTOFromDomain(task domain.Task) DTOTaskResponse {
	return DTOTaskResponse{
		Id:           task.Id,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserId: task.AuthorUserId,
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []DTOTaskResponse {
	tasksDTO := make([]DTOTaskResponse, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}
	return tasksDTO
}
