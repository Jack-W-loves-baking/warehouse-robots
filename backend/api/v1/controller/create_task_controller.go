package controller

import "net/http"

// ICreateTaskController to handle the incoming request to create a task
type ICreateTaskController interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
