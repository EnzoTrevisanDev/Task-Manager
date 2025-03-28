import React from 'react';
import { LayoutDashboard, Folder, BarChart, Settings, UserCircle } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import '../styles/Sidebar.css';

function Sidebar({ isSidebarCollapsed, setIsSidebarCollapsed }) {
  const location = useLocation();
  const currentPath = location.pathname;

  return (
    <div className={`sidebar ${isSidebarCollapsed ? 'collapsed' : ''}`}>
      {/* Sidebar header with toggle button */}
      <div className="sidebar-header">
        <h4 className={`menu-title ${isSidebarCollapsed ? 'hidden' : ''}`}>Optima</h4>
        <button
          onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
          className="sidebar-toggle"
        >
          {isSidebarCollapsed ? '→' : '←'}
        </button>
      </div>
      <div className="sidebar-menu">
        <ul>
          <li className={currentPath === '/dashboard' ? 'active' : ''}>
            <LayoutDashboard size={20} />
            <span className={`menu-text ${isSidebarCollapsed ? 'hidden' : ''}`}>
              <Link to="/dashboard">Dashboard</Link>
            </span>
          </li>
          <li className={currentPath === '/projects' ? 'active' : ''}>
            <Folder size={20} />
            <span className={`menu-text ${isSidebarCollapsed ? 'hidden' : ''}`}>
              <Link to="/projects">Projects</Link>
            </span>
          </li>
          <li className={currentPath === '/analytics' ? 'active' : ''}>
            <BarChart size={20} />
            <span className={`menu-text ${isSidebarCollapsed ? 'hidden' : ''}`}>
              <Link to="/analytics">Analytics</Link>
            </span>
          </li>
          <li className={currentPath === '/settings' ? 'active' : ''}>
            <Settings size={20} />
            <span className={`menu-text ${isSidebarCollapsed ? 'hidden' : ''}`}>
              <Link to="/settings">Settings</Link>
            </span>
          </li>
        </ul>
      </div>
      <div className="sidebar-user">
        <div className="user-info">
          <UserCircle size={24} className="user-avatar" />
          <span className={`menu-text user-email ${isSidebarCollapsed ? 'hidden' : ''}`}>
            john.doe@example.com
          </span>
        </div>
      </div>
    </div>
  );
}

export default Sidebar;