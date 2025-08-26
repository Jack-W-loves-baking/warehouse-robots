package service

import (
	"warehouse-robots/backend/api/dtos"
)

// IRetrieveTaskService exposes read-only access to task state.
// Implementations should fetch a task snapshot by its ID and map the
// domain model into the public DTO returned to API callers.
type IRetrieveTaskService interface {
	// RetrieveTaskById returns a TaskInfo snapshot for the given task ID.
	// If the task does not exist, an error is returned.
	// Returns
	// - TaskInfo
	//
	// Error Returns:
	// - ErrTaskNotFound: task id not found in db.
	//
	RetrieveTaskById(taskID string) (*dtos.TaskInfo, error)
}
