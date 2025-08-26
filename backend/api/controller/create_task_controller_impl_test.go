package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"warehouse-robots/backend/api/dtos"
)

// Mock service for testing
type mockCreateTaskService struct{}

func (m *mockCreateTaskService) CreateTask(robotID string, req dtos.CreateTaskRequest) (*dtos.TaskInfo, error) {
	return &dtos.TaskInfo{
		TaskID:   "test-123",
		RobotID:  robotID,
		Commands: req.Commands,
	}, nil
}

func TestCreateTaskController_Success(t *testing.T) {
	controller := NewCreateTaskController(&mockCreateTaskService{})

	reqBody := dtos.CreateTaskRequest{Commands: "NESW"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/robots/1/tasks", bytes.NewBuffer(jsonBody))
	req.SetPathValue("robotId", "1")
	w := httptest.NewRecorder()

	controller.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestCreateTaskController_InvalidJSON(t *testing.T) {
	controller := NewCreateTaskController(&mockCreateTaskService{})

	req := httptest.NewRequest("POST", "/robots/1/tasks", bytes.NewBufferString("invalid json"))
	req.SetPathValue("robotId", "1")
	w := httptest.NewRecorder()

	controller.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestValidateCommands(t *testing.T) {
	controller := &CreateTaskControllerImpl{}

	tests := []struct {
		commands    string
		expectError bool
	}{
		{"NESW", false},
		{"nesw", false},
		{"N E S W", false},
		{"", true},
		{"   ", true},
		{"NXS", true},
		{"N123", true},
	}

	for _, test := range tests {
		err := controller.validateCommands(test.commands)
		if test.expectError && err == nil {
			t.Errorf("Expected error for commands: %s", test.commands)
		}
		if !test.expectError && err != nil {
			t.Errorf("Unexpected error for commands: %s, error: %v", test.commands, err)
		}
	}
}
