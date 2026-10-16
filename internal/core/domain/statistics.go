package domain

import "time"

type Statistics struct {
	TaskCreted                int
	TaskComleted              int
	TaskCompletedRate         *float64
	TaskAverageCompletionTime *time.Duration
}

func NewStatistics(
	taskCreted int,
	taskComleted int,
	taskCompletedRate *float64,
	taskAverageCompletedTime *time.Duration,
) Statistics {
	return Statistics{
		TaskCreted:                taskCreted,
		TaskComleted:              taskComleted,
		TaskCompletedRate:         taskCompletedRate,
		TaskAverageCompletionTime: taskAverageCompletedTime,
	}
}
