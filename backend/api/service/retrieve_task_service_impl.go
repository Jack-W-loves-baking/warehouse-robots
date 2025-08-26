package service

import (
	"warehouse-robots/backend/api/dao"
	"warehouse-robots/backend/api/dtos"
	"warehouse-robots/backend/api/model"
)

// RetrieveTaskServiceImpl is the default implementation of IRetrieveTaskService.
// It reads task records from the data in the memory and converts them to DTOs suitable
// for API responses.
type RetrieveTaskServiceImpl struct {
	repository dao.ITaskRepository
}

// NewRetrieveTaskService constructor
func NewRetrieveTaskService(repository dao.ITaskRepository) IRetrieveTaskService {
	return &RetrieveTaskServiceImpl{
		repository: repository,
	}
}

// RetrieveTaskById loads a task by ID, then maps the domain model to dtos.TaskInfo.
// It returns a not-found error if the task cannot be located.
func (s *RetrieveTaskServiceImpl) RetrieveTaskById(taskId string) (*dtos.TaskInfo, error) {
	task, err := s.repository.GetById(taskId)
	if err != nil {
		return nil, model.ErrTaskNotFound
	}

	taskInfo := &dtos.TaskInfo{
		TaskID:    task.TaskID,
		RobotID:   task.RobotID,
		Status:    mapToDtoStatus(task.Status),
		Commands:  task.Commands,
		Error:     task.Error,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	if task.CurrentPosition != nil {
		taskInfo.CurrentState = &dtos.RobotState{
			X:        task.CurrentPosition.X,
			Y:        task.CurrentPosition.Y,
			HasCrate: task.CurrentPosition.HasCrate,
		}
	}

	return taskInfo, nil
}

// mapToDtoStatus converts a domain TaskStatus into its DTO equivalent.
func mapToDtoStatus(s model.TaskStatus) dtos.TaskStatus {
	return dtos.TaskStatus(s)
}
