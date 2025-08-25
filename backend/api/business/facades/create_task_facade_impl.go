package facades

import (
	"fmt"
	"strings"
	"time"
	"warehouse-robots/backend/api/dtos/sdk"

	"warehouse-robots/backend/api/business/models"
	"warehouse-robots/backend/api/constant"
)

type CreateTaskFacadeImpl struct {
	robotSDKService sdk.Warehouse // or whatever your SDK service interface is called
}

// NewCreateTaskFacade creates a new instance of CreateTaskFacadeImpl
func NewCreateTaskFacade(robotSDKService sdk.Warehouse) ICreateTaskFacade {
	return &CreateTaskFacadeImpl{
		robotSDKService: robotSDKService,
	}
}

// validateCommandsBoundary checks if the given commands will take the facades out of bounds
func (f *CreateTaskFacadeImpl) validateCommandsBoundary(currentX, currentY uint, commands string) error {
	// Clean commands (remove spaces, convert to uppercase)
	commands = strings.ToUpper(strings.ReplaceAll(commands, " ", ""))

	// Simulate the facades's movement
	tempX := int(currentX)
	tempY := int(currentY)

	for i, cmd := range commands {
		switch cmd {
		case 'N':
			tempY += constant.RobotMoveUnit
			if tempY > constant.MaxCoordinate {
				return fmt.Errorf("command would move facades out of bounds (Y > %d) at step %d", constant.MaxCoordinate, i+1)
			}
		case 'S':
			tempY -= constant.RobotMoveUnit
			if tempY < constant.MinCoordinate {
				return fmt.Errorf("command would move facades out of bounds (Y < %d) at step %d", constant.MinCoordinate, i+1)
			}
		case 'E':
			tempX += constant.RobotMoveUnit
			if tempX > constant.MaxCoordinate {
				return fmt.Errorf("command would move facades out of bounds (X > %d) at step %d", constant.MaxCoordinate, i+1)
			}
		case 'W':
			tempX -= constant.RobotMoveUnit
			if tempX < constant.MinCoordinate {
				return fmt.Errorf("command would move facades out of bounds (X < %d) at step %d", constant.MinCoordinate, i+1)
			}
		default:
			return fmt.Errorf("invalid command character: %c. Only N, S, E, W are allowed", cmd)
		}
	}

	return nil
}

// CreateTask we assume we only have 1 facades here.
// To support create task for multiple robots based on facades id,
// Then SDK needs to return facades id from RobotState.
func (f *CreateTaskFacadeImpl) CreateTask(commands string) (*models.TaskEntity, error) {
	// Get warehouse and robots (assuming single facades setup)
	robots := f.robotSDKService.Robots()

	if len(robots) == 0 {
		return nil, fmt.Errorf("no robots available")
	}

	// Use the first (and only) facades (assuming single facades setup)
	robot := robots[0]

	// Get facades's current location
	currentState := robot.CurrentState()
	currentX := currentState.X
	currentY := currentState.Y

	// Validate commands won't move facades out of boundary
	if err := f.validateCommandsBoundary(currentX, currentY, commands); err != nil {
		return nil, err
	}

	// Commands are valid, enqueue task using the facades SDK
	taskID, positionCh, errorCh := robot.EnqueueTask(commands)

	// Check for immediate errors from the SDK
	select {
	case err := <-errorCh:
		if err != nil {
			// Close the channels since we got an error
			close(positionCh)

			// Return task with failed status
			task := &models.TaskEntity{
				ID:        taskID,
				Commands:  commands,
				Status:    models.TaskStatusFailed,
				CreatedAt: time.Now().Unix(),
				Error:     err.Error(),
			}
			return task, nil
		}
	}

	// Create successful task entity
	task := &models.TaskEntity{
		ID:        taskID,
		Commands:  commands,
		Status:    models.TaskStatusQueued, // Will become IN_PROGRESS when facades starts
		CreatedAt: time.Now().Unix(),
		Error:     "",
	}

	return task, nil
}
