package manager

import (
	"context"
	"fmt"
	"sync"
	"time"
	"warehouse-robots/backend/api/dao"
	"warehouse-robots/backend/api/model"
)

// TaskMonitor manages the lifecycle of goroutines that watch robot task channels.
// It ensures updates (status, position, errors) are persisted into the repository
// and provides graceful shutdown.
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

// StartMonitoring creates a goroutine that listens for position and error events
// from a robot task. It registers a cancel function to allow external shutdown.
// Each monitor has a maximum lifetime of 30 minutes.
func (tm *TaskMonitor) StartMonitoring(
	taskID string,
	positionChan <-chan model.RobotState,
	errorChan <-chan error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)

	tm.mu.Lock()
	tm.monitors[taskID] = cancel
	tm.mu.Unlock()

	// Increment WaitGroup counter before starting the goroutine.
	// This ensures Shutdown() can wait for this monitor to exit
	tm.wg.Add(1)
	go tm.monitorTask(ctx, taskID, positionChan, errorChan)
}

// monitorTask is the goroutine that listens to channels
// It updates task state in the repository until the task completes, fails, or times out.
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

			status := model.TaskStatusPending
			if task != nil && task.Status == model.TaskStatusCompleted {
				status = task.Status // Don't override completed status
			}

			err = tm.repository.UpdatePosition(taskID, pos, status)
			if err != nil {
				fmt.Printf("Error updating position for task %s: %v\n", taskID, err)
			}

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

// cleanup removes the monitor for a task and cancels its context.
func (tm *TaskMonitor) cleanup(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancel, exists := tm.monitors[taskID]; exists {
		cancel()
		delete(tm.monitors, taskID)
	}
}

// CancelTask stops monitoring and cancels the task explicitly.
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
