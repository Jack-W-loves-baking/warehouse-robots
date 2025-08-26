package service

import (
	"log"
	"strconv"
	"time"

	"warehouse-robots/backend/api/dao"
	"warehouse-robots/backend/api/manager"
	"warehouse-robots/backend/api/model"
)

type CancelTaskServiceImpl struct {
	warehouse   model.Warehouse
	repository  dao.ITaskRepository
	taskMonitor *manager.TaskMonitor
}

// NewCancelTaskService constructor
func NewCancelTaskService(
	warehouse model.Warehouse,
	repository dao.ITaskRepository) ICancelTaskService {
	return &CancelTaskServiceImpl{
		warehouse:   warehouse,
		repository:  repository,
		taskMonitor: manager.NewTaskMonitor(repository),
	}
}

// CancelTaskById cancels a task.
// Rules:
//   - If task is TERMINAL (COMPLETED/FAILED/CANCELLED): reject.
//   - If task is RUNNING or PENDING: attempt SDK CancelTask with retries; on success, stop monitor and mark CANCELLED.
//     On repeated failure, mark FAILED with an explanatory error.
func (s *CancelTaskServiceImpl) CancelTaskById(taskId string) error {
	task, err := s.repository.GetById(taskId)
	if err != nil {
		log.Printf("get task: %w", err)
		return model.ErrTaskNotFound
	}

	switch task.Status {
	// If the task is COMPLETED, FAILED OR CANCELLED,we reject
	case model.TaskStatusCompleted, model.TaskStatusFailed, model.TaskStatusCancelled:
		log.Printf("task %s is already %s", taskId, task.Status)
		return model.ErrTaskProcessed

	case model.TaskStatusPending:
		robot, err := s.getRobotByRobotID(task.RobotID)
		if err != nil {
			log.Printf("resolve robot %q: %w", task.RobotID, err)
			return model.ErrRobotNotFound
		}

		// Retry SDK cancel a few times
		const maxRetries = 3
		baseDelay := 100 * time.Millisecond
		var lastErr error
		for i := 0; i < maxRetries; i++ {
			if err := robot.CancelTask(taskId); err == nil {
				// SDK accepted so we need to update the task status to CANCELLED
				if monErr := s.taskMonitor.CancelTask(taskId); monErr != nil {
					// If no monitor found, still update status explicitly
					_ = s.repository.UpdateStatus(taskId, model.TaskStatusCancelled, "cancelled by user")
				}
				return nil
			} else {
				lastErr = err
				time.Sleep(time.Duration(1<<i) * baseDelay)
			}
		}

		// Could not cancel in SDK even after retry, then dont do anything
		log.Printf("sdk cancel failed for task %s after %d retries: %v", taskId, maxRetries, lastErr)
		return model.ErrSDKFailedToCancel
	}

	return nil
}

// getRobotByID resolves a robot from the warehouse by numeric string ID.
// The robotID is expected to be a base-10 string representing a zero-based index
// into the slice returned by warehouse.Robots() (e.g., "0", "1", ...).
func (s *CancelTaskServiceImpl) getRobotByRobotID(robotID string) (model.Robot, error) {
	robotIndex, err := strconv.Atoi(robotID)
	if err != nil {
		return nil, model.ErrRobotIDInvalid
	}

	robots := s.warehouse.Robots()
	if robotIndex < 0 || robotIndex >= len(robots) {
		log.Printf("robot index %d out of range (0-%d)\", robotIndex, len(robots)-1")
		return nil, model.ErrRobotIDInvalid
	}
	return robots[robotIndex], nil
}
