interface CreateTaskResponse {
  task_id: string;
  robot_id: string;
  commands: string;
  status: string;
  currentState: {
    x: number;
    y: number;
    has_crate: boolean;
  };
  error: string;
  create_at: string;
  update_at: string;
}

interface CreateTaskRequest {
  commands: string;
}
