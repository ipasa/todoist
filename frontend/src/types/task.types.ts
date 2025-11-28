export interface CreateTaskRequest {
  title: string;
  description?: string;
  priority?: number;
  project_id?: string;
  due_date?: string;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  status?: string;
  priority?: number;
  project_id?: string;
  due_date?: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  status: string;
  priority: number;
  user_id: string;
  project_id?: string;
  due_date?: string;
  created_at: string;
  updated_at: string;
}

export interface TaskResponse {
  data: Task[];
  total: number;
  page: number;
  limit: number;
}
