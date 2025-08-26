package service

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/dao"
	"warehouse-robots/backend/api/dtos"
	"warehouse-robots/backend/api/manager"
	"warehouse-robots/backend/api/model"
)

// CreateTaskServiceImpl coordinates validation, enqueue, and monitoring of robot tasks.
// It retrieves the target robot from the warehouse SDK, validates the command
// plan against warehouse bounds, persists a task record, and starts background
// monitoring to keep the task status and position up to date.
type CreateTaskServiceImpl struct {
	warehouse   model.Warehouse
	repository  dao.ITaskRepository
	taskMonitor *manager.TaskMonitor
}

// NewCreateTaskService constructs a CreateTaskServiceImpl with the provided
// warehouse SDK handle and task repository.
func NewCreateTaskService(
	warehouse model.Warehouse,
	repository dao.ITaskRepository,
) *CreateTaskServiceImpl {
	return &CreateTaskServiceImpl{
		warehouse:   warehouse,
		repository:  repository,
		taskMonitor: manager.NewTaskMonitor(repository),
	}
}

// CreateTask validates and enqueues a new task for the given robot.
// Returns a TaskInfo snapshot for the newly created task or an error.
func (s *CreateTaskServiceImpl) CreateTask(robotID string, req dtos.CreateTaskRequest) (*dtos.TaskInfo, error) {
	robot, err := s.getRobotByID(robotID)
	if err != nil {
		log.Printf("robot not found: %w", err)
		return nil, model.ErrRobotNotFound
	}

	tasks, err := s.repository.GetByRobotId(robotID)
	if err != nil {
		return nil, model.ErrTaskNotFound
	}

	startPos, err := s.calculateStartPosition(robotID, tasks)
	if err != nil {
		return nil, err
	}

	if err := s.validateBoundary(startPos, req.Commands); err != nil {
		return nil, err
	}

	taskID, posCh, errCh := robot.EnqueueTask(req.Commands)

	task := &model.Task{
		TaskID:          taskID,
		RobotID:         robotID,
		Commands:        req.Commands,
		Status:          model.TaskStatusPending,
		CurrentPosition: nil, // updated by monitor
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repository.Create(task); err != nil {
		// The SDK has already accepted the task; still start monitoring, but return the persistence error.
		s.taskMonitor.StartMonitoring(taskID, posCh, errCh)
		log.Printf("repo.Create task=%s robot=%s failed: %v", taskID, robotID, err)
		return nil, err
	}

	s.taskMonitor.StartMonitoring(taskID, posCh, errCh)

	return &dtos.TaskInfo{
		TaskID:    taskID,
		RobotID:   robotID,
		Status:    dtos.TaskStatusPending,
		Commands:  req.Commands,
		CreatedAt: task.CreatedAt,
	}, nil
}

// getRobotByID resolves a robot from the warehouse by numeric string ID.
// The robotID is expected to be a base-10 string representing a zero-based index
// into the slice returned by warehouse.Robots() (e.g., "0", "1", ...).
func (s *CreateTaskServiceImpl) getRobotByID(robotID string) (model.Robot, error) {
	robotIndex, err := strconv.Atoi(robotID)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID %q: %w", robotID, err)
	}

	robots := s.warehouse.Robots()
	if robotIndex < 0 || robotIndex >= len(robots) {
		return nil, fmt.Errorf("robot index %d out of range (0-%d)", robotIndex, len(robots)-1)
	}
	return robots[robotIndex], nil
}

// calculateStartPosition determines the robot’s starting point when queuing a new task.
//
// Policy:
//   - Only one active task per robot is allowed; if a PENDING/RUNNING task exists, reject.
//   - Use the most recent TERMINAL task to derive the next start:
//   - COMPLETED or FAILED or CANCELLED → use its last known CurrentPosition.
//   - If no prior task exists, default to (0,0,false).
//
// Returns the computed starting position or an error if the request should be rejected.
func (s *CreateTaskServiceImpl) calculateStartPosition(robotID string, tasks []*model.Task) (*model.Position, error) {
	for _, task := range tasks {
		if task.Status == model.TaskStatusPending {
			log.Printf("task %s is pending", task.TaskID)
			return nil, model.ErrTaskQueueFull
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
	})

	for _, task := range tasks {
		switch task.Status {
		case model.TaskStatusCompleted, model.TaskStatusFailed, model.TaskStatusCancelled:
			if task.CurrentPosition != nil {
				return &model.Position{
					X:        task.CurrentPosition.X,
					Y:        task.CurrentPosition.Y,
					HasCrate: task.CurrentPosition.HasCrate,
				}, nil
			}
		}
	}

	return &model.Position{X: 0, Y: 0, HasCrate: false}, nil
}

// validateBoundary simulates the command sequence from a starting position and
// ensures every step remains within the warehouse bounds.
//
// Commands are case-insensitive single letters; whitespace is ignored.
// Valid moves: N (y+1), S (y-1), E (x+1), W (x-1).
// Returns an error on the first out-of-bounds move or invalid command.
func (s *CreateTaskServiceImpl) validateBoundary(start *model.Position, commands string) error {
	x, y := int(start.X), int(start.Y)

	maxX := constant.WarehouseSizeX - 1
	maxY := constant.WarehouseSizeY - 1
	minX := 0
	minY := 0

	commands = strings.ToUpper(strings.ReplaceAll(commands, " ", ""))

	for i, cmd := range commands {
		originalX, originalY := x, y

		switch cmd {
		case 'N':
			y++
		case 'S':
			y--
		case 'E':
			x++
		case 'W':
			x--
		}

		if x < minX {

			return model.ErrBoundary
		}
		if x > maxX {
			log.Printf("boundary violation: command '%c' at index %d would move robot out of bounds (too far WEST): (%d,%d)->(%d,%d), x_min=%d", cmd, i, originalX, originalY, x, y, minX)
			return model.ErrBoundary
		}
		if y < minY {
			log.Printf("boundary violation: command '%c' at index %d would move robot out of bounds (too far SOUTH): (%d,%d)->(%d,%d), x_min=%d", cmd, i, originalX, originalY, x, y, minX)
			return model.ErrBoundary
		}
		if y > maxY {
			log.Printf("boundary violation: command '%c' at index %d would move robot out of bounds (too far NORTH): (%d,%d)->(%d,%d), x_min=%d", cmd, i, originalX, originalY, x, y, minX)
			return model.ErrBoundary
		}
	}

	return nil
}
