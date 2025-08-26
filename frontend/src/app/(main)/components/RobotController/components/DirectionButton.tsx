import React from "react";
import { DirectionEnum } from "../../../types/types";

interface DirectionButtonProps {
  direction: DirectionEnum;
  onClick: (direction: DirectionEnum) => void;
}

const DirectionButton: React.FC<DirectionButtonProps> = ({
  direction,
  onClick,
}) => {
  const colors: Record<DirectionEnum, string> = {
    N: "bg-blue-600 hover:bg-blue-700",
    S: "bg-red-600 hover:bg-red-700",
    E: "bg-amber-600 hover:bg-amber-700",
    W: "bg-green-600 hover:bg-green-700",
  };

  return (
    <button
      onClick={() => onClick(direction)}
      className={`w-20 h-20 ${colors[direction]} text-white font-bold text-2xl rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200 active:scale-95`}
    >
      {direction}
    </button>
  );
};

export default DirectionButton;
