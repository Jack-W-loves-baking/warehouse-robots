import React from "react";

const RobotIcon: React.FC = () => (
  <div className="absolute inset-0 flex items-center justify-center animate-pulse">
    <div className="relative">
      <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg shadow-lg">
        <div className="flex justify-center pt-1.5 gap-1">
          <div className="w-1.5 h-1.5 bg-white rounded-full"></div>
          <div className="w-1.5 h-1.5 bg-white rounded-full"></div>
        </div>
        <div className="w-4 h-0.5 bg-white/70 mx-auto mt-1 rounded"></div>
      </div>
      <div className="absolute -bottom-1 left-1/2 -translate-x-1/2 w-6 h-1 bg-black/30 rounded-full blur-sm"></div>
    </div>
  </div>
);

export default RobotIcon;
