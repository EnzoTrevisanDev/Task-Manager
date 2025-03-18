import React from 'react';
import Sidebar from '../components/Sidebar';

function Projects() {
  return (
    <div className="dashboard">
      <Sidebar />
      <div className="dashboard-content">
        <h1>Projects Page</h1>
        <p>This is where you’ll manage your projects.</p>
      </div>
    </div>
  );
}

export default Projects;