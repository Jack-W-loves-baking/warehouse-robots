import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import {
  CreateTaskRequest,
  CreateTaskResponse,
} from "./robotwarehouseApi.type";

export const robotwarehouseApi = createApi({
  reducerPath: "warehouseRobot",
  baseQuery: fetchBaseQuery({
    baseUrl: `/api/`,
  }),
  endpoints: (builder) => ({
    createTask: builder.mutation<CreateTaskResponse, CreateTaskRequest>({
      query: ({ commands, robotId }) => ({
        url: `robots/${robotId}/tasks`,
        method: "POST",
        body: {
          commands,
        },
      }),
    }),
    deleteTask: builder.mutation<void, string>({
      query: (taskId) => ({
        url: `tasks/${taskId}`,
        method: "DELETE",
      }),
    }),
    getTask: builder.query<CreateTaskResponse, string>({
      query: (taskId) => ({
        url: `tasks/${taskId}`,
      }),
    }),
  }),
});

export const { useCreateTaskMutation, useGetTaskQuery, useDeleteTaskMutation } =
  robotwarehouseApi;
