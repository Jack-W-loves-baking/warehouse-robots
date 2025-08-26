package sdkService

import (
	"warehouse-robots/backend/api/model"
	"warehouse-robots/backend/config"
	mockSdk "warehouse-robots/backend/infra/sdkService/mock"
)

type RobotSDKFactory struct {
	config *config.Config
}

// NewRobotSDKFactory creates a new factory with configuration
func NewRobotSDKFactory(cfg *config.Config) *RobotSDKFactory {
	return &RobotSDKFactory{
		config: cfg,
	}
}

// CreateRobotSDKService creates either mock or real SDK service based on configuration
func (f *RobotSDKFactory) CreateRobotSDKService() model.Warehouse {
	if f.config.Robot.EnableMock {
		return mockSdk.NewMockWarehouse()
	}

	// Return real implementation when available
	return mockSdk.NewMockWarehouse()
}
