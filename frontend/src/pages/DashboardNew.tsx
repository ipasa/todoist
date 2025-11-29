import { useNavigate } from "react-router-dom";
import { useAuthStore } from "@/store/authStore";
import { TaskListRefined } from "@/components/TaskListRefined";
import { Sidebar } from "@/components/Sidebar";
import { TaskDetailModal } from "@/components/TaskDetailModal";
import { useState } from "react";
import type { Task } from "@/types/task.types";

export function DashboardNew() {
  const navigate = useNavigate();
  const { logout } = useAuthStore();
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [activeView, setActiveView] = useState("inbox");
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [taskCount, setTaskCount] = useState(0);

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  const handleTaskCreated = () => {
    setRefreshTrigger((prev) => prev + 1);
  };

  const handleTaskClick = (task: Task) => {
    setSelectedTask(task);
  };

  const handleTaskUpdate = (updatedTask: Task) => {
    setRefreshTrigger((prev) => prev + 1);
    setSelectedTask(updatedTask);
  };

  const handleTaskDelete = () => {
    setRefreshTrigger((prev) => prev + 1);
  };

  const getViewTitle = () => {
    switch (activeView) {
      case "inbox": return "Inbox";
      case "today": return "Today";
      case "upcoming": return "Upcoming";
      case "filters": return "Filters & Labels";
      default: return "Inbox";
    }
  };

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <Sidebar
        activeView={activeView}
        onViewChange={setActiveView}
        taskCount={taskCount}
      />

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Navigation */}
        <nav className="bg-white border-b border-gray-200">
          <div className="px-6 py-3">
            <div className="flex justify-between items-center">
              <div className="flex items-center space-x-3">
                <h1 className="text-2xl font-bold text-gray-900">{getViewTitle()}</h1>
              </div>
              <div className="flex items-center space-x-4">
                <button className="p-2 hover:bg-gray-100 rounded">
                  <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                  </svg>
                </button>
                <button className="p-2 hover:bg-gray-100 rounded">
                  <svg className="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                  </svg>
                </button>
                <button
                  onClick={handleLogout}
                  className="text-sm text-gray-600 hover:text-gray-900"
                >
                  Logout
                </button>
              </div>
            </div>
          </div>
        </nav>

        {/* Date Header */}
        {activeView === "today" && (
          <div className="px-6 py-3 bg-white border-b border-gray-200">
            <p className="text-sm text-gray-500">
              {new Date().toLocaleDateString("en-US", {
                weekday: "long",
                month: "long",
                day: "numeric"
              })}
            </p>
          </div>
        )}

        {/* Task List */}
        <main className="flex-1 overflow-y-auto px-6 py-6">
          <TaskListRefined
            refreshTrigger={refreshTrigger}
            onTaskCreated={handleTaskCreated}
            onTaskClick={handleTaskClick}
            onTaskCountChange={setTaskCount}
          />
        </main>
      </div>

      {/* Task Detail Modal */}
      {selectedTask && (
        <TaskDetailModal
          task={selectedTask}
          isOpen={!!selectedTask}
          onClose={() => setSelectedTask(null)}
          onUpdate={handleTaskUpdate}
          onDelete={handleTaskDelete}
        />
      )}
    </div>
  );
}
