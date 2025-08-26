package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"warehouse-robots/backend/api/dao"
	"warehouse-robots/backend/api/model"
)

// TaskMonitor manages background goroutines that listen to channels
type TaskMonitor struct {
	repository dao.ITaskRepository
	monitors   map[string]context.CancelFunc
	mu         sync.Mutex
	wg         sync.WaitGroup
}

func NewTaskMonitor(repo dao.ITaskRepository) *TaskMonitor {
	return &TaskMonitor{
		repository: repo,
		monitors:   make(map[string]context.CancelFunc),
	}
}

// StartMonitoring starts a goroutine to monitor task channels
func (tm *TaskMonitor) StartMonitoring(
	taskID string,
	positionChan <-chan model.RobotState,
	errorChan <-chan error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)

	tm.mu.Lock()
	tm.monitors[taskID] = cancel
	tm.mu.Unlock()

	tm.wg.Add(1)
	go tm.monitorTask(ctx, taskID, positionChan, errorChan)
}

// monitorTask is the goroutine that listens to channels
func (tm *TaskMonitor) monitorTask(
	ctx context.Context,
	taskID string,
	positionChan <-chan model.RobotState,
	errorChan <-chan error,
) {
	defer tm.wg.Done()
	defer tm.cleanup(taskID)

	for {
		select {
		case position, ok := <-positionChan:
			if !ok {
				// Channel closed - task completed successfully
				err := tm.repository.UpdateStatus(taskID, model.TaskStatusCompleted, "")
				if err != nil {
					// Log error but don't return - channel is closed anyway
					fmt.Printf("Error updating status to completed: %v\n", err)
				}
				return
			}

			// Update position in repository
			pos := &model.Position{
				X:        position.X,
				Y:        position.Y,
				HasCrate: position.HasCrate,
			}

			// First position update means task is running
			task, err := tm.repository.GetById(taskID)
			if err != nil {
				// Task might have been deleted, just log and continue
				fmt.Printf("Error getting task %s: %v\n", taskID, err)
				continue
			}

			status := model.TaskStatusRunning
			if task != nil && task.Status == model.TaskStatusCompleted {
				status = task.Status // Don't override completed status
			}

			err = tm.repository.UpdatePosition(taskID, pos, status)
			if err != nil {
				fmt.Printf("Error updating position for task %s: %v\n", taskID, err)
			}
			tm.printTaskSnapshot(taskID)

		case err, ok := <-errorChan:
			if ok && err != nil {
				// Error received - task failed
				updateErr := tm.repository.UpdateStatus(taskID, model.TaskStatusFailed, err.Error())
				if updateErr != nil {
					fmt.Printf("Error updating status to failed: %v\n", updateErr)
				}
				return
			}

		case <-ctx.Done():
			// Timeout or cancelled
			if ctx.Err() == context.DeadlineExceeded {
				err := tm.repository.UpdateStatus(taskID, model.TaskStatusFailed, "task timeout")
				if err != nil {
					fmt.Printf("Error updating status to timeout: %v\n", err)
				}
			}
			// If cancelled, status should already be updated elsewhere
			return
		}
	}
}

func (tm *TaskMonitor) cleanup(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancel, exists := tm.monitors[taskID]; exists {
		cancel()
		delete(tm.monitors, taskID)
	}
}

// CancelTask cancels monitoring for a specific task
func (tm *TaskMonitor) CancelTask(taskID string) error {
	tm.mu.Lock()

	cancel, exists := tm.monitors[taskID]
	if !exists {
		tm.mu.Unlock()
		return fmt.Errorf("no monitor found for task %s", taskID)
	}

	cancel()
	delete(tm.monitors, taskID)
	tm.mu.Unlock() // Unlock before calling repository to avoid potential deadlock

	// Update status in repository
	return tm.repository.UpdateStatus(taskID, model.TaskStatusCancelled, "cancelled by user")
}

// Shutdown gracefully stops all monitors
func (tm *TaskMonitor) Shutdown(ctx context.Context) error {
	tm.mu.Lock()
	// Cancel all monitors
	for _, cancel := range tm.monitors {
		cancel()
	}
	tm.mu.Unlock()

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		tm.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (tm *TaskMonitor) printTaskSnapshot(taskID string) {
	task, err := tm.repository.GetById(taskID)
	if err != nil {
		fmt.Printf("[task %s] snapshot error: %v\n", taskID, err)
		return
	}
	b, _ := json.MarshalIndent(task, "", "  ")
	fmt.Printf("[task %s] %s\n", taskID, string(b))
}
