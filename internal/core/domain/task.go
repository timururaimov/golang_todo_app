package domain

import (
	"fmt"
	"time"

	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

type Task struct {
	Id           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserId int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserId int,
) Task {
	return Task{
		Id:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserId: authorUserId,
	}
}

func NewTaskUnInizilized(
	title string,
	description *string,
	authorUserId int,
) Task {
	return NewTask(
		UnInitializedId,
		UnInitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserId,
	)

}

func (t *Task) Validate() error {

	titleLenth := len([]rune(t.Title))
	if titleLenth < 1 || titleLenth > 100 {
		return fmt.Errorf(
			"invalid title lenth : %d:%w",
			titleLenth,
			core_errors.ErrInvalidArgument,
		)
	}
	if t.Description != nil {
		descriptionLenth := len([]rune(*t.Description))
		if descriptionLenth < 1 || descriptionLenth > 1000 {
			return fmt.Errorf(
				"invalid description lenth : %d:%w",
				descriptionLenth,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"CompletedAt cant be nil : %w",
				core_errors.ErrInvalidArgument,
			)
		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"Completed can't be befor vretedAt : %w",
				core_errors.ErrInvalidArgument,
			)
		}

	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"CompletedAt must be nil if Completed == false : %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	return nil

}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("Title cant be to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf("Completed cant be to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch : %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value
		if tmp.Completed {
			completedAT := time.Now()
			tmp.CompletedAt = &completedAT
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task : %w", err)
	}

	*t = tmp
	return nil
}
