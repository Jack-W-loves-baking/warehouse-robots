import react from "react";
import CancelDialog from "./CancelDialog";
import { motion } from "framer-motion";

// Control Buttons Component
interface ControlButtonsProps {
  onRun: () => void;
  onCancel: () => void;
  isRunning: boolean;
  onClear: () => void;
  isLoading?: boolean;
}

const ControlButtons: React.FC<ControlButtonsProps> = ({
  onRun,
  onCancel,
  isRunning,
  onClear,
  isLoading = false,
}) => (
  <>
    <div className="flex gap-4 justify-center w-full flex-col">
      <button
        onClick={onRun}
        disabled={isRunning || isLoading}
        className="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-700 disabled:cursor-not-allowed text-white font-semibold rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200 "
      >
        Run Command
      </button>

      <button
        onClick={onClear}
        className="px-8 py-3 bg-gray-200 hover:bg-gray-300 text-black font-semibold rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200"
      >
        Clear Commands
      </button>
      <CancelDialog onConfirm={onCancel} />
    </div>
  </>
);

export default ControlButtons;
