"use client";

import React, { useState, useEffect } from "react";

import RobotController from "./components/RobotController";
import { TaskStatus } from "../_common/redux/robotwarehouseApi/robotwarehouseApi.type";
import {
  useCreateTaskMutation,
  useDeleteTaskMutation,
  useGetTaskQuery,
} from "../_common/redux/robotwarehouseApi";
import RobotGrid from "./components/RobotGrid/RobotGrid";
import { Position } from "./types";
import { ROBOT_ID } from "./constants";

export default function Home() {
  const [commands, setCommands] = useState<string>("");
  const [currentTaskId, setCurrentTaskId] = useState<string>("");
  const [robotPosition, setRobotPosition] = useState<Position>({ x: 0, y: 0 });
  const [taskStatus, setTaskStatus] = useState<TaskStatus>("IDLE");

  const [pollingInterval, setPollingInterval] = useState(1000);

  const [createTask, { isLoading: isCreating }] = useCreateTaskMutation();
  const [deleteTask, { isLoading: isDeleting }] = useDeleteTaskMutation();

  // Poll for task status when we have an active task
  const { data: taskData } = useGetTaskQuery(currentTaskId, {
    pollingInterval,
    skip: !currentTaskId,
  });

  useEffect(() => {
    if (currentTaskId) {
      setPollingInterval(1000);

      if (
        taskData?.status &&
        ["COMPLETED", "FAILED", "CANCELLED"].includes(taskData?.status)
      ) {
        setPollingInterval(0);
      }
    }
  }, [taskData, currentTaskId]);

  useEffect(() => {
    if (taskData) {
      setRobotPosition({
        x: taskData.current_state.x,
        y: taskData.current_state.y,
      });

      if (taskData.status === "COMPLETED") {
        setTaskStatus("COMPLETED");
        setCurrentTaskId("");
      }
    }
  }, [taskData]);

  const handleRunTask = async () => {
    if (!commands) return;

    try {
      setTaskStatus("PENDING");
      const result = await createTask({
        robotId: ROBOT_ID,
        commands,
      }).unwrap();

      if (result?.task_id) {
        setCurrentTaskId(result.task_id);
      }
    } catch (error) {
      alert("Out of boundary");
      setTaskStatus("IDLE");
    }
  };

  const handleCancelTask = async () => {
    if (!currentTaskId) return;

    try {
      await deleteTask(currentTaskId as any).unwrap();
      alert("task deleted");
      setTaskStatus("CANCELLED");
      setCurrentTaskId("");
      setTimeout(() => setTaskStatus("IDLE"), 3000);
    } catch (error) {
      alert("failed to delete task");
      console.error("Failed to cancel task:", error);
    }
  };

  return (
    <div className="min-h-screen bg-white relative">
      <div className="w-full h-screen flex items-center justify-center p-8">
        <div style={{ width: 600, height: 600 }}>
          <RobotGrid robotPosition={robotPosition} taskStatus={taskStatus} />
        </div>
      </div>

      <RobotController
        commands={commands}
        onCommandAdd={setCommands}
        onRun={handleRunTask}
        onCancel={handleCancelTask}
        isTaskRunning={taskStatus === "PENDING"}
        isLoading={isCreating || isDeleting}
      />
    </div>
  );
}
