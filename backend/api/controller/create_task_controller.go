package controller

import "net/http"

// ICreateTaskController processes POST /robots/{robotId}/tasks requests.
//
// Request:
//   - Path:   robotId (string) resolved via r.PathValue("robotId").
//   - Body:   dtos.CreateTaskRequest (JSON).
//
// Responses:
//   - 201 Created: on successful creation, returns dtos.TaskInfo.
//   - 400 Bad Request: invalid JSON or invalid command sequence.
//   - 429 Too many requests: the queue has not Terminated, so we cannot queue a new task.
//   - 503 Service Unavailable: no robots available.
//   - 500 Internal Server Error: unexpected failures.
//
// Error bodies are standardized via ControllerHelper.
type ICreateTaskController interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
