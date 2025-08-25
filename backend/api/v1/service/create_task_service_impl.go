package service

import (
	facade "warehouse-robots/backend/api/business/facades"
	"warehouse-robots/backend/api/business/models"
	dtos "warehouse-robots/backend/api/dtos/api"
)

type CreateTaskServiceImpl struct {
	createTaskFacade facade.ICreateTaskFacade
}

func NewCreateTaskService(createTaskFacade facade.ICreateTaskFacade) ICreateTaskService {
	return &CreateTaskServiceImpl{
		createTaskFacade: createTaskFacade, // Fixed: was using 'facade' instead of 'createTaskFacade'
	}
}

func (s *CreateTaskServiceImpl) CreateTask(req dtos.CreateTaskRequest) (*dtos.TaskInfo, error) {
	// Call the business facade to create the task
	taskEntity, err := s.createTaskFacade.CreateTask(req.Commands)
	if err != nil {
		return nil, err
	}

	// Map business model (TaskEntity) to API DTO (TaskInfo)
	taskInfo := s.mapTaskEntityToDTO(taskEntity)

	return taskInfo, nil
}

// mapTaskEntityToDTO maps business model TaskEntity to API DTO TaskInfo
func (s *CreateTaskServiceImpl) mapTaskEntityToDTO(taskEntity *models.TaskEntity) *dtos.TaskInfo {
	return &dtos.TaskInfo{
		ID:       taskEntity.ID,
		Commands: taskEntity.Commands,
		Status:   string(s.mapTaskStatus(taskEntity.Status)),
		Error:    taskEntity.Error,
	}
}

func (s *CreateTaskServiceImpl) mapTaskStatus(businessStatus models.TaskStatus) dtos.TaskStatus {
	switch businessStatus {
	case models.TaskStatusQueued:
		return dtos.TaskStatusQueued
	case models.TaskStatusInProgress:
		return dtos.TaskStatusInProgress
	case models.TaskStatusCompleted:
		return dtos.TaskStatusCompleted
	case models.TaskStatusFailed:
		return dtos.TaskStatusFailed
	case models.TaskStatusCancelled:
		return dtos.TaskStatusCancelled
	default:
		// Default to queued if unknown status
		return dtos.TaskStatusQueued
	}
}
