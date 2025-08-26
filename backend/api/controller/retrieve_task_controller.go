package controller

import "net/http"

// IRetrieveTaskController to handle the incoming request to retrieve the status of a task
type IRetrieveTaskController interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
