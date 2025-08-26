package model

import (
	"errors"
	"warehouse-robots/backend/api/constant"
)

var (
	ErrValidation        = errors.New(constant.ErrorCodeValidation)
	ErrBoundary          = errors.New(constant.ErrorCodeBoundary)
	ErrRobotIDInvalid    = errors.New(constant.ErrorCodeRobotIdInvalid)
	ErrRobotNotFound     = errors.New(constant.ErrorCodeRobotNotFound)
	ErrTaskNotFound      = errors.New(constant.ErrorCodeTaskNotFound)
	ErrRobotBusy         = errors.New(constant.ErrorCodeRobotBusy)
	ErrTaskQueueFull     = errors.New(constant.ErrorCodeTaskQueueFull)
	ErrInternal          = errors.New(constant.ErrorCodeInternal)
	ErrTaskProcessed     = errors.New(constant.ErrorCodeTaskAlreadyDone)
	ErrSDKFailedToCancel = errors.New(constant.ErrorSDKFailedToCancel)
)
