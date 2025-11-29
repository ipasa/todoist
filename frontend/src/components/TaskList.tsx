import { useState, useEffect } from "react";
import type { Task } from "@/types/task.types";
import { taskApi } from "@/api/task.api";

interface TaskListProps {
  refreshTrigger?: number;
}

export function TaskList({ refreshTrigger }: TaskListProps) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchTasks();
  }, [refreshTrigger]);

  const fetchTasks = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await taskApi.getUserTasks();
      setTasks(response.data || []);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to fetch tasks");
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteTask = async (id: string) => {
    if (!window.confirm("Are you sure you want to delete this task?")) {
      return;
    }

    try {
      await taskApi.deleteTask(id);
      setTasks(tasks.filter((task) => task.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to delete task");
    }
  };

  const getPriorityLabel = (priority: number) => {
    switch (priority) {
      case 1:
        return "Low";
      case 2:
        return "Medium";
      case 3:
        return "High";
      case 4:
        return "Urgent";
      default:
        return "Unknown";
    }
  };

  const getPriorityColor = (priority: number) => {
    switch (priority) {
      case 1:
        return "bg-gray-100 text-gray-800";
      case 2:
        return "bg-blue-100 text-blue-800";
      case 3:
        return "bg-orange-100 text-orange-800";
      case 4:
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "pending":
        return "bg-gray-100 text-gray-800";
      case "in_progress":
        return "bg-blue-100 text-blue-800";
      case "completed":
        return "bg-green-100 text-green-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <p className="text-sm text-red-800">{error}</p>
          <button
            onClick={fetchTasks}
            className="mt-2 text-sm font-medium text-red-600 hover:text-red-500"
          >
            Try again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-4">Your Tasks</h2>

      {tasks.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500">No tasks yet. Create your first task!</p>
        </div>
      ) : (
        <div className="overflow-hidden">
          <ul className="divide-y divide-gray-200">
            {tasks.map((task) => (
              <li key={task.id} className="py-4">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {task.title}
                    </p>
                    <p className="text-sm text-gray-500 truncate">
                      {task.description}
                    </p>
                    <div className="mt-2 flex items-center space-x-2">
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityColor(task.priority)}`}
                      >
                        {getPriorityLabel(task.priority)}
                      </span>
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(task.status)}`}
                      >
                        {task.status.replace("_", " ")}
                      </span>
                      {task.due_date && (
                        <span className="text-xs text-gray-500">
                          Due: {formatDate(task.due_date)}
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="flex-shrink-0 ml-4">
                    <button
                      onClick={() => handleDeleteTask(task.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
