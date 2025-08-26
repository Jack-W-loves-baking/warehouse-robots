package service

import (
	"warehouse-robots/backend/api/dtos"
)

// ICreateTaskService coordinates validation, enqueueing, persistence, and monitoring
// of robot tasks. Implementations are expected to:
//   - Resolve the target robot from the warehouse/SDK.
//   - Derive the starting position from the most recent terminal task and reject
//     creation if there is an active (pending/running) task.
//   - Validate the command sequence against warehouse bounds.
//   - Enqueue the commands to the SDK and persist a PENDING task record.
//   - Start background monitoring to keep task status/position up to date.
type ICreateTaskService interface {
	// CreateTask creates and enqueues a task for the given robot.
	//
	// Parameters:
	//   - robotID: identifier of the target robot.
	//   - req:     command payload to execute.
	//
	// Returns:
	//   - TaskInfo snapshot of the newly created task on success.
	//
	// Error Returns
	//	 - ErrRobotNotFound: robot not found
	//   - ErrTaskNotFound: task not found by the robot id
	//   - ErrTaskQueueFull: task is pending, but we want to queue another one.
	//	 - ErrBoundary: the robot will move out of the boundary if execute the given command.
	CreateTask(robotID string, req dtos.CreateTaskRequest) (*dtos.TaskInfo, error)
}
