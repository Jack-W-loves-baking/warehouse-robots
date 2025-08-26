package dao

import (
	"fmt"
	"sync"
	"time"
	"warehouse-robots/backend/api/model"
)

// InMemoryTaskRepository is a thread-safe in-memory implementation of ITaskRepository.
// It uses a sync.RWMutex to protect concurrent access to the tasks map.
// This should be replaced with a real database in prod.
type InMemoryTaskRepository struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex // protects concurrent reads/writes to tasks
}

func NewInMemoryTaskRepository() ITaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*model.Task),
	}
}

// Create adds a new task to the repository.
// It returns an error if the taskID already exists.
// A copy of the task is stored to prevent external code from mutating internal state.
func (r *InMemoryTaskRepository) Create(task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.TaskID]; exists {
		return fmt.Errorf("task %s already exists", task.TaskID)
	}

	// Store a copy to avoid external modifications
	taskCopy := *task
	r.tasks[task.TaskID] = &taskCopy

	return nil
}

// GetById retrieves a task by ID.
func (r *InMemoryTaskRepository) GetById(taskID string) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	// Return a copy to avoid race conditions
	taskCopy := *task
	if task.CurrentPosition != nil {
		posCopy := *task.CurrentPosition
		taskCopy.CurrentPosition = &posCopy
	}

	return &taskCopy, nil
}

// Update replaces an existing task with the provided one.
// The task must already exist in the repository.
func (r *InMemoryTaskRepository) Update(task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.TaskID]; !exists {
		return fmt.Errorf("task %s not found", task.TaskID)
	}

	task.UpdatedAt = time.Now()
	taskCopy := *task
	r.tasks[task.TaskID] = &taskCopy

	return nil
}

// GetByRobotId returns all tasks belonging to a given robot.
func (r *InMemoryTaskRepository) GetByRobotId(robotID string) ([]*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []*model.Task
	for _, task := range r.tasks {
		if task.RobotID == robotID {
			taskCopy := *task
			if task.CurrentPosition != nil {
				posCopy := *task.CurrentPosition
				taskCopy.CurrentPosition = &posCopy
			}
			tasks = append(tasks, &taskCopy)
		}
	}

	return tasks, nil
}

// UpdateStatus updates the status of an existing task.
func (r *InMemoryTaskRepository) UpdateStatus(taskID string, status model.TaskStatus, errorMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.Status = status
	if errorMsg != "" {
		task.Error = errorMsg
	}
	task.UpdatedAt = time.Now()

	return nil
}

// UpdatePosition updates the current position of the task and its status.
func (r *InMemoryTaskRepository) UpdatePosition(taskID string, position *model.Position, status model.TaskStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.CurrentPosition = position
	task.Status = status
	task.UpdatedAt = time.Now()

	return nil
}
