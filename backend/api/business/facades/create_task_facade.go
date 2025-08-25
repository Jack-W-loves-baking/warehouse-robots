package facades

import (
	robotModels "warehouse-robots/backend/api/business/models"
)

// ICreateTaskFacade Facade interface to create a task and queue in the sdk.
type ICreateTaskFacade interface {
	CreateTask(commands string) (*robotModels.TaskEntity, error)
}
