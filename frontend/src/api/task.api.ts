import { apiClient } from './client';
import type { CreateTaskRequest, UpdateTaskRequest, Task, TaskResponse } from '@/types/task.types';

export const taskApi = {
  createTask: async (data: CreateTaskRequest): Promise<Task> => {
    const response = await apiClient.post<Task>('/tasks', data);
    return response.data;
  },

  getUserTasks: async (params?: {
    page?: number;
    limit?: number;
    status?: string;
    priority?: number;
    project_id?: string;
  }): Promise<TaskResponse> => {
    const response = await apiClient.get<TaskResponse>('/tasks', { params });
    return response.data;
  },

  getTask: async (id: string): Promise<Task> => {
    const response = await apiClient.get<Task>(`/tasks/${id}`);
    return response.data;
  },

  updateTask: async (id: string, data: UpdateTaskRequest): Promise<Task> => {
    const response = await apiClient.put<Task>(`/tasks/${id}`, data);
    return response.data;
  },

  deleteTask: async (id: string): Promise<void> => {
    await apiClient.delete(`/tasks/${id}`);
  },
};
