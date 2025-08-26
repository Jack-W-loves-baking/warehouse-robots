package controller

import (
	"net/http"
	"strings"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/helper"
	retrieveTask "warehouse-robots/backend/api/service"
)

type RetrieveTaskControllerImpl struct {
	Service retrieveTask.IRetrieveTaskService
	Helper  *helper.ControllerHelper
}

// NewRetrieveTaskController constructor
func NewRetrieveTaskController(service retrieveTask.IRetrieveTaskService) IRetrieveTaskController {
	return &RetrieveTaskControllerImpl{
		Service: service,
		Helper:  helper.NewControllerHelper(),
	}
}

func (c *RetrieveTaskControllerImpl) Handle(w http.ResponseWriter, r *http.Request) {
	// Get task id from request url.
	taskId := r.PathValue("taskId")

	if taskId == "" {
		c.Helper.SendErrorResponse(w, http.StatusBadRequest,
			constant.ErrorCodeValidation, "Task ID is required", "")
		return
	}

	taskInfo, err := c.Service.RetrieveTaskById(taskId)

	if err != nil {
		// Map error to appropriate HTTP status and error code
		statusCode, errorCode := c.mapErrorToHTTPStatus(err)
		c.Helper.SendErrorResponse(w, statusCode, errorCode, err.Error(), "")
		return
	}

	// Return successful response (use 200 OK for GET requests)
	c.Helper.SendSuccessResponse(w, http.StatusOK, taskInfo)
}

func (c *RetrieveTaskControllerImpl) mapErrorToHTTPStatus(err error) (int, string) {
	errorMsg := err.Error()

	switch {
	case strings.Contains(errorMsg, "not found"):
		return http.StatusNotFound, constant.ErrorCodeTaskNotFound
	default:
		return http.StatusInternalServerError, constant.ErrorCodeInternal
	}
}
