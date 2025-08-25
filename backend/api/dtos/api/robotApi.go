package api

// Api dtos - Request and Response objects based on swagger API specification.
// Ideally this file should be auto generated.

// RobotState represents the state of a facades in API responses
type RobotState struct {
	X        uint
	Y        uint
	HasCrate bool
}

// RobotInfo represents a facades in API responses
type RobotInfo struct {
	ID       string
	Position RobotState
}

// CreateTaskRequest is the payload for creating a facades task
type CreateTaskRequest struct {
	Commands string
}

// TaskInfo contains information about a facades task (single facades system)
type TaskInfo struct {
	ID       string
	Status   string
	Commands string
	Error    string
}

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Code    string
	Message string
	Details string
}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusQueued     TaskStatus = "QUEUED"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusFailed     TaskStatus = "FAILED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)
