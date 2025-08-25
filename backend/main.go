package main

import (
	"fmt"
	"log"
	"net/http"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/middleware"
	"warehouse-robots/backend/binder"
	"warehouse-robots/backend/config"
)

/**
 * Main entrance for the app
 */
func main() {
	// Load configuration
	config := config.Load()

	container := binder.NewContainer(config)

	// Register routes
	mux := http.NewServeMux()

	mux.HandleFunc(constant.RouteCreateTask, container.CreateTaskController.Handle)

	// Apply middleware stack with configuration
	handler := middleware.Chain(mux,
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware(config),
		middleware.JSONMiddleware,
	)

	// Start server with configured address
	address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Starting server on %s", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
