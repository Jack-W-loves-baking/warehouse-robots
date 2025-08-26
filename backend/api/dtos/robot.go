package dtos

import (
	"time"
)

// Api model - Request and Response objects based on swagger API specification.
// Ideally this file should be auto generated.

// RobotState represents the state of a facades in API responses
type RobotState struct {
	X        uint `json:"x"`
	Y        uint `json:"y"`
	HasCrate bool `json:"has_crate"`
}

// CreateTaskRequest is the payload for creating a facades task
type CreateTaskRequest struct {
	Commands string
}

// TaskInfo contains information about a facades task (single facades system)
type TaskInfo struct {
	TaskID       string      `json:"task_id"`
	RobotID      string      `json:"robot_id"`
	Commands     string      `json:"commands"`
	Status       TaskStatus  `json:"status"`
	CurrentState *RobotState `json:"current_state,omitempty"`
	Error        string      `json:"error,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusFailed    TaskStatus = "FAILED"
	TaskStatusCancelled TaskStatus = "CANCELLED"
)
