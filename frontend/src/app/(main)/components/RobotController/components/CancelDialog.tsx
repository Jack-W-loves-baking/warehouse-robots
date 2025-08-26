import React, { useState } from "react";

interface CancelDialogProps {
  onConfirm: () => void;
}

const CancelDialog: React.FC<CancelDialogProps> = ({ onConfirm }) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleConfirm = () => {
    onConfirm();
    setIsOpen(false);
  };

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="px-8 py-3 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200"
      >
        Cancel Task
      </button>

      {isOpen && (
        <>
          {/* Backdrop */}
          <div
            className="fixed inset-0 bg-black/50 z-40"
            onClick={() => setIsOpen(false)}
          />

          {/* Modal */}
          <div className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-full max-w-md p-6 bg-slate-800 rounded-xl shadow-2xl border border-slate-700">
            <div className="mb-4">
              <h2 className="text-xl font-bold text-white mb-2">
                Are you sure?
              </h2>
              <p className="text-slate-400">
                This will cancel the current task. The robot will stop moving
                and the task cannot be resumed.
              </p>
            </div>

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setIsOpen(false)}
                className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white font-medium rounded-lg transition-colors"
              >
                No, keep running
              </button>
              <button
                onClick={handleConfirm}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors"
              >
                Yes, cancel task
              </button>
            </div>
          </div>
        </>
      )}
    </>
  );
};

export default CancelDialog;
