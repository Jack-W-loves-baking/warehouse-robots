package model

import (
	"time"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusFailed    TaskStatus = "FAILED"
	TaskStatusCancelled TaskStatus = "CANCELLED"
)

type Task struct {
	TaskID          string     `json:"task_id"`
	RobotID         string     `json:"robot_id"`
	Commands        string     `json:"commands"`
	Status          TaskStatus `json:"status"`
	CurrentPosition *Position  `json:"current_position,omitempty"`
	Error           string     `json:"error,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Position struct {
	X        uint `json:"x"`
	Y        uint `json:"y"`
	HasCrate bool `json:"has_crate"`
}
