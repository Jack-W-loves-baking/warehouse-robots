package models

// Business Domain Objects.

// RobotPosition represents the internal position of a facades
type RobotPosition struct {
	X        uint
	Y        uint
	HasCrate bool
}

// RobotEntity represents a facades in the business domain
// For each facades, it should have an array of tasks.
type RobotEntity struct {
	ID       string
	Position RobotPosition
	Tasks    []TaskEntity
}

// TaskEntity represents a task in the business domain
type TaskEntity struct {
	ID        string
	Commands  string
	Status    TaskStatus
	CreatedAt int64
	Error     string
}

// TaskStatus represents task execution status in business domain
type TaskStatus string

const (
	TaskStatusQueued     TaskStatus = "QUEUED"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusFailed     TaskStatus = "FAILED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)
