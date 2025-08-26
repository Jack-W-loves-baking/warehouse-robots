package controller

import (
	"net/http"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/helper"
	cancelTask "warehouse-robots/backend/api/service"
)

type CancelTaskControllerImpl struct {
	Service cancelTask.ICancelTaskService
	Helper  *helper.ControllerHelper
}

// NewCancelTaskController constructor
func NewCancelTaskController(service cancelTask.ICancelTaskService) ICancelTaskController {
	return &CancelTaskControllerImpl{
		Service: service,
		Helper:  helper.NewControllerHelper(),
	}
}

// Handle cancel controller handler
func (c *CancelTaskControllerImpl) Handle(w http.ResponseWriter, r *http.Request) {
	// Get task id from request url.
	taskId := r.PathValue("taskId")

	if taskId == "" {
		c.Helper.SendErrorResponse(w, http.StatusBadRequest,
			constant.ErrorCodeValidation, "Task ID is required", "")
		return
	}

	// call service layer
	err := c.Service.CancelTaskById(taskId)

	if err != nil {
		// Map error to appropriate HTTP status and error code
		statusCode, errorCode := helper.MapErrorToHTTPStatus(err)
		c.Helper.SendErrorResponse(w, statusCode, errorCode, err.Error(), "")
		return
	}

	// Return successful response (use 204 OK for Delete requests)
	c.Helper.SendNoContentResponse(w)
}
