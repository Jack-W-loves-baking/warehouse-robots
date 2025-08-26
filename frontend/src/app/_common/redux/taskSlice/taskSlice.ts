import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface TaskSlice {
  tasksList: Task[];
}

interface Task {
  taskId: string;
  robotId: string;
  commands: string;
  status: string;
}
const initialState: TaskSlice = {
  tasksList: [],
};

export const taskSlice = createSlice({
  name: "task",
  initialState,
  reducers: {
    updatedList: (state, action: PayloadAction<Task>) => {
      state.tasksList = [...state.tasksList, action.payload];
    },
  },
});
