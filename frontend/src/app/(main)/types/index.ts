import { TaskStatus } from "@/app/_common/redux/robotwarehouseApi/robotwarehouseApi.type";

export interface Position {
  x: number;
  y: number;
}

export interface TaskData {
  position: Position;
  status: TaskStatus;
}

export enum DirectionEnum {
  NORTH = "N",
  SOUTH = "S",
  EAST = "E",
  WEST = "W",
}
