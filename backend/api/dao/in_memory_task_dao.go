package dao

import (
	"fmt"
	"sync"
	"time"
	"warehouse-robots/backend/api/model"
)

type InMemoryTaskRepository struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex
}

func NewInMemoryTaskRepository() ITaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*model.Task),
	}
}

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

func (r *InMemoryTaskRepository) GetAll() ([]*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskCopy := *task
		if task.CurrentPosition != nil {
			posCopy := *task.CurrentPosition
			taskCopy.CurrentPosition = &posCopy
		}
		tasks = append(tasks, &taskCopy)
	}

	return tasks, nil
}
