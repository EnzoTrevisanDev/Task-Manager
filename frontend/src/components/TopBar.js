import React, { useState } from 'react';
import { Bell, LogOut } from 'lucide-react';
import '../styles/TopBar.css';

function TopBar({ isSidebarCollapsed }) {
  const [isNotificationOpen, setIsNotificationOpen] = useState(false);

  // Mock notification data (replace with real data from backend later)
  const notifications = [
    { id: 1, message: 'Project "Website Redesign" updated', time: '2 hours ago' },
    { id: 2, message: 'New task assigned to you', time: '1 day ago' },
    { id: 3, message: 'New task assigned to you', time: '2 days ago' },
  ];

  const toggleNotifications = () => {
    setIsNotificationOpen(!isNotificationOpen);
  };

  return (
    <div className={`top-bar ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
      <div className="top-bar-left">
        <h1>Dashboard</h1>
      </div>
      <div className="top-bar-right">
        <div className="user-info">
          <div className="notification-wrapper">
            <Bell size={20} onClick={toggleNotifications} style={{ cursor: 'pointer' }} />
            <span className="notification-count">{notifications.length}</span>
            {isNotificationOpen && (
              <div className="notification-dropdown">
                <h4>Notifications</h4>
                <ul>
                  {notifications.length > 0 ? (
                    notifications.map((notification) => (
                      <li key={notification.id}>
                        <span className="notification-message">{notification.message}</span>
                        <span className="notification-time">{notification.time}</span>
                      </li>
                    ))
                  ) : (
                    <li>No new updates</li>
                  )}
                </ul>
              </div>
            )}
          </div>
          <span className="username">enzo</span>
          <button className="logout-btn">
            <LogOut size={20} />
            <span>Logout</span>
          </button>
        </div>
      </div>
    </div>
  );
}

export default TopBar;