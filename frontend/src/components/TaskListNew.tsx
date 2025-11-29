import { useState, useEffect } from "react";
import type { Task } from "@/types/task.types";
import { taskApi } from "@/api/task.api";

interface TaskListNewProps {
  refreshTrigger?: number;
  onTaskCreated?: () => void;
  onTaskClick?: (task: Task) => void;
  onTaskCountChange?: (count: number) => void;
}

export function TaskListNew({ refreshTrigger, onTaskCreated, onTaskClick, onTaskCountChange }: TaskListNewProps) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newTaskTitle, setNewTaskTitle] = useState("");
  const [newTaskDescription, setNewTaskDescription] = useState("");
  const [newTaskDueDate, setNewTaskDueDate] = useState("");
  const [newTaskPriority, setNewTaskPriority] = useState(1);
  const [isAdding, setIsAdding] = useState(false);
  const [showAddForm, setShowAddForm] = useState(false);
  const [showDescription, setShowDescription] = useState(false);
  const [completingTasks, setCompletingTasks] = useState<Set<string>>(new Set());

  useEffect(() => {
    fetchTasks();
  }, [refreshTrigger]);

  const fetchTasks = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await taskApi.getUserTasks();
      const taskList = response.data || [];
      setTasks(taskList);
      onTaskCountChange?.(taskList.filter(t => t.status !== "completed").length);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to fetch tasks");
    } finally {
      setLoading(false);
    }
  };

  const handleQuickAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTaskTitle.trim()) return;

    setIsAdding(true);
    try {
      // Convert due_date to RFC3339 format if provided
      let dueDate = undefined;
      if (newTaskDueDate) {
        const date = new Date(newTaskDueDate);
        dueDate = date.toISOString();
      }

      const newTask = await taskApi.createTask({
        title: newTaskTitle,
        description: newTaskDescription,
        priority: newTaskPriority,
        due_date: dueDate,
      });
      setTasks([newTask, ...tasks]);

      // Reset form
      setNewTaskTitle("");
      setNewTaskDescription("");
      setNewTaskDueDate("");
      setNewTaskPriority(1);
      setShowAddForm(false);
      setShowDescription(false);

      onTaskCreated?.();
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to create task");
    } finally {
      setIsAdding(false);
    }
  };

  const handleToggleComplete = async (task: Task) => {
    const newStatus = task.status === "completed" ? "pending" : "completed";

    // Add to completing set for animation
    setCompletingTasks(prev => new Set(prev).add(task.id));

    try {
      await taskApi.updateTask(task.id, { status: newStatus });

      // Update local state
      setTasks(tasks.map(t =>
        t.id === task.id ? { ...t, status: newStatus } : t
      ));

      // Remove from completing set after animation
      setTimeout(() => {
        setCompletingTasks(prev => {
          const next = new Set(prev);
          next.delete(task.id);
          return next;
        });
      }, 300);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to update task");
      setCompletingTasks(prev => {
        const next = new Set(prev);
        next.delete(task.id);
        return next;
      });
    }
  };

  const handleDeleteTask = async (id: string) => {
    if (!window.confirm("Delete this task?")) return;

    try {
      await taskApi.deleteTask(id);
      setTasks(tasks.filter((task) => task.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to delete task");
    }
  };

  const getPriorityColor = (priority: number) => {
    switch (priority) {
      case 4:
        return "text-red-500 hover:text-red-600";
      case 3:
        return "text-orange-500 hover:text-orange-600";
      case 2:
        return "text-blue-500 hover:text-blue-600";
      default:
        return "text-gray-400 hover:text-gray-500";
    }
  };

  const getPriorityColorBg = (priority: number) => {
    switch (priority) {
      case 4:
        return "text-red-600 bg-red-50";
      case 3:
        return "text-orange-600 bg-orange-50";
      case 2:
        return "text-blue-600 bg-blue-50";
      default:
        return "text-gray-500";
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-sm text-red-800">{error}</p>
        <button
          onClick={fetchTasks}
          className="mt-2 text-sm font-medium text-red-600 hover:text-red-500"
        >
          Try again
        </button>
      </div>
    );
  }

  const activeTasks = tasks.filter(t => t.status !== "completed");
  const completedTasks = tasks.filter(t => t.status === "completed");

  return (
    <div className="space-y-4">
      {/* Quick Add Button/Form */}
      {!showAddForm ? (
        <button
          onClick={() => setShowAddForm(true)}
          className="w-full text-left px-4 py-3 text-gray-400 hover:text-gray-600 hover:bg-gray-50 rounded-lg border-2 border-transparent hover:border-gray-200 transition-all group"
        >
          <div className="flex items-center space-x-3">
            <div className="w-5 h-5 rounded-full border-2 border-gray-300 group-hover:border-red-500 transition-colors" />
            <span className="text-sm">Add task</span>
          </div>
        </button>
      ) : (
        <form onSubmit={handleQuickAdd} className="bg-white border-2 border-gray-200 rounded-lg p-4 shadow-sm">
          <input
            type="text"
            value={newTaskTitle}
            onChange={(e) => setNewTaskTitle(e.target.value)}
            placeholder="Task name"
            autoFocus
            className="w-full text-sm border-none focus:outline-none focus:ring-0 p-0 mb-2"
          />

          {showDescription && (
            <textarea
              value={newTaskDescription}
              onChange={(e) => setNewTaskDescription(e.target.value)}
              placeholder="Description"
              rows={2}
              className="w-full text-xs text-gray-600 border-none focus:outline-none focus:ring-0 p-0 mb-2 resize-none"
            />
          )}

          <div className="flex items-center justify-between mt-3">
            <div className="flex items-center space-x-1">
              {/* Due Date Picker */}
              <label className={`text-xs px-2 py-1 hover:bg-gray-100 rounded flex items-center space-x-1 cursor-pointer ${
                newTaskDueDate ? "text-green-600 bg-green-50" : "text-gray-500"
              }`}>
                <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <span>
                  {newTaskDueDate
                    ? new Date(newTaskDueDate + 'T00:00:00').toLocaleDateString("en-US", { month: "short", day: "numeric" })
                    : "Due date"
                  }
                </span>
                <input
                  type="date"
                  value={newTaskDueDate}
                  onChange={(e) => setNewTaskDueDate(e.target.value)}
                  className="absolute opacity-0 w-0 h-0"
                />
              </label>

              {/* Priority Selector */}
              <label className={`text-xs px-2 py-1 hover:bg-gray-100 rounded flex items-center space-x-1 cursor-pointer ${
                newTaskPriority > 1 ? getPriorityColorBg(newTaskPriority) : "text-gray-500"
              }`}>
                <svg className="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M3 6l7-4 7 4v8l-7 4-7-4V6z" />
                </svg>
                <span>P{5 - newTaskPriority}</span>
                <select
                  value={newTaskPriority}
                  onChange={(e) => setNewTaskPriority(Number(e.target.value))}
                  className="absolute opacity-0 w-0 h-0"
                >
                  <option value="1">P4</option>
                  <option value="2">P3</option>
                  <option value="3">P2</option>
                  <option value="4">P1</option>
                </select>
              </label>

              {/* Description Toggle */}
              {!showDescription && (
                <button
                  type="button"
                  onClick={() => setShowDescription(true)}
                  className="text-xs text-gray-500 hover:text-gray-700 px-2 py-1 hover:bg-gray-100 rounded flex items-center space-x-1"
                >
                  <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16m-7 6h7" />
                  </svg>
                  <span>Desc</span>
                </button>
              )}
            </div>

            <div className="flex items-center space-x-2">
              <button
                type="button"
                onClick={() => {
                  setShowAddForm(false);
                  setNewTaskTitle("");
                  setNewTaskDescription("");
                  setNewTaskDueDate("");
                  setNewTaskPriority(1);
                  setShowDescription(false);
                }}
                className="text-xs px-3 py-1.5 text-gray-600 hover:bg-gray-100 rounded"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isAdding || !newTaskTitle.trim()}
                className="text-xs px-3 py-1.5 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isAdding ? "Adding..." : "Add task"}
              </button>
            </div>
          </div>
        </form>
      )}

      {/* Active Tasks */}
      <div className="space-y-1">
        {activeTasks.length === 0 && !showAddForm && (
          <div className="text-center py-12">
            <p className="text-gray-400 text-sm">No tasks. Enjoy your day!</p>
          </div>
        )}

        {activeTasks.map((task) => (
          <div
            key={task.id}
            className={`group bg-white hover:bg-gray-50 border border-gray-200 rounded-lg p-4 transition-all ${
              completingTasks.has(task.id) ? "opacity-50 scale-98" : ""
            }`}
          >
            <div className="flex items-start space-x-3">
              {/* Checkbox */}
              <button
                onClick={() => handleToggleComplete(task)}
                className={`mt-0.5 flex-shrink-0 w-5 h-5 rounded-full border-2 ${
                  task.status === "completed"
                    ? "bg-gray-400 border-gray-400"
                    : "border-gray-300 hover:border-gray-400"
                } transition-all flex items-center justify-center`}
              >
                {task.status === "completed" && (
                  <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                  </svg>
                )}
              </button>

              {/* Task Content */}
              <div
                className="flex-1 min-w-0 cursor-pointer"
                onClick={() => onTaskClick?.(task)}
              >
                <p className={`text-sm ${task.status === "completed" ? "line-through text-gray-400" : "text-gray-900"}`}>
                  {task.title}
                </p>
                {task.description && (
                  <p className="text-xs text-gray-500 mt-1">{task.description}</p>
                )}
                <div className="flex items-center space-x-2 mt-2">
                  {task.due_date && (
                    <span className="text-xs text-gray-500">
                      {new Date(task.due_date).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                    </span>
                  )}
                  {task.priority > 1 && (
                    <svg className={`w-4 h-4 ${getPriorityColor(task.priority)}`} fill="currentColor" viewBox="0 0 20 20">
                      <path d="M3 6l7-4 7 4v8l-7 4-7-4V6z" />
                    </svg>
                  )}
                </div>
              </div>

              {/* Delete Button (shown on hover) */}
              <button
                onClick={() => handleDeleteTask(task.id)}
                className="flex-shrink-0 opacity-0 group-hover:opacity-100 text-gray-400 hover:text-red-500 transition-all"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Completed Tasks Section */}
      {completedTasks.length > 0 && (
        <div className="mt-8 pt-6 border-t border-gray-200">
          <details className="group">
            <summary className="cursor-pointer text-sm font-medium text-gray-600 hover:text-gray-900 flex items-center space-x-2">
              <svg className="w-4 h-4 transition-transform group-open:rotate-90" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
              </svg>
              <span>Completed ({completedTasks.length})</span>
            </summary>
            <div className="mt-3 space-y-1">
              {completedTasks.map((task) => (
                <div
                  key={task.id}
                  className="group bg-white hover:bg-gray-50 border border-gray-100 rounded-lg p-4 transition-all"
                >
                  <div className="flex items-start space-x-3">
                    <button
                      onClick={() => handleToggleComplete(task)}
                      className="mt-0.5 flex-shrink-0 w-5 h-5 rounded-full bg-gray-400 border-2 border-gray-400 transition-all flex items-center justify-center hover:bg-gray-500"
                    >
                      <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                      </svg>
                    </button>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm line-through text-gray-400">{task.title}</p>
                      {task.description && (
                        <p className="text-xs text-gray-400 mt-1">{task.description}</p>
                      )}
                    </div>
                    <button
                      onClick={() => handleDeleteTask(task.id)}
                      className="flex-shrink-0 opacity-0 group-hover:opacity-100 text-gray-400 hover:text-red-500 transition-all"
                    >
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </details>
        </div>
      )}
    </div>
  );
}
