package service

// ICancelTaskService LOGIC:
//   - If the task is already terminal (COMPLETED, FAILED, or CANCELLED): reject.
//   - If the task is PENDING or RUNNING: invoke the SDK's CancelTask with up to
//     three retries using backoff.
//   - On success: stop monitoring and persist status = CANCELLED.
//   - On repeated failure: we dont do anything.
type ICancelTaskService interface {
	// CancelTaskById cancels the task in both sdk and updates status in the db
	//
	//Error returns
	//
	//	 - ErrTaskNotFound - when repository cannot find the task id
	//	 - ErrTaskProcessed - when task status is COMPLETED, FAILED or CANCELLED
	//	 - ErrRobotNotFound - when robot is not found.
	//	 - ErrSDKFailedToCancel - when sdk failed to cancel task even after retry
	CancelTaskById(taskID string) error
}
