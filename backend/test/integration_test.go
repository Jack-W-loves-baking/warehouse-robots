package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"warehouse-robots/backend/api/dtos"
	"warehouse-robots/backend/binder"
	"warehouse-robots/backend/config"
)

func TestIntegration_CreateTask_HappyFlow(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	requestBody := dtos.CreateTaskRequest{
		Commands: "NNNN",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/robots/0/tasks", bytes.NewBuffer(jsonBody))
	req.SetPathValue("robotId", "0")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	container.CreateTaskController.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
		return
	}

	var createResponse dtos.TaskInfo
	if err := json.Unmarshal(w.Body.Bytes(), &createResponse); err != nil {
		t.Errorf("Failed to unmarshal create response: %v", err)
		return
	}

	if createResponse.TaskID == "" {
		t.Error("Expected task ID to be set")
	}

	if createResponse.RobotID != "0" {
		t.Errorf("Expected robot ID '0', got '%s'", createResponse.RobotID)
	}

	if createResponse.Commands != "NNNN" {
		t.Errorf("Expected commands 'NNNN', got '%s'", createResponse.Commands)
	}

	if createResponse.Status != dtos.TaskStatusPending {
		t.Errorf("Expected status 'PENDING', got '%s'", createResponse.Status)
	}

	time.Sleep(100 * time.Millisecond)

	// Now retrieve the task to see its updated status
	getReq := httptest.NewRequest("GET", "/api/tasks/"+createResponse.TaskID, nil)
	getReq.SetPathValue("taskId", createResponse.TaskID)

	getW := httptest.NewRecorder()
	container.RetrieveTaskController.Handle(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected get status code %d, got %d", http.StatusOK, getW.Code)
		t.Errorf("Get response body: %s", getW.Body.String())
		return
	}

	var getResponse dtos.TaskInfo
	if err := json.Unmarshal(getW.Body.Bytes(), &getResponse); err != nil {
		t.Errorf("Failed to unmarshal get response: %v", err)
		return
	}

	if getResponse.TaskID != createResponse.TaskID {
		t.Errorf("Expected same task ID, got create='%s' get='%s'", createResponse.TaskID, getResponse.TaskID)
	}

	// Task should still be pending or might have started processing
	if getResponse.Status != dtos.TaskStatusPending && getResponse.Status != dtos.TaskStatusCompleted {
		t.Errorf("Expected status 'PENDING' or 'COMPLETED', got '%s'", getResponse.Status)
	}
}

func TestIntegration_CreateTask_BoundaryViolation(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	// Try to create a task that would move robot out of bounds
	requestBody := dtos.CreateTaskRequest{
		Commands: "SSSSSSSSSSSSSSS", // Move south 15 times from (0,0) - should fail
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/robots/0/tasks", bytes.NewBuffer(jsonBody))
	req.SetPathValue("robotId", "0")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	container.CreateTaskController.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for boundary violation, got %d", http.StatusBadRequest, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
		return
	}

	var errorResponse dtos.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
		return
	}

	if errorResponse.Code == "" {
		t.Error("Expected error code to be set")
	}
}

func TestIntegration_CreateTask_InvalidRobotId(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	requestBody := dtos.CreateTaskRequest{
		Commands: "NNNN",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/robots/999/tasks", bytes.NewBuffer(jsonBody))
	req.SetPathValue("robotId", "999")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	container.CreateTaskController.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for invalid robot, got %d", http.StatusNotFound, w.Code)
		return
	}

	var errorResponse dtos.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
		return
	}

	if errorResponse.Code == "" {
		t.Error("Expected error code to be set")
	}
}

func TestIntegration_CreateTask_InvalidCommands(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	testCases := []struct {
		name     string
		commands string
	}{
		{"empty commands", ""},
		{"invalid character", "NXXX"},
		{"mixed valid and invalid", "NNSZ"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBody := dtos.CreateTaskRequest{Commands: tc.commands}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/api/robots/0/tasks", bytes.NewBuffer(jsonBody))
			req.SetPathValue("robotId", "0")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			container.CreateTaskController.Handle(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status code %d for %s, got %d", http.StatusBadRequest, tc.name, w.Code)
				return
			}

			var errorResponse dtos.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
				t.Errorf("Failed to unmarshal error response: %v", err)
				return
			}

			if errorResponse.Code == "" {
				t.Error("Expected error code to be set")
			}
		})
	}
}

func TestIntegration_GetTask_NotFound(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	req := httptest.NewRequest("GET", "/api/tasks/nonexistent", nil)
	req.SetPathValue("taskId", "nonexistent")

	w := httptest.NewRecorder()
	container.RetrieveTaskController.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		return
	}

	var errorResponse dtos.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
		return
	}

	if errorResponse.Code == "" {
		t.Error("Expected error code to be set")
	}
}

func TestIntegration_CancelTask_NotFound(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	req := httptest.NewRequest("DELETE", "/api/tasks/nonexistent", nil)
	req.SetPathValue("taskId", "nonexistent")

	w := httptest.NewRecorder()
	container.CancelTaskController.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		return
	}
}

func TestIntegration_CreateTask_RobotBusy(t *testing.T) {
	cfg := &config.Config{
		Robot: config.RobotConfig{
			EnableMock: true,
		},
	}

	container := binder.NewContainer(cfg)

	// Create first task
	requestBody1 := dtos.CreateTaskRequest{Commands: "NN"}
	jsonBody1, _ := json.Marshal(requestBody1)

	req1 := httptest.NewRequest("POST", "/api/robots/0/tasks", bytes.NewBuffer(jsonBody1))
	req1.SetPathValue("robotId", "0")
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	container.CreateTaskController.Handle(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("Expected first task creation to succeed, got %d", w1.Code)
		return
	}

	// Try to create second task immediately (should be rejected if robot is busy)
	requestBody2 := dtos.CreateTaskRequest{Commands: "SS"}
	jsonBody2, _ := json.Marshal(requestBody2)

	req2 := httptest.NewRequest("POST", "/api/robots/0/tasks", bytes.NewBuffer(jsonBody2))
	req2.SetPathValue("robotId", "0")
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	container.CreateTaskController.Handle(w2, req2)

	// This should fail with a conflict, bad request, or too many requests
	if w2.Code == http.StatusCreated {
		t.Error("Expected second task creation to fail due to robot being busy")
		return
	}

	if w2.Code != http.StatusBadRequest && w2.Code != http.StatusConflict && w2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, %d, or %d for busy robot, got %d",
			http.StatusBadRequest, http.StatusConflict, http.StatusTooManyRequests, w2.Code)
		t.Errorf("Response: %s", w2.Body.String())
	}
}
