package controller

import (
	"encoding/json"
	"fmt"
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

// Handle processes POST /robots/{robotId}/tasks requests.
//
// Request:
//   - Path:   robotId (string) resolved via r.PathValue("robotId").
//   - Body:   dtos.CreateTaskRequest (JSON).
//
// Responses:
//   - 201 Created: on successful creation, returns dtos.TaskInfo.
//   - 400 Bad Request: invalid JSON or invalid command sequence.
//   - 409 Conflict: the robot is working and task cannot be cancelled.
//   - 429 Too many requests: the queue has not Terminated, so we cannot queue a new task.
//   - 503 Service Unavailable: no robots available.
//   - 500 Internal Server Error: unexpected failures.
//
// Error bodies are standardized via ControllerHelper.
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

	taskInfo, err := c.Service.CreateTask(robotId, req)
	if err != nil {
		statusCode, errorCode := c.mapErrorToHTTPStatus(err)
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
		return fmt.Errorf("commands cannot be empty")
	}

	for i, cmd := range commands {
		if cmd != 'N' && cmd != 'S' && cmd != 'E' && cmd != 'W' {
			return fmt.Errorf("invalid command character '%c' at position %d. Only N, S, E, W are allowed", cmd, i+1)
		}
	}
	return nil
}

// mapErrorToHTTPStatus translates service-layer errors into HTTP status codes
// and system error codes used by API clients for programmatic handling.
func (c *CreateTaskControllerImpl) mapErrorToHTTPStatus(err error) (int, string) {
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
