export interface CreateTaskResponse {
  task_id: string;
  robot_id: string;
  commands: string;
  status: TaskStatus;
  current_state: {
    x: number;
    y: number;
    has_crate: boolean;
  };
  error: string;
  create_at: string;
  update_at: string;
}

export interface CreateTaskRequest {
  commands: string;
  robotId: string;
}

export type TaskStatus =
  | "PENDING"
  | "COMPLETED"
  | "FAILED"
  | "CANCELLED"
  | "IDLE";
