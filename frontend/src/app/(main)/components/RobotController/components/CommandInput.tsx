import React from "react";

interface CommandInputProps {
  value: string;
  onChange: (value: string) => void;
}

const CommandInput: React.FC<CommandInputProps> = ({ value, onChange }) => (
  <div className="mb-6">
    <label className="text-slate-400 text-sm font-medium mb-2 block">
      Command
    </label>
    <div className="bg-slate-900/50 border-2 border-slate-700 rounded-xl p-4 min-h-16 flex items-center">
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value.toUpperCase())}
        placeholder="Enter commands (NSEW)"
        className="bg-transparent text-slate-300 font-mono text-lg outline-none w-full placeholder-slate-600"
      />
    </div>
  </div>
);

export default CommandInput;
