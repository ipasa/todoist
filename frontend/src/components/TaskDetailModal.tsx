import { useState, useEffect } from "react";
import type { Task } from "@/types/task.types";
import { taskApi } from "@/api/task.api";

interface TaskDetailModalProps {
  task: Task;
  isOpen: boolean;
  onClose: () => void;
  onUpdate: (updatedTask: Task) => void;
  onDelete: () => void;
}

export function TaskDetailModal({ task, isOpen, onClose, onUpdate, onDelete }: TaskDetailModalProps) {
  const [editedTask, setEditedTask] = useState<Task>(task);
  const [isEditing, setIsEditing] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    setEditedTask(task);
  }, [task]);

  if (!isOpen) return null;

  const handleSave = async () => {
    setIsSaving(true);
    try {
      let dueDate = editedTask.due_date;
      if (dueDate && !dueDate.includes('T')) {
        const date = new Date(dueDate + 'T00:00:00');
        dueDate = date.toISOString();
      }

      const updated = await taskApi.updateTask(task.id, {
        title: editedTask.title,
        description: editedTask.description,
        priority: editedTask.priority,
        status: editedTask.status,
        due_date: dueDate,
      });
      onUpdate(updated);
      setIsEditing(false);
    } catch (err) {
      console.error("Failed to update task", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleComplete = async () => {
    const newStatus = editedTask.status === "completed" ? "pending" : "completed";
    setIsSaving(true);
    try {
      const updated = await taskApi.updateTask(task.id, { status: newStatus });
      onUpdate(updated);
      setEditedTask({ ...editedTask, status: newStatus });
    } catch (err) {
      console.error("Failed to toggle completion", err);
    } finally {
      setIsSaving(false);
    }
  };

  const getPriorityLabel = (priority: number) => {
    return `P${5 - priority}`;
  };

  const getPriorityColor = (priority: number) => {
    switch (priority) {
      case 4: return "text-red-600 bg-red-50";
      case 3: return "text-orange-600 bg-orange-50";
      case 2: return "text-blue-600 bg-blue-50";
      default: return "text-gray-600 bg-gray-50";
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-start justify-center z-50 pt-16" onClick={onClose}>
      <div
        className="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[80vh] overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200">
          <div className="flex items-center space-x-2 text-sm text-gray-500">
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
            <span>Inbox</span>
          </div>
          <div className="flex items-center space-x-2">
            <button className="p-2 hover:bg-gray-100 rounded">
              <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
              </svg>
            </button>
            <button className="p-2 hover:bg-gray-100 rounded">
              <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <button className="p-2 hover:bg-gray-100 rounded">
              <svg className="w-5 h-5 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
              </svg>
            </button>
            <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded">
              <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="overflow-y-auto max-h-[calc(80vh-120px)]">
          <div className="px-6 py-4">
            {/* Task Title */}
            <div className="flex items-start space-x-3 mb-6">
              <button
                onClick={handleComplete}
                disabled={isSaving}
                className={`mt-1 flex-shrink-0 w-6 h-6 rounded-full border-2 transition-all flex items-center justify-center ${
                  editedTask.status === "completed"
                    ? "bg-gray-400 border-gray-400"
                    : "border-gray-300 hover:border-gray-400"
                }`}
              >
                {editedTask.status === "completed" && (
                  <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                  </svg>
                )}
              </button>
              <div className="flex-1">
                {isEditing ? (
                  <input
                    type="text"
                    value={editedTask.title}
                    onChange={(e) => setEditedTask({ ...editedTask, title: e.target.value })}
                    className="w-full text-xl font-semibold border-none focus:outline-none focus:ring-0 p-0"
                  />
                ) : (
                  <h2
                    className={`text-xl font-semibold cursor-pointer hover:text-gray-600 ${
                      editedTask.status === "completed" ? "line-through text-gray-400" : ""
                    }`}
                    onClick={() => setIsEditing(true)}
                  >
                    {editedTask.title}
                  </h2>
                )}
              </div>
            </div>

            {/* Description */}
            <div className="mb-6">
              {isEditing || editedTask.description ? (
                <div>
                  <div className="flex items-center space-x-2 text-sm text-gray-500 mb-2">
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16m-7 6h7" />
                    </svg>
                    <span>Description</span>
                  </div>
                  {isEditing ? (
                    <textarea
                      value={editedTask.description || ""}
                      onChange={(e) => setEditedTask({ ...editedTask, description: e.target.value })}
                      rows={3}
                      className="w-full text-sm text-gray-700 border border-gray-300 rounded-lg p-2 focus:outline-none focus:ring-2 focus:ring-red-500"
                      placeholder="Add description..."
                    />
                  ) : (
                    <p
                      className="text-sm text-gray-700 cursor-pointer hover:bg-gray-50 p-2 rounded"
                      onClick={() => setIsEditing(true)}
                    >
                      {editedTask.description}
                    </p>
                  )}
                </div>
              ) : (
                <button
                  onClick={() => setIsEditing(true)}
                  className="flex items-center space-x-2 text-sm text-gray-500 hover:text-gray-700"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                  <span>Add description</span>
                </button>
              )}
            </div>

            {/* Sub-tasks placeholder */}
            <div className="mb-6">
              <button className="flex items-center space-x-2 text-sm text-gray-500 hover:text-gray-700">
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                <span>Add sub-task</span>
              </button>
            </div>

            {/* Comments placeholder */}
            <div className="border-t border-gray-200 pt-4">
              <div className="flex items-center space-x-3">
                <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-xs text-white font-semibold">
                  H
                </div>
                <input
                  type="text"
                  placeholder="Comment"
                  className="flex-1 text-sm border-none focus:outline-none focus:ring-0 bg-transparent"
                />
                <button className="p-2 hover:bg-gray-100 rounded">
                  <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="absolute right-0 top-0 w-64 h-full border-l border-gray-200 bg-white px-4 py-6 space-y-4">
          {/* Project */}
          <div>
            <div className="text-xs font-medium text-gray-500 mb-2">Project</div>
            <button className="flex items-center space-x-2 text-sm text-gray-700 hover:bg-gray-50 px-2 py-1 rounded w-full">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
              </svg>
              <span>Inbox</span>
            </button>
          </div>

          {/* Due Date */}
          <div>
            <div className="text-xs font-medium text-gray-500 mb-2">Date</div>
            <label className="flex items-center space-x-2 text-sm text-gray-700 hover:bg-gray-50 px-2 py-1 rounded w-full cursor-pointer">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span>
                {editedTask.due_date
                  ? new Date(editedTask.due_date).toLocaleDateString("en-US", { month: "short", day: "numeric" })
                  : "No date"}
              </span>
              <input
                type="date"
                value={editedTask.due_date ? editedTask.due_date.split('T')[0] : ""}
                onChange={(e) => setEditedTask({ ...editedTask, due_date: e.target.value })}
                className="absolute opacity-0 w-0 h-0"
              />
            </label>
          </div>

          {/* Priority */}
          <div>
            <div className="text-xs font-medium text-gray-500 mb-2">Priority</div>
            <label className={`flex items-center space-x-2 text-sm px-2 py-1 rounded w-full cursor-pointer ${getPriorityColor(editedTask.priority)}`}>
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M3 6l7-4 7 4v8l-7 4-7-4V6z" />
              </svg>
              <span>{getPriorityLabel(editedTask.priority)}</span>
              <select
                value={editedTask.priority}
                onChange={(e) => setEditedTask({ ...editedTask, priority: Number(e.target.value) })}
                className="absolute opacity-0 w-0 h-0"
              >
                <option value="1">P4</option>
                <option value="2">P3</option>
                <option value="3">P2</option>
                <option value="4">P1</option>
              </select>
            </label>
          </div>

          {/* Labels */}
          <div>
            <div className="text-xs font-medium text-gray-500 mb-2">Labels</div>
            <button className="flex items-center space-x-2 text-sm text-gray-500 hover:bg-gray-50 px-2 py-1 rounded w-full">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
            </button>
          </div>

          {/* Actions */}
          <div className="pt-4 border-t border-gray-200 space-y-2">
            {isEditing && (
              <button
                onClick={handleSave}
                disabled={isSaving}
                className="w-full px-3 py-2 bg-red-500 text-white text-sm rounded hover:bg-red-600 disabled:opacity-50"
              >
                {isSaving ? "Saving..." : "Save"}
              </button>
            )}
            <button
              onClick={() => {
                if (confirm("Delete this task?")) {
                  onDelete();
                  onClose();
                }
              }}
              className="w-full px-3 py-2 text-sm text-red-600 hover:bg-red-50 rounded"
            >
              Delete task
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
