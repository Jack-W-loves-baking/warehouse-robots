import React from "react";
import RobotIcon from "./RobotIcon";

// Grid Cell Component
interface GridCellProps {
  x: number;
  y: number;
  hasRobot: boolean;
  showCoordinates: boolean;
}

const GridCell: React.FC<GridCellProps> = ({
  x,
  y,
  hasRobot,
  showCoordinates,
}) => (
  <div className="relative flex items-center justify-center border border-black bg-white transition-all duration-200">
    {hasRobot && <RobotIcon />}
    {showCoordinates && (
      <span className="absolute bottom-0.5 right-1 text-xs text-black opacity-0 group-hover:opacity-100 transition-opacity">
        {x},{y}
      </span>
    )}
  </div>
);

export default GridCell;
