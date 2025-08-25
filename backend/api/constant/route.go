package constant

// API route constants
const (
	RouteGetRobots      = "GET /api/v1/robots"
	RouteGetRobotById   = "GET /api/v1/robots/{robotId}"
	RouteCreateTask     = "POST /api/v1/robots/{robotId}/tasks"
	RouteGetTaskById    = "GET /api/v1/tasks/{taskId}"
	RouteDeleteTaskById = "DELETE /api/v1/tasks/{taskId}"
)
