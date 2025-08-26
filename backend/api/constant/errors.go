package constant

// Error codes
const (
	// Request/validation
	ErrorCodeValidation     = "VALIDATION_ERROR"
	ErrorCodeBoundary       = "BOUNDARY_ERROR"
	ErrorCodeRobotIdInvalid = "ROBOT_ID_INVALID"

	// Lookup
	ErrorCodeTaskNotFound  = "TASK_NOT_FOUND"
	ErrorCodeRobotNotFound = "ROBOT_NOT_FOUND"

	// State
	ErrorCodeRobotBusy       = "ROBOT_BUSY"
	ErrorCodeTaskAlreadyDone = "TASK_ALREADY_TERMINAL"

	// Queue/capacity
	ErrorCodeTaskQueueFull = "TASK_QUEUE_FULL"

	// Sdk
	ErrorSDKFailedToCancel = "SDK_CANCEL_FAILED"

	// General
	ErrorCodeInternal = "INTERNAL_ERROR"
)
