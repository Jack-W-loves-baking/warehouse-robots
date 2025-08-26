import React from "react";

import { TaskStatus } from "@/app/_common/redux/robotwarehouseApi/robotwarehouseApi.type";
import GridCell from "./components/GridCell";
import { Position } from "../../types";

interface RobotGridProps {
  robotPosition: Position;
  taskStatus: TaskStatus;
}

const RobotGrid: React.FC<RobotGridProps> = ({ robotPosition, taskStatus }) => {
  const cells: any[] = [];

  for (let y = 9; y >= 0; y--) {
    for (let x = 0; x < 10; x++) {
      cells.push(
        <GridCell
          key={`${x}-${y}`}
          x={x}
          y={y}
          hasRobot={robotPosition.x === x && robotPosition.y === y}
          showCoordinates={true}
        />
      );
    }
  }

  return (
    <div className="bg-white border border-black rounded-lg p-6 shadow-md">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-2xl font-bold text-black">Robot Grid</h2>
          <p className="text-gray-700 text-sm mt-1">
            Position: ({robotPosition.x}, {robotPosition.y})
          </p>
        </div>
      </div>

      <div className="flex-[2] flex items-center justify-center p-6">
        <div className="w-full h-full max-w-5xl aspect-square">
          <div className="grid grid-cols-10 gap-1 w-full h-full">{cells}</div>
        </div>
      </div>
    </div>
  );
};

export default RobotGrid;
