// src/Dashboard.js
import React, { useState } from 'react';
import { Clock, AlertCircle, CheckCircle } from 'lucide-react';
import '../styles/Dashboard.css';
import Sidebar from '../components/Sidebar';
import TopBar from '../components/TopBar';

function Dashboard() {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  const stats = {
    tasksDueToday: 3,
    inProgress: 5,
    completed: 12,
  };

  const recentTasks = [
    { title: 'Design System Update', due: 'Due in 2 days' },
    { title: 'Client Meeting Preparation', due: 'Due today' },
    { title: 'Weekly Report', due: 'Due in 4 days' },
  ];

  return (
    <div className="dashboard">
      <Sidebar
        className="sidebar"
        isSidebarCollapsed={isSidebarCollapsed}
        setIsSidebarCollapsed={setIsSidebarCollapsed}
      />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="top-bar-left">
          <h1>Dashboard</h1>
        </div>
        <div className="main-content">
          <div className="stats">
            <div className="stat-card">
              <div className="stat-icon">
                <Clock size={24} className="tasks-due-icon" />
              </div>
              <h4>Tasks Due Today</h4>
              <p>{stats.tasksDueToday}</p>
            </div>
            <div className="stat-card">
              <div className="stat-icon">
                <AlertCircle size={24} className="in-progress-icon" />
              </div>
              <h4>In Progress</h4>
              <p>{stats.inProgress}</p>
            </div>
            <div className="stat-card">
              <div className="stat-icon">
                <CheckCircle size={24} className="completed-icon" />
              </div>
              <h4>Completed</h4>
              <p>{stats.completed}</p>
            </div>
          </div>
          <div className="recent-tasks">
            <h3>Recent Tasks</h3>
            <ul>
              {recentTasks.map((task, index) => (
                <li key={index}>
                  <span className="task-title">{task.title}</span>
                  <span className="task-due">{task.due}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;