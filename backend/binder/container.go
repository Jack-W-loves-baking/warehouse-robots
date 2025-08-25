package binder

import (
	"warehouse-robots/backend/api/business/facades"
	"warehouse-robots/backend/api/dtos/sdk"
	"warehouse-robots/backend/api/v1/controller"
	"warehouse-robots/backend/api/v1/service"
	"warehouse-robots/backend/infra/sdkService"

	"warehouse-robots/backend/config"
)

// Container holds all dependencies
type Container struct {
	Config *config.Config

	// SDK Layer
	RobotSDKService sdk.Warehouse
	SDKFactory      *sdkService.RobotSDKFactory

	// Business Layer Facades
	CreateTaskFacade facades.ICreateTaskFacade

	// Service Layer
	CreateTaskService service.ICreateTaskService

	// Controller Layer
	CreateTaskController controller.ICreateTaskController
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) *Container {
	container := &Container{
		Config: cfg,
	}

	container.bindSDKLayer()
	container.bindBusinessLayer()
	container.bindServiceLayer()
	container.bindControllerLayer()

	return container
}

// bindSDKLayer sets up SDK services
func (c *Container) bindSDKLayer() {
	c.SDKFactory = sdkService.NewRobotSDKFactory(c.Config)
	c.RobotSDKService = c.SDKFactory.CreateRobotSDKService()
}

// bindBusinessLayer sets up business logic facades
func (c *Container) bindBusinessLayer() {
	c.CreateTaskFacade = facades.NewCreateTaskFacade(c.RobotSDKService)
}

// bindServiceLayer sets up service layer
func (c *Container) bindServiceLayer() {
	c.CreateTaskService = service.NewCreateTaskService(c.CreateTaskFacade)
}

// bindControllerLayer sets up controller layer
func (c *Container) bindControllerLayer() {
	c.CreateTaskController = controller.NewCreateTaskController(c.CreateTaskService)
}
