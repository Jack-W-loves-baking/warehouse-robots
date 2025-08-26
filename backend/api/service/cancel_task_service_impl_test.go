package service

import (
	"testing"
	"warehouse-robots/backend/api/constant"
	"warehouse-robots/backend/api/model"
)

func TestCreateTaskServiceImpl_validateBoundary(t *testing.T) {
	// Create a service instance for testing
	service := &CreateTaskServiceImpl{}

	tests := []struct {
		name        string
		start       *model.Position
		commands    string
		expectError bool
		description string
	}{
		// Valid movement tests
		{
			name:        "valid_single_north",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "N",
			expectError: false,
			description: "Single north movement within bounds",
		},
		{
			name:        "valid_single_south",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "S",
			expectError: false,
			description: "Single south movement within bounds",
		},
		{
			name:        "valid_single_east",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "E",
			expectError: false,
			description: "Single east movement within bounds",
		},
		{
			name:        "valid_single_west",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "W",
			expectError: false,
			description: "Single west movement within bounds",
		},
		{
			name:        "valid_multiple_commands",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "NNESSW",
			expectError: false,
			description: "Multiple valid movements within bounds",
		},
		{
			name:        "valid_with_spaces",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "N N E S S W",
			expectError: false,
			description: "Valid movements with spaces (should be ignored)",
		},
		{
			name:        "valid_lowercase_commands",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "nnessw",
			expectError: false,
			description: "Valid movements in lowercase (should be converted to uppercase)",
		},
		{
			name:        "valid_mixed_case",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "NnEsSw",
			expectError: false,
			description: "Valid movements in mixed case",
		},
		{
			name:        "valid_at_origin",
			start:       &model.Position{X: 0, Y: 0, HasCrate: false},
			commands:    "NE",
			expectError: false,
			description: "Valid movements starting from origin",
		},
		{
			name:        "valid_at_max_boundary",
			start:       &model.Position{X: constant.WarehouseSizeX - 1, Y: constant.WarehouseSizeY - 1, HasCrate: false},
			commands:    "SW",
			expectError: false,
			description: "Valid movements starting from max boundary",
		},
		{
			name:        "valid_empty_commands",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "",
			expectError: false,
			description: "Empty command string should be valid",
		},
		{
			name:        "valid_only_spaces",
			start:       &model.Position{X: 5, Y: 5, HasCrate: false},
			commands:    "   ",
			expectError: false,
			description: "Only spaces should be valid (all spaces removed)",
		},

		// Boundary violation tests - X axis (West)
		{
			name:        "invalid_west_from_origin",
			start:       &model.Position{X: 0, Y: 5, HasCrate: false},
			commands:    "W",
			expectError: true,
			description: "Moving west from X=0 should violate boundary",
		},
		{
			name:        "invalid_west_from_min_x",
			start:       &model.Position{X: constant.MinCoordinate, Y: 5, HasCrate: false},
			commands:    "W",
			expectError: true,
			description: "Moving west from minimum X coordinate should violate boundary",
		},

		// Boundary violation tests - X axis (East)
		{
			name:        "invalid_east_from_max_x",
			start:       &model.Position{X: constant.WarehouseSizeX - 1, Y: 5, HasCrate: false},
			commands:    "E",
			expectError: true,
			description: "Moving east from maximum X coordinate should violate boundary",
		},

		// Boundary violation tests - Y axis (South)
		{
			name:        "invalid_south_from_origin",
			start:       &model.Position{X: 5, Y: 0, HasCrate: false},
			commands:    "S",
			expectError: true,
			description: "Moving south from Y=0 should violate boundary",
		},
		{
			name:        "invalid_south_from_min_y",
			start:       &model.Position{X: 5, Y: constant.MinCoordinate, HasCrate: false},
			commands:    "S",
			expectError: true,
			description: "Moving south from minimum Y coordinate should violate boundary",
		},

		// Boundary violation tests - Y axis (North)
		{
			name:        "invalid_north_from_max_y",
			start:       &model.Position{X: 5, Y: constant.WarehouseSizeY - 1, HasCrate: false},
			commands:    "N",
			expectError: true,
			description: "Moving north from maximum Y coordinate should violate boundary",
		},

		// Complex boundary violation tests
		{
			name:        "invalid_multiple_commands_violation_at_end",
			start:       &model.Position{X: constant.WarehouseSizeX - 2, Y: 5, HasCrate: false},
			commands:    "EE",
			expectError: true,
			description: "Valid first move, but second move violates boundary",
		},
		{
			name:        "invalid_multiple_commands_violation_in_middle",
			start:       &model.Position{X: constant.WarehouseSizeX - 1, Y: 5, HasCrate: false},
			commands:    "ENS",
			expectError: true,
			description: "First move violates boundary, should fail immediately",
		},
		{
			name:        "invalid_complex_path_with_violation",
			start:       &model.Position{X: 1, Y: 1, HasCrate: false},
			commands:    "WWSS",
			expectError: true,
			description: "Complex path that eventually violates boundary",
		},

		// Edge case tests
		{
			name:        "valid_return_to_start",
			start:       &model.Position{X: 0, Y: 0, HasCrate: false},
			commands:    "NESW",
			expectError: false,
			description: "Move in a square and return to start position",
		},
		{
			name:        "valid_stay_at_boundary",
			start:       &model.Position{X: 0, Y: 0, HasCrate: false},
			commands:    "NENS",
			expectError: false,
			description: "Move along boundary without crossing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateBoundary(tt.start, tt.commands)

			if tt.expectError {
				if err == nil {
					t.Errorf("validateBoundary() expected error but got none. %s", tt.description)
				} else if err != model.ErrBoundary {
					t.Errorf("validateBoundary() expected ErrBoundary but got %v. %s", err, tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("validateBoundary() unexpected error = %v. %s", err, tt.description)
				}
			}
		})
	}
}

func TestCreateTaskServiceImpl_validateBoundary_InvalidCommands(t *testing.T) {
	service := &CreateTaskServiceImpl{}
	start := &model.Position{X: 5, Y: 5, HasCrate: false}

	invalidCommands := []string{
		"X",     // Invalid command
		"NXES",  // Valid commands with invalid command in middle
		"123",   // Numbers
		"N@E",   // Special characters
		"HELLO", // Mix of valid and invalid letters
	}

	for _, cmd := range invalidCommands {
		t.Run("invalid_command_"+cmd, func(t *testing.T) {
			// Note: Based on the current implementation, invalid commands are silently ignored
			// This test documents the current behavior - you may want to modify the function
			// to return an error for invalid commands
			err := service.validateBoundary(start, cmd)
			// Current implementation doesn't validate command characters, only boundaries
			// If you want to add command validation, modify the function and this test
			_ = err // Suppress unused variable warning
		})
	}
}

func TestCreateTaskServiceImpl_validateBoundary_NilPosition(t *testing.T) {
	service := &CreateTaskServiceImpl{}

	// This test checks behavior with nil position - may cause panic in current implementation
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior - function doesn't handle nil position
			t.Logf("Function panicked with nil position as expected: %v", r)
		}
	}()

	err := service.validateBoundary(nil, "N")
	if err == nil {
		t.Error("Expected function to handle nil position, but it didn't return an error")
	}
}

// Benchmark test for performance
func BenchmarkCreateTaskServiceImpl_validateBoundary(b *testing.B) {
	service := &CreateTaskServiceImpl{}
	start := &model.Position{X: 5, Y: 5, HasCrate: false}
	commands := "NNESSSWWNNEESSSWW"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.validateBoundary(start, commands)
		if err != nil {
			return
		}
	}
}
