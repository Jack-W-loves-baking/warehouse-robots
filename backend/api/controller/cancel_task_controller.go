package controller

import "net/http"

// ICancelTaskController handles HTTP requests to cancel a task.
//
// DELETE Request:
//   - Path:   taskId resolved via r.PathValue("taskId").
//
// Responses:
//   - 204 No Content: if we successfully delete the running task.
//   - 404 Bad Request: Task id not found in the database.
//   - 500 Internal Server Error: unexpected failures.
//
// The controller translates service-layer errors into appropriate HTTP responses.
type ICancelTaskController interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
