package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/model"

	"warehouse-robots/backend/api/dtos"
)

// ControllerHelper provides common functionality for all controllers
type ControllerHelper struct{}

// NewControllerHelper creates a new instance of ControllerHelper
func NewControllerHelper() *ControllerHelper {
	return &ControllerHelper{}
}

// SendErrorResponse sends a standardized error response
func (h *ControllerHelper) SendErrorResponse(w http.ResponseWriter, statusCode int, code, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dtos.ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}

// SendSuccessResponse sends a successful response with data
func (h *ControllerHelper) SendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendNoContentResponse sends a 204 No Content response
func (h *ControllerHelper) SendNoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// MapErrorToHTTPStatus maps typed service errors to (HTTP status, API error code).
// It relies on errors.Is(err, model.ErrX) checks (sentinel errors).
func MapErrorToHTTPStatus(err error) (int, string) {
	switch {
	// 400
	case errors.Is(err, model.ErrValidation):
		return http.StatusBadRequest, constant.ErrorCodeValidation
	case errors.Is(err, model.ErrRobotIDInvalid):
		return http.StatusBadRequest, constant.ErrorCodeRobotIdInvalid
	case errors.Is(err, model.ErrBoundary):
		return http.StatusBadRequest, constant.ErrorCodeBoundary

	// 404
	case errors.Is(err, model.ErrTaskNotFound):
		return http.StatusNotFound, constant.ErrorCodeTaskNotFound
	case errors.Is(err, model.ErrRobotNotFound):
		return http.StatusNotFound, constant.ErrorCodeRobotNotFound

	// 409
	case errors.Is(err, model.ErrRobotBusy):
		return http.StatusConflict, constant.ErrorCodeRobotBusy
	case errors.Is(err, model.ErrTaskProcessed): // already terminal
		return http.StatusConflict, constant.ErrorCodeTaskAlreadyDone

	// 429
	case errors.Is(err, model.ErrTaskQueueFull):
		return http.StatusTooManyRequests, constant.ErrorCodeTaskQueueFull

	// 502
	case errors.Is(err, model.ErrSDKFailedToCancel):
		return http.StatusBadGateway, constant.ErrorSDKFailedToCancel

	// 500
	default:
		return http.StatusInternalServerError, constant.ErrorCodeInternal
	}
}
