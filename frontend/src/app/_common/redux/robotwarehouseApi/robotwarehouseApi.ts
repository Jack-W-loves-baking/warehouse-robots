import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const robotwarehouseApi = createApi({
  reducerPath: "warehouseRobot",
  baseQuery: fetchBaseQuery({
    baseUrl: `/api/`,
  }),
  endpoints: (builder) => ({
    createTask: builder.mutation<CreateTaskResponse, CreateTaskRequest>({
      query: (robotId) => ({
        url: `${robotId}/tasks`,
        method: "POST",
      }),
    }),
    deleteTask: builder.mutation<void, void>({
      query: (taskId) => ({
        url: `tasks/${taskId}`,
        method: "DELETE",
      }),
    }),
    getTask: builder.query<CreateTaskResponse, void>({
      query: (taskId) => ({
        url: `tasks/${taskId}`,
      }),
    }),
  }),
});

export const { useCreateTaskMutation, useGetTaskQuery, useDeleteTaskMutation } =
  robotwarehouseApi;
