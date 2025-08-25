package service

import (
	dtos "warehouse-robots/backend/api/dtos/api"
)

// ICreateTaskService defines the interface for create task service
type ICreateTaskService interface {
	CreateTask(req dtos.CreateTaskRequest) (*dtos.TaskInfo, error)
}
