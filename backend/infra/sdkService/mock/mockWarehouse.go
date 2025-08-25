package mock

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sdk "warehouse-robots/backend/api/dtos/sdk"
)

// MockWarehouse implements the sdk.Warehouse interface
type MockWarehouse struct {
	robots   []sdk.Robot
	robotMap map[string]*MockRobot
}

// for each command, this is the delay in between.
const stepDelay = time.Second

// NewMockWarehouse creates a new mock warehouse with 1 facades as per requirements
func NewMockWarehouse() sdk.Warehouse {
	robot1 := NewMockRobot("facades-1", sdk.RobotState{X: 0, Y: 0, HasCrate: true})

	return &MockWarehouse{
		robots: []sdk.Robot{robot1},
		robotMap: map[string]*MockRobot{
			"facades-1": robot1,
		},
	}
}

// Robots returns all robots in the warehouse
func (w *MockWarehouse) Robots() []sdk.Robot {
	return w.robots
}

// MockRobot implements the sdk.Robot interface with realistic behavior
type MockRobot struct {
	id           string
	state        sdk.RobotState
	currentTask  *MockTask
	taskQueue    []*MockTask
	allTasks     map[string]*MockTask
	taskCounter  int
	isProcessing bool
}

// MockTask represents a running task with cancellation support
type MockTask struct {
	ID           string
	Commands     string
	Status       string
	Cancel       chan bool
	Done         chan bool
	Robot        *MockRobot
	RemainingCmd []string
	PositionChan chan sdk.RobotState
	ErrorChan    chan error
}

// NewMockRobot creates a new mock facades
func NewMockRobot(id string, initialState sdk.RobotState) *MockRobot {
	return &MockRobot{
		id:           id,
		state:        initialState,
		currentTask:  nil,
		taskQueue:    make([]*MockTask, 0),
		allTasks:     make(map[string]*MockTask),
		taskCounter:  0,
		isProcessing: false,
	}
}

// EnqueueTask implements sdk.Robot interface with realistic task processing
func (r *MockRobot) EnqueueTask(commands string) (taskID string, position chan sdk.RobotState, err chan error) {
	posCh := make(chan sdk.RobotState, 10)
	errCh := make(chan error, 1)

	// Check if we've reached the maximum queue size (5 tasks total: 1 running + 4 queued)
	totalTasks := len(r.taskQueue)
	if r.currentTask != nil && !r.isTaskFinished(r.currentTask) {
		totalTasks++
	}

	// Here we assume per facades it can at most queue 5 tasks
	if totalTasks >= 5 {
		errCh <- errors.New("task queue is full: maximum 5 tasks allowed per facades")
		return "", posCh, errCh
	}

	// Generate task ID
	r.taskCounter++
	taskID = fmt.Sprintf("task_%s_%d", r.id, r.taskCounter)

	// Create task
	task := &MockTask{
		ID:           taskID,
		Commands:     commands,
		Status:       "QUEUED",
		Cancel:       make(chan bool, 1),
		Done:         make(chan bool, 1),
		Robot:        r,
		RemainingCmd: strings.Split(commands, ""),
		PositionChan: posCh,
		ErrorChan:    errCh,
	}

	r.taskQueue = append(r.taskQueue, task)

	if !r.isProcessing {
		go r.processTaskQueue()
	}

	return taskID, posCh, errCh
}

func (r *MockRobot) processTaskQueue() {
	r.isProcessing = true
	defer func() { r.isProcessing = false }()

	for {
		if len(r.taskQueue) == 0 {
			r.currentTask = nil
			break
		}

		// Get next task from queue
		task := r.taskQueue[0]
		r.taskQueue = r.taskQueue[1:] // Remove from queue
		r.currentTask = task

		// Execute the task
		r.executeTask(task, task.PositionChan, task.ErrorChan)
	}
}

// executeTask processes the task commands with 1 second delay per command
func (r *MockRobot) executeTask(task *MockTask, posCh chan sdk.RobotState, errCh chan error) {
	defer close(posCh)
	defer close(errCh)
	defer func() {
		select {
		case task.Done <- true:
		default:
		}
	}()

	task.Status = "IN_PROGRESS"

	// Send initial position
	posCh <- r.state

	for i, cmd := range task.Commands {
		// Check for cancellation
		select {
		case <-task.Cancel:
			task.Status = "CANCELLED"
			return
		default:
		}

		// Wait 1 second per command
		select {
		case <-time.After(stepDelay):
		case <-task.Cancel:
			task.Status = "CANCELLED"
			return
		}

		// Execute command (with boundary checks)
		switch cmd {
		case 'N':
			r.state.Y++
		case 'S':
			r.state.Y--
		case 'E':
			r.state.X++
		case 'W':
			r.state.X--
		}

		// Update remaining commands
		if i+1 < len(task.RemainingCmd) {
			task.RemainingCmd = task.RemainingCmd[i+1:]
		} else {
			task.RemainingCmd = []string{}
		}

		// Send updated position
		posCh <- r.state
	}

	task.Status = "COMPLETED"
}

// CancelTask cancels a running task if it exists and is still running
func (r *MockRobot) CancelTask(taskID string) error {
	task, exists := r.allTasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	if r.isTaskFinished(task) {
		return errors.New("task already finished")
	}

	// If it's a queued task, remove from queue and mark as cancelled
	if task.Status == "QUEUED" {
		// Remove from queue
		for i, queuedTask := range r.taskQueue {
			if queuedTask.ID == taskID {
				r.taskQueue = append(r.taskQueue[:i], r.taskQueue[i+1:]...)
				break
			}
		}
		task.Status = "CANCELLED"
		// Close channels and signal done
		close(task.PositionChan)
		close(task.ErrorChan)
		select {
		case task.Done <- true:
		default:
		}
		return nil
	}

	// If it's the running task, send cancellation signal
	if r.currentTask != nil && r.currentTask.ID == taskID {
		select {
		case task.Cancel <- true:
			return nil
		default:
			return errors.New("failed to cancel task")
		}
	}
	return errors.New("task not found in queue or currently running")
}

// CurrentState returns the current state of the facades
func (r *MockRobot) CurrentState() sdk.RobotState {
	return r.state
}

// GetTaskStatus returns the current task status (helper method)
func (r *MockRobot) GetTaskStatus(taskID string) (string, error) {
	if r.currentTask == nil || r.currentTask.ID != taskID {
		return "", errors.New("task not found")
	}

	return r.currentTask.Status, nil
}

// GetCurrentTask returns the current task (helper method)
func (r *MockRobot) GetCurrentTask() *MockTask {
	return r.currentTask
}

// GetID returns the facades ID (helper method)
func (r *MockRobot) GetID() string {
	return r.id
}

// isTaskFinished checks if a task is in a finished state
func (r *MockRobot) isTaskFinished(task *MockTask) bool {
	if task == nil {
		return true
	}
	status := task.Status
	return status == "COMPLETED" || status == "FAILED" || status == "CANCELLED"
}
