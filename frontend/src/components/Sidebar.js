import React from 'react';
import { LayoutDashboard, Folder, BarChart, Settings, Menu, X, UserCircle } from 'lucide-react';
import { Link } from 'react-router-dom';
import '../styles/Sidebar.css';

function Sidebar({ isSidebarCollapsed, setIsSidebarCollapsed }) {
  const toggleSidebar = () => {
    setIsSidebarCollapsed(!isSidebarCollapsed);
  };

  return (
    <div className={`sidebar ${isSidebarCollapsed ? 'collapsed' : ''}`}>
      <div className="sidebar-content">
        <div className="sidebar-header">
          <div className="sidebar-logo">
            <h2>Optima</h2>
          </div>
          <button className="sidebar-toggle" onClick={toggleSidebar}>
            {isSidebarCollapsed ? <Menu size={20} /> : <X size={20} />}
          </button>
        </div>
        <div className="sidebar-menu">
          <h4 className="menu-title">MAIN</h4>
          <ul>
            <li className="active">
              <LayoutDashboard size={20} />
              <span className="menu-text">
                <Link to="/dashboard">Dashboard</Link>
              </span>
            </li>
            <li>
              <Folder size={20} />
              <span className="menu-text">
                <Link to="/projects">Projects</Link>
              </span>
            </li>
            <li>
              <BarChart size={20} />
              <span className="menu-text">
                <Link to="/analytics">Analytics</Link>
              </span>
            </li>
            <li>
              <Settings size={20} />
              <span className="menu-text">
                <Link to="/settings">Settings</Link>
              </span>
            </li>
          </ul>
        </div>
        <div className="sidebar-projects">
          <h4>PROJECTS</h4>
          <ul>
            <li className="projects-link">
              <Folder size={20} />
              <span className="menu-text">
                <Link to="/projects">Active Projects</Link>
              </span>
            </li>
            <li className="projects-link">
              <Folder size={20} />
              <span className="menu-text">
                <Link to="/archived">Archived</Link>
              </span>
            </li>
          </ul>
        </div>
      </div>
      <div className="sidebar-user">
        <div className="user-info">
          <UserCircle size={24} className="user-avatar" />
          <span className="menu-text user-email">john.doe@example.com</span>
        </div>
      </div>
    </div>
  );
}

export default Sidebar;