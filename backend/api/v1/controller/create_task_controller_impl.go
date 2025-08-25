package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/dtos/api"
	"warehouse-robots/backend/api/v1/helper"
	createTask "warehouse-robots/backend/api/v1/service"
)

type Impl struct {
	Service createTask.ICreateTaskService
	Helper  *helper.ControllerHelper
}

// NewCreateTaskController constructor
func NewCreateTaskController(service createTask.ICreateTaskService) ICreateTaskController {
	return &Impl{
		Service: service,
		Helper:  helper.NewControllerHelper(),
	}
}

// Handle Implementation for the controller
func (c *Impl) Handle(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req api.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Helper.SendErrorResponse(w, http.StatusBadRequest,
			constant.ErrorCodeValidation, "Invalid JSON format", err.Error())
		return
	}

	// Basic validation
	if err := c.validateCommands(req.Commands); err != nil {
		c.Helper.SendErrorResponse(w, http.StatusBadRequest,
			constant.ErrorCodeValidation, err.Error(), "")
		return
	}

	// Call service layer
	taskInfo, err := c.Service.CreateTask(req)
	if err != nil {
		// Map error to appropriate HTTP status and error code
		statusCode, errorCode := c.mapErrorToHTTPStatus(err)
		c.Helper.SendErrorResponse(w, statusCode, errorCode, err.Error(), "")
		return
	}

	// Return successful response
	c.Helper.SendSuccessResponse(w, http.StatusCreated, taskInfo)
}

// validateCommands basic validation to check if commands is not empty and only has NWSE
func (c *Impl) validateCommands(commands string) error {
	// Clean commands (remove spaces, convert to uppercase)
	commands = strings.ToUpper(strings.ReplaceAll(commands, " ", ""))

	// Check if empty
	if len(commands) == 0 {
		return fmt.Errorf("commands cannot be empty")
	}

	// Validate each character
	for i, cmd := range commands {
		if cmd != 'N' && cmd != 'S' && cmd != 'E' && cmd != 'W' {
			return fmt.Errorf("invalid command character '%c' at position %d. Only N, S, E, W are allowed", cmd, i+1)
		}
	}

	return nil
}

// mapErrorToHTTPStatus maps service errors to HTTP status codes and error codes
func (c *Impl) mapErrorToHTTPStatus(err error) (int, string) {
	errorMsg := err.Error()

	switch {
	case strings.Contains(errorMsg, "invalid command"):
		return http.StatusBadRequest, constant.ErrorCodeValidation
	case strings.Contains(errorMsg, "out of bounds"):
		return http.StatusBadRequest, constant.ErrorCodeBoundary
	case strings.Contains(errorMsg, "no robots"):
		return http.StatusServiceUnavailable, constant.ErrorCodeNoRobots
	case strings.Contains(errorMsg, "busy"):
		return http.StatusConflict, constant.ErrorCodeRobotBusy
	case strings.Contains(errorMsg, "queue full"):
		return http.StatusTooManyRequests, constant.ErrorCodeTaskQueueFull
	default:
		return http.StatusInternalServerError, constant.ErrorCodeInternal
	}
}
