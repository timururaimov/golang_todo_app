package tasks_transport_http

import (
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

type DTOTaskResponse struct {
	Id           int        `json:"id" example:"5"`
	Version      int        `json:"version" example:"1"`
	Title        string     `json:"title" example:"Сделат домашку"`
	Description  *string    `json:"description" example:"Сделат домашку в 6:00"`
	Completed    bool       `json:"completed" example:"false"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-02-26T10:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"null"`
	AuthorUserId int        `json:"author_user_id" example:"2"`
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
