package dao

import (
	"warehouse-robots/backend/api/model"
)

type ITaskRepository interface {
	// Create stores a new task
	Create(task *model.Task) error

	// GetById Get retrieves a task by ID
	GetById(taskID string) (*model.Task, error)

	// GetByRobotId get tasks info by robot id.
	// one robot could have multiple tasks
	GetByRobotId(robotID string) ([]*model.Task, error)

	// Update updates an existing task
	Update(task *model.Task) error

	// UpdatePosition updates only the position and status
	UpdatePosition(taskID string, position *model.Position, status model.TaskStatus) error

	// UpdateStatus updates only the status and error if any
	UpdateStatus(taskID string, status model.TaskStatus, errorMsg string) error
}
