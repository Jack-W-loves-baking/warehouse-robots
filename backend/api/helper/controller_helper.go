package helper

import (
	"encoding/json"
	"net/http"

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
