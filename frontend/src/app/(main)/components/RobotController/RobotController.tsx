"use client";

import React, { useState } from "react";
import { motion } from "framer-motion";

import DirectionButton from "./components/DirectionButton";
import CommandInput from "./components/CommandInput";
import ControlButtons from "./components/ControlButtons";
import { DirectionEnum } from "../../types";

interface RobotControllerProps {
  commands: string;
  onCommandAdd: (commands: string) => void;
  onRun: () => void;
  onCancel: () => void;
  isTaskRunning: boolean;
  isLoading?: boolean;
}

const RobotController = ({
  onCommandAdd,
  onRun,
  onCancel,
  commands,
  isTaskRunning,
  isLoading,
}: RobotControllerProps) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleDirectionClick = (direction: DirectionEnum) => {
    onCommandAdd(commands + direction);
  };

  const handleClear = () => {
    onCommandAdd("");
  };

  const togglePanel = () => {
    setIsOpen(!isOpen);
  };

  return (
    <motion.div
      initial={{ x: "85%" }}
      animate={{ x: isOpen ? 0 : "85%" }}
      transition={{
        type: "spring",
        stiffness: 300,
        damping: 30,
      }}
      className="fixed right-0 top-0 h-full w-96 bg-gradient-to-t from-slate-900 via-slate-800 to-slate-800/90 backdrop-blur-xl border-l border-slate-700 shadow-2xl z-50"
    >
      <button
        onClick={togglePanel}
        className="absolute left-0 top-1/2 -translate-x-full -translate-y-1/2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white p-3 rounded-l-lg shadow-lg transition-all duration-200 hover:scale-105"
      >
        <motion.div
          animate={{ rotate: isOpen ? 0 : 180 }}
          transition={{ duration: 0.2 }}
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
        </motion.div>
      </button>
      <div className="h-full overflow-y-auto">
        <div className="px-8 py-8">
          <motion.h1
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: isOpen ? 1 : 0.7, y: 0 }}
            transition={{ delay: 0.1 }}
            className="text-3xl font-bold text-white text-center mb-6"
          >
            Robot Controller
          </motion.h1>

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: isOpen ? 1 : 0.8 }}
            transition={{ delay: 0.2 }}
          >
            <CommandInput value={commands} onChange={onCommandAdd} />
          </motion.div>

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: isOpen ? 1 : 0.8 }}
            transition={{ delay: 0.3 }}
            className="flex flex-col justify-center gap-4 w-full"
          >
            <ControlButtons
              onRun={onRun}
              onCancel={onCancel}
              onClear={handleClear}
              isRunning={isTaskRunning}
              isLoading={isLoading}
            />
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{
              opacity: isOpen ? 1 : 0.7,
              scale: isOpen ? 1 : 0.9,
            }}
            transition={{ delay: 0.4 }}
            className="mt-8 flex justify-center"
          >
            <div className="grid grid-cols-3 gap-3 w-72">
              <div></div>
              <DirectionButton
                direction={DirectionEnum.NORTH}
                onClick={handleDirectionClick}
              />
              <div></div>

              <DirectionButton
                direction={DirectionEnum.WEST}
                onClick={handleDirectionClick}
              />
              <div className="w-20 h-20 bg-slate-700/50 rounded-full flex items-center justify-center">
                <div className="w-10 h-10 bg-slate-600 rounded-full shadow-inner"></div>
              </div>
              <DirectionButton
                direction={DirectionEnum.EAST}
                onClick={handleDirectionClick}
              />

              <div></div>
              <DirectionButton
                direction={DirectionEnum.SOUTH}
                onClick={handleDirectionClick}
              />
              <div></div>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: isOpen ? 1 : 0.8 }}
            transition={{ delay: 0.5 }}
            className="flex justify-center mt-6"
          ></motion.div>
        </div>
      </div>
    </motion.div>
  );
};

export default RobotController;
