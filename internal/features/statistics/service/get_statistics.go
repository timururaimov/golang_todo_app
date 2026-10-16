package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil || to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("to must be after from %w", core_errors.ErrInvalidArgument)
		}
	}
	tasks, err := s.statisticsRepository.GetTasks(
		ctx,
		userId,
		from,
		to,
	)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("failed to get task from repository: %w", err)
	}

	statistics := calcStatistics(tasks)
	return statistics, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}
	//Сколко созда
	tasksCreted := len(tasks)
	//сколко выполнено
	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++

		}
		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	//Пронцент выполненых задач
	tasksCopletedRate := float64(tasksCompleted) / float64(tasksCreted) * 100

	//среднее время выролнение
	var taskAverageComplettionTime *time.Duration
	if tasksCompleted > 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		taskAverageComplettionTime = &avg

	}
	return domain.NewStatistics(
		tasksCreted,
		tasksCompleted,
		&tasksCopletedRate,
		taskAverageComplettionTime,
	)
}
