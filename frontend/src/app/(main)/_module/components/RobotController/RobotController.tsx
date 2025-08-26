"use client";

import React, { useState } from "react";

const ConsoleController = () => {
  const [inputValue, setInputValue] = useState("");

  const handleDirectionClick = (direction: string) => {
    const command = `MOVE_${direction}`;
    setInputValue((prev) => prev + (prev ? ", " : "") + command);
  };

  const clearInput = () => {
    setInputValue("");
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-900 p-8">
      <div className="bg-gray-800 rounded-2xl p-8 shadow-2xl border border-gray-700">
        <h1 className="text-2xl font-bold text-white text-center mb-8">
          Console Controller
        </h1>

        {/* Input Display */}
        <div className="mb-8">
          <label className="block text-gray-300 text-sm font-medium mb-2">
            Command Input:
          </label>
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Commands will appear here..."
          />
          <button
            onClick={clearInput}
            className="mt-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition-colors"
          >
            Clear
          </button>
        </div>

        {/* Controller Layout */}
        <div className="relative w-64 h-64 mx-auto">
          {/* North Button */}
          <button
            onClick={() => handleDirectionClick("NORTH")}
            className="absolute top-0 left-1/2 transform -translate-x-1/2 w-16 h-16 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-lg shadow-lg transition-all duration-150 flex items-center justify-center group"
          >
            <span className="text-lg">N</span>
            <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-black text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity">
              North
            </div>
          </button>

          {/* West Button */}
          <button
            onClick={() => handleDirectionClick("WEST")}
            className="absolute left-0 top-1/2 transform -translate-y-1/2 w-16 h-16 bg-green-600 hover:bg-green-700 active:bg-green-800 text-white font-bold rounded-lg shadow-lg transition-all duration-150 flex items-center justify-center group"
          >
            <span className="text-lg">W</span>
            <div className="absolute -left-8 top-1/2 transform -translate-y-1/2 bg-black text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity">
              West
            </div>
          </button>

          {/* Center Circle (decorative) */}
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-20 h-20 bg-gray-700 border-4 border-gray-600 rounded-full flex items-center justify-center">
            <div className="w-8 h-8 bg-gray-600 rounded-full"></div>
          </div>

          {/* East Button */}
          <button
            onClick={() => handleDirectionClick("EAST")}
            className="absolute right-0 top-1/2 transform -translate-y-1/2 w-16 h-16 bg-yellow-600 hover:bg-yellow-700 active:bg-yellow-800 text-white font-bold rounded-lg shadow-lg transition-all duration-150 flex items-center justify-center group"
          >
            <span className="text-lg">E</span>
            <div className="absolute -right-8 top-1/2 transform -translate-y-1/2 bg-black text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity">
              East
            </div>
          </button>

          {/* South Button */}
          <button
            onClick={() => handleDirectionClick("SOUTH")}
            className="absolute bottom-0 left-1/2 transform -translate-x-1/2 w-16 h-16 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white font-bold rounded-lg shadow-lg transition-all duration-150 flex items-center justify-center group"
          >
            <span className="text-lg">S</span>
            <div className="absolute -bottom-8 left-1/2 transform -translate-x-1/2 bg-black text-white text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity">
              South
            </div>
          </button>
        </div>

        {/* Command History */}
        {/* {commandHistory.length > 0 && (
          <div className="mt-8">
            <h3 className="text-gray-300 text-sm font-medium mb-2">
              Command History:
            </h3>
            <div className="bg-gray-700 rounded-lg p-3 max-h-24 overflow-y-auto">
              <div className="text-gray-300 text-sm font-mono">
                {commandHistory.join(' → ')}
              </div>
            </div>
          </div>
        )} */}
      </div>
    </div>
  );
};

export default ConsoleController;
