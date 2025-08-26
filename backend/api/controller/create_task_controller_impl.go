package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/dtos"
	"warehouse-robots/backend/api/helper"
	createTask "warehouse-robots/backend/api/service"
)

// CreateTaskControllerImpl handles HTTP requests that create tasks for robots.
// Responsibilities:
//   - Parse and validate incoming request payloads.
//   - Resolve the robot ID from the request path.
//   - Delegate business logic to ICreateTaskService.
//   - Map domain/service errors to HTTP responses with structured error bodies.
type CreateTaskControllerImpl struct {
	Service createTask.ICreateTaskService
	Helper  *helper.ControllerHelper
}

// NewCreateTaskController constructs a CreateTaskControllerImpl with the given service.
// The returned value satisfies ICreateTaskController.
func NewCreateTaskController(service createTask.ICreateTaskService) ICreateTaskController {
	return &CreateTaskControllerImpl{
		Service: service,
		Helper:  helper.NewControllerHelper(),
	}
}

// Handle for the endpoint
func (c *CreateTaskControllerImpl) Handle(w http.ResponseWriter, r *http.Request) {
	robotId := r.PathValue("robotId")

	var req dtos.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Helper.SendErrorResponse(w, http.StatusBadRequest,
			constant.ErrorCodeValidation, "Invalid JSON format", err.Error())
		return
	}

	if err := c.validateCommands(req.Commands); err != nil {
		c.Helper.SendErrorResponse(w, http.StatusBadRequest,
			constant.ErrorCodeValidation, err.Error(), "")
		return
	}

	// call service layer.
	taskInfo, err := c.Service.CreateTask(robotId, req)

	if err != nil {
		statusCode, errorCode := helper.MapErrorToHTTPStatus(err)
		c.Helper.SendErrorResponse(w, statusCode, errorCode, err.Error(), "")
		return
	}

	c.Helper.SendSuccessResponse(w, http.StatusCreated, taskInfo)
}

// validateCommands ensures the command string is non-empty and contains only
// the supported movement directives: N, S, E, W (case-insensitive; whitespace ignored).
func (c *CreateTaskControllerImpl) validateCommands(commands string) error {
	commands = strings.ToUpper(strings.ReplaceAll(commands, " ", ""))

	if len(commands) == 0 {
		log.Printf("commands cannot be empty")
		return fmt.Errorf("commands cannot be empty")
	}

	for i, cmd := range commands {
		if cmd != 'N' && cmd != 'S' && cmd != 'E' && cmd != 'W' {
			return fmt.Errorf("invalid command character '%c' at position %d. Only N, S, E, W are allowed", cmd, i+1)
		}
	}
	return nil
}
