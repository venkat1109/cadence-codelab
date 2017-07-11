package eats

import (
	s "go.uber.org/cadence/.gen/go/shared"
)

type (
	// TaskStatus type for status value.
	TaskStatus string

	// Task models a job task.
	Task struct {
		ID       int64
		Name     string
		Status   TaskStatus
		SubTasks []*Task
	}

	// TaskGroupStatus type for status of a task group.
	TaskGroupStatus string

	// TaskGroup models a group of tasks.
	TaskGroup struct {
		ID      string
		RunID   string
		Status  TaskGroupStatus
		Tasks   []*Task
		TaskMap map[int64]*Task
		History *s.History
	}
)
