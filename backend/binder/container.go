package binder

import (
	controller "warehouse-robots/backend/api/controller"
	"warehouse-robots/backend/api/dao"
	"warehouse-robots/backend/api/manager"
	"warehouse-robots/backend/api/model"
	service "warehouse-robots/backend/api/service"
	"warehouse-robots/backend/infra/sdkService"

	"warehouse-robots/backend/config"
)

// Container holds all dependencies
type Container struct {
	Config *config.Config

	// SDK Layer
	RobotSDKService model.Warehouse
	SDKFactory      *sdkService.RobotSDKFactory

	// Repository Layer
	TaskRepository dao.ITaskRepository

	// Manager Layer
	TaskMonitor *manager.TaskMonitor

	// Service Layer
	CreateTaskService   service.ICreateTaskService
	RetrieveTaskService service.IRetrieveTaskService
	CancelTaskService   service.ICancelTaskService

	// Controller Layer
	CreateTaskController   controller.ICreateTaskController
	RetrieveTaskController controller.IRetrieveTaskController
	CancelTaskController   controller.ICancelTaskController
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	container := &Container{
		Config: cfg,
	}

	container.bindSDKLayer()
	container.bindDataLayer()
	container.bindManagerLayer()
	container.bindServiceLayer()
	container.bindControllerLayer()

	return container
}

// bindSDKLayer sets up SDK services
func (c *Container) bindSDKLayer() {
	c.SDKFactory = sdkService.NewRobotSDKFactory(c.Config)
	c.RobotSDKService = c.SDKFactory.CreateRobotSDKService()
}

// bindDataLayer sets up data access layer
func (c *Container) bindDataLayer() {
	// Create the shared repository instance
	c.TaskRepository = dao.NewInMemoryTaskRepository()
}

// bindManagerLayer sets up manager layer
func (c *Container) bindManagerLayer() {
	// TaskMonitor needs repository
	c.TaskMonitor = manager.NewTaskMonitor(c.TaskRepository)
}

// bindServiceLayer sets up service layer
func (c *Container) bindServiceLayer() {
	c.CreateTaskService = service.NewCreateTaskService(c.RobotSDKService,
		c.TaskRepository)
	c.RetrieveTaskService = service.NewRetrieveTaskService(c.TaskRepository)
	c.CancelTaskService = service.NewCancelTaskService(c.RobotSDKService,
		c.TaskRepository)
}

// bindControllerLayer sets up controller layer
func (c *Container) bindControllerLayer() {
	c.CreateTaskController = controller.NewCreateTaskController(c.CreateTaskService)
	c.RetrieveTaskController = controller.NewRetrieveTaskController(c.RetrieveTaskService)
	c.CancelTaskController = controller.NewCancelTaskController(c.CancelTaskService)
}
