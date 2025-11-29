import { useState, useEffect } from "react";
import type { Task } from "@/types/task.types";
import { taskApi } from "@/api/task.api";

interface TaskListRefinedProps {
  refreshTrigger?: number;
  onTaskCreated?: () => void;
  onTaskClick?: (task: Task) => void;
  onTaskCountChange?: (count: number) => void;
}

export function TaskListRefined({ refreshTrigger, onTaskCreated, onTaskClick, onTaskCountChange }: TaskListRefinedProps) {
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
    setCompletingTasks(prev => new Set(prev).add(task.id));

    try {
      await taskApi.updateTask(task.id, { status: newStatus });
      setTasks(tasks.map(t =>
        t.id === task.id ? { ...t, status: newStatus } : t
      ));

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
    try {
      await taskApi.deleteTask(id);
      setTasks(tasks.filter((task) => task.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to delete task");
    }
  };

  const getPriorityColor = (priority: number) => {
    switch (priority) {
      case 4: return "border-red-500";
      case 3: return "border-orange-400";
      case 2: return "border-blue-500";
      default: return "border-transparent";
    }
  };

  const getPriorityColorBg = (priority: number) => {
    switch (priority) {
      case 4: return "text-red-600 bg-red-50";
      case 3: return "text-orange-600 bg-orange-50";
      case 2: return "text-blue-600 bg-blue-50";
      default: return "text-gray-500";
    }
  };

  const groupTasksByDate = (taskList: Task[]) => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const groups: { [key: string]: Task[] } = {
      overdue: [],
      today: [],
      upcoming: [],
      nodate: [],
    };

    taskList.forEach(task => {
      if (task.status === "completed") return;

      if (!task.due_date) {
        groups.nodate.push(task);
        return;
      }

      const dueDate = new Date(task.due_date);
      dueDate.setHours(0, 0, 0, 0);

      if (dueDate < today) {
        groups.overdue.push(task);
      } else if (dueDate.getTime() === today.getTime()) {
        groups.today.push(task);
      } else {
        groups.upcoming.push(task);
      }
    });

    return groups;
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
  const groupedTasks = groupTasksByDate(activeTasks);

  const renderTask = (task: Task) => (
    <div
      key={task.id}
      className={`group bg-white border-l-2 ${getPriorityColor(task.priority)} px-3 py-2.5 rounded-lg hover:shadow-sm transition-all duration-150 ${
        completingTasks.has(task.id) ? "opacity-50" : ""
      }`}
    >
      <div className="flex items-start space-x-3">
        <button
          onClick={() => handleToggleComplete(task)}
          className="mt-0.5 flex-shrink-0 w-[18px] h-[18px] rounded-full border-2 border-gray-400 hover:border-gray-500 transition-colors"
        />

        <div
          className="flex-1 min-w-0 cursor-pointer"
          onClick={() => onTaskClick?.(task)}
        >
          <p className="text-[13px] text-gray-800 leading-relaxed">
            {task.title}
          </p>
          {task.description && (
            <p className="text-xs text-gray-500 mt-1 leading-relaxed">{task.description}</p>
          )}
          {task.due_date && (
            <div className="flex items-center space-x-1 mt-1.5">
              <span className="text-[11px] text-green-700 font-medium">
                {new Date(task.due_date).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
              </span>
            </div>
          )}
        </div>

        <button
          onClick={(e) => {
            e.stopPropagation();
            handleDeleteTask(task.id);
          }}
          className="flex-shrink-0 opacity-0 group-hover:opacity-100 text-gray-400 hover:text-gray-600 transition-opacity p-1"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      {!showAddForm ? (
        <button
          onClick={() => setShowAddForm(true)}
          className="group flex items-center space-x-2 text-gray-500 hover:text-red-600 transition-colors"
        >
          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
          </svg>
          <span className="text-[13px] font-medium">Add task</span>
        </button>
      ) : (
        <form onSubmit={handleQuickAdd} className="bg-white border border-gray-200 rounded-lg p-3 shadow-sm">
          <input
            type="text"
            value={newTaskTitle}
            onChange={(e) => setNewTaskTitle(e.target.value)}
            placeholder="Task name"
            autoFocus
            className="w-full text-[13px] border-none focus:outline-none focus:ring-0 p-0 mb-2 placeholder-gray-400"
          />

          {showDescription && (
            <textarea
              value={newTaskDescription}
              onChange={(e) => setNewTaskDescription(e.target.value)}
              placeholder="Description"
              rows={2}
              className="w-full text-xs text-gray-600 border-none focus:outline-none focus:ring-0 p-0 mb-2 resize-none placeholder-gray-400"
            />
          )}

          <div className="flex items-center justify-between mt-2 pt-2 border-t border-gray-100">
            <div className="flex items-center space-x-1">
              <label className={`text-xs px-2 py-1 hover:bg-gray-100 rounded flex items-center space-x-1 cursor-pointer transition-colors ${
                newTaskDueDate ? "text-green-700 bg-green-50" : "text-gray-600"
              }`}>
                <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <span className="text-[11px] font-medium">
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

              <label className={`text-xs px-2 py-1 hover:bg-gray-100 rounded flex items-center space-x-1 cursor-pointer transition-colors ${
                newTaskPriority > 1 ? getPriorityColorBg(newTaskPriority) : "text-gray-600"
              }`}>
                <svg className="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M3 6l7-4 7 4v8l-7 4-7-4V6z" />
                </svg>
                <span className="text-[11px] font-medium">P{5 - newTaskPriority}</span>
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

              {!showDescription && (
                <button
                  type="button"
                  onClick={() => setShowDescription(true)}
                  className="text-xs text-gray-600 hover:bg-gray-100 px-2 py-1 rounded transition-colors"
                >
                  <span className="text-[11px] font-medium">Description</span>
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
                className="text-xs px-3 py-1.5 text-gray-600 hover:bg-gray-100 rounded transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isAdding || !newTaskTitle.trim()}
                className="text-xs px-3 py-1.5 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isAdding ? "Adding..." : "Add task"}
              </button>
            </div>
          </div>
        </form>
      )}

      {/* Date-grouped Tasks */}
      {activeTasks.length === 0 && !showAddForm ? (
        <div className="text-center py-16">
          <div className="mb-4">
            <svg className="w-16 h-16 mx-auto text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-sm text-gray-400">All clear</p>
          <p className="text-xs text-gray-400 mt-1">Looks like everything's organized in the right place.</p>
        </div>
      ) : (
        <>
          {groupedTasks.overdue.length > 0 && (
            <div>
              <h3 className="text-xs font-bold text-red-600 mb-2 uppercase tracking-wide">Overdue</h3>
              <div className="space-y-1">
                {groupedTasks.overdue.map(renderTask)}
              </div>
            </div>
          )}

          {groupedTasks.today.length > 0 && (
            <div>
              <h3 className="text-xs font-bold text-gray-700 mb-2 uppercase tracking-wide">
                {new Date().toLocaleDateString("en-US", { month: "short", day: "numeric" })} · Today
              </h3>
              <div className="space-y-1">
                {groupedTasks.today.map(renderTask)}
              </div>
            </div>
          )}

          {groupedTasks.upcoming.length > 0 && (
            <div>
              <h3 className="text-xs font-bold text-gray-700 mb-2 uppercase tracking-wide">Upcoming</h3>
              <div className="space-y-1">
                {groupedTasks.upcoming.map(renderTask)}
              </div>
            </div>
          )}

          {groupedTasks.nodate.length > 0 && (
            <div>
              <div className="space-y-1">
                {groupedTasks.nodate.map(renderTask)}
              </div>
            </div>
          )}
        </>
      )}

      {completedTasks.length > 0 && (
        <div className="mt-8 pt-6 border-t border-gray-200">
          <details className="group">
            <summary className="cursor-pointer text-xs font-bold text-gray-500 hover:text-gray-700 flex items-center space-x-2 uppercase tracking-wide">
              <svg className="w-3 h-3 transition-transform group-open:rotate-90" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
              </svg>
              <span>{completedTasks.length} Completed</span>
            </summary>
            <div className="mt-3 space-y-1">
              {completedTasks.map((task) => (
                <div
                  key={task.id}
                  className="group bg-white px-3 py-2.5 rounded-lg hover:shadow-sm transition-all duration-150"
                >
                  <div className="flex items-start space-x-3">
                    <button
                      onClick={() => handleToggleComplete(task)}
                      className="mt-0.5 flex-shrink-0 w-[18px] h-[18px] rounded-full bg-gray-400 border-2 border-gray-400 transition-all flex items-center justify-center hover:bg-gray-500"
                    >
                      <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                      </svg>
                    </button>
                    <div className="flex-1 min-w-0">
                      <p className="text-[13px] line-through text-gray-400">{task.title}</p>
                      {task.description && (
                        <p className="text-xs text-gray-400 mt-1">{task.description}</p>
                      )}
                    </div>
                    <button
                      onClick={() => handleDeleteTask(task.id)}
                      className="flex-shrink-0 opacity-0 group-hover:opacity-100 text-gray-400 hover:text-gray-600 transition-opacity p-1"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
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
