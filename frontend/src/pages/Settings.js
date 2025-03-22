import React, { useState } from 'react';
import { User, Upload } from 'lucide-react';
import Sidebar from '../components/Sidebar';
import TopBar from '../components/TopBar';
import '../styles/Settings.css';

function Settings() {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeTab, setActiveTab] = useState('User Profile');
  const [userProfile, setUserProfile] = useState({
    fullName: 'John Doe',
    email: 'john.doe@example.com',
    bio: 'Product Manager with 5+ years of experience',
  });
  const [notifications, setNotifications] = useState({
    email: true,
    push: false,
    tasks: true,
    digest: false,
    mentions: true,
  });
  const [appearance, setAppearance] = useState({
    theme: 'System',
    compactMode: false,
    highContrast: false,
  });

  // Handle User Profile form changes
  const handleProfileChange = (e) => {
    const { name, value } = e.target;
    setUserProfile((prev) => ({ ...prev, [name]: value }));
  };

  // Handle User Profile form submission
  const handleProfileSubmit = (e) => {
    e.preventDefault();
    alert('Profile updated successfully!');
  };

  // Handle Notification toggles
  const handleNotificationToggle = (key) => {
    setNotifications((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  // Handle Appearance changes
  const handleThemeChange = (theme) => {
    setAppearance((prev) => ({ ...prev, theme }));
  };

  const handleAppearanceToggle = (key) => {
    setAppearance((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  return (
    <div className="dashboard">
      <Sidebar isSidebarCollapsed={isSidebarCollapsed} setIsSidebarCollapsed={setIsSidebarCollapsed} />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="settings-content">
          <div className="tabs">
            <button
              className={activeTab === 'User Profile' ? 'active' : ''}
              onClick={() => setActiveTab('User Profile')}
            >
              User Profile
            </button>
            <button
              className={activeTab === 'Notifications' ? 'active' : ''}
              onClick={() => setActiveTab('Notifications')}
            >
              Notifications
            </button>
            <button
              className={activeTab === 'Appearance' ? 'active' : ''}
              onClick={() => setActiveTab('Appearance')}
            >
              Appearance
            </button>
          </div>

          <div className="tab-content">
            {activeTab === 'User Profile' && (
              <div className="user-profile">
                <h2>User Profile</h2>
                <p>Manage your personal information and preferences.</p>
                <form onSubmit={handleProfileSubmit}>
                  <div className="avatar-section">
                    <div className="avatar">
                      <User size={48} />
                    </div>
                    <button type="button" className="change-avatar-btn">
                      <Upload size={16} /> Change Avatar
                    </button>
                  </div>
                  <div className="form-group">
                    <label>Full Name</label>
                    <input
                      type="text"
                      name="fullName"
                      value={userProfile.fullName}
                      onChange={handleProfileChange}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Email</label>
                    <input
                      type="email"
                      name="email"
                      value={userProfile.email}
                      onChange={handleProfileChange}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Bio</label>
                    <input
                      type="text"
                      name="bio"
                      value={userProfile.bio}
                      onChange={handleProfileChange}
                    />
                  </div>
                  <button type="submit" className="save-btn">
                    Save Changes
                  </button>
                </form>
              </div>
            )}

            {activeTab === 'Notifications' && (
              <div className="notifications">
                <h2>Notification Preferences</h2>
                <p>Configure how you receive notifications.</p>
                <div className="notification-option">
                  <label>Email Notifications</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={notifications.email}
                      onChange={() => handleNotificationToggle('email')}
                    />
                    <span className="slider"></span>
                  </label>
                </div>
                <div className="notification-option">
                  <label>Push Notifications</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={notifications.push}
                      onChange={() => handleNotificationToggle('push')}
                    />
                    <span className="slider"></span>
                  </label>
                </div>
                <div className="notification-option">
                  <label>Task Reminders</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={notifications.tasks}
                      onChange={() => handleNotificationToggle('tasks')}
                    />
                    <span className="slider"></span>
                  </label>
                </div>
                <div className="notification-option">
                  <label>Weekly Digest</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={notifications.digest}
                      onChange={() => handleNotificationToggle('digest')}
                    />
                    <span className="slider"></span>
                  </label>
                </div>
                <div className="notification-option">
                  <label>Mentions</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={notifications.mentions}
                      onChange={() => handleNotificationToggle('mentions')}
                    />
                    <span className="slider"></span>
                  </label>
                </div>
              </div>
            )}

            {activeTab === 'Appearance' && (
              <div className="appearance">
                <h2>Appearance</h2>
                <p>Customize the look and feel of your workspace.</p>
                <div className="appearance-option">
                  <label>Theme</label>
                  <div className="theme-buttons">
                    <button
                      className={appearance.theme === 'Light' ? 'active' : ''}
                      onClick={() => handleThemeChange('Light')}
                    >
                      Light
                    </button>
                    <button
                      className={appearance.theme === 'Dark' ? 'active' : ''}
                      onClick={() => handleThemeChange('Dark')}
                    >
                      Dark
                    </button>
                    <button
                      className={appearance.theme === 'System' ? 'active' : ''}
                      onClick={() => handleThemeChange('System')}
                    >
                      System
                    </button>
                  </div>
                </div>
                <div className="appearance-option">
                  <label>Compact Mode</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={appearance.compactMode}
                      onChange={() => handleAppearanceToggle('compactMode')}
                    />
                    <span className="slider"></span>
                  </label>
                  <p className="description">Use a more compact user interface</p>
                </div>
                <div className="appearance-option">
                  <label>High Contrast</label>
                  <label className="switch">
                    <input
                      type="checkbox"
                      checked={appearance.highContrast}
                      onChange={() => handleAppearanceToggle('highContrast')}
                    />
                    <span className="slider"></span>
                  </label>
                  <p className="description">Increase contrast for better readability</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Settings;