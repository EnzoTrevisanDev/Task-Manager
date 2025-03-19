import React, { useState } from 'react';
import { Star, Clock, AlertCircle, CheckCircle } from 'lucide-react';
import { Link } from 'react-router-dom'; // For navigation
import '../styles/Projects.css';
import Sidebar from '../components/Sidebar';
import TopBar from '../components/TopBar';

function Projects() {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeTab, setActiveTab] = useState('All Projects');
  const [activeCategory, setActiveCategory] = useState('All');
  const [activeStatus, setActiveStatus] = useState('All');

  // Mock project data
  const projects = [
    {
      id: 1,
      title: 'Marketing Campaign',
      description: 'Q1 2025 Digital Marketing Campaign',
      tasks: 12,
      category: 'Marketing',
      status: 'In Progress',
      isFavorite: true,
    },
    {
      id: 2,
      title: 'Website Redesign',
      description: 'Company Website overhaul project',
      tasks: 8,
      category: 'Development',
      status: 'Planned',
      isFavorite: true,
    },
    {
      id: 3,
      title: 'Product Launch',
      description: 'New product launch preparation',
      tasks: 15,
      category: 'Product',
      status: 'In Progress',
      isFavorite: true,
    },
    {
      id: 4,
      title: 'Customer Research',
      description: 'User feedback analysis and report',
      tasks: 12,
      category: 'Research',
      status: 'Completed',
      isFavorite: true,
    },
  ];

  // Filter projects based on category, status, and favorites
  const filteredProjects = projects.filter((project) => {
    if (activeTab === 'Favorites') return project.isFavorite;
    const categoryMatch = activeCategory === 'All' || project.category === activeCategory;
    const statusMatch = activeStatus === 'All' || project.status === activeStatus;
    return categoryMatch && statusMatch;
  });

  return (
    <div className="dashboard">
      <Sidebar isSidebarCollapsed={isSidebarCollapsed} setIsSidebarCollapsed={setIsSidebarCollapsed} />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="projects-content">
          <div className="projects-filters">
            <div className="filter-tabs">
              <button
                className={activeTab === 'All Projects' ? 'active' : ''}
                onClick={() => setActiveTab('All Projects')}
              >
                All Projects
              </button>
              <button
                className={activeTab === 'Favorites' ? 'active' : ''}
                onClick={() => setActiveTab('Favorites')}
              >
                Favorites
              </button>
            </div>
            {activeTab === 'All Projects' && (
              <div className="sub-filters">
                <div className="category-filters">
                  <button
                    className={activeCategory === 'All' ? 'active' : ''}
                    onClick={() => setActiveCategory('All')}
                  >
                    All
                  </button>
                  <button
                    className={activeCategory === 'Marketing' ? 'active' : ''}
                    onClick={() => setActiveCategory('Marketing')}
                  >
                    Marketing
                  </button>
                  <button
                    className={activeCategory === 'Development' ? 'active' : ''}
                    onClick={() => setActiveCategory('Development')}
                  >
                    Development
                  </button>
                  <button
                    className={activeCategory === 'Product' ? 'active' : ''}
                    onClick={() => setActiveCategory('Product')}
                  >
                    Product
                  </button>
                  <button
                    className={activeCategory === 'Research' ? 'active' : ''}
                    onClick={() => setActiveCategory('Research')}
                  >
                    Research
                  </button>
                </div>
                <div className="status-filters">
                  <button
                    className={activeStatus === 'All' ? 'active' : ''}
                    onClick={() => setActiveStatus('All')}
                  >
                    All
                  </button>
                  <button
                    className={activeStatus === 'Planned' ? 'active' : ''}
                    onClick={() => setActiveStatus('Planned')}
                  >
                    Planned
                  </button>
                  <button
                    className={activeStatus === 'In Progress' ? 'active' : ''}
                    onClick={() => setActiveStatus('In Progress')}
                  >
                    In Progress
                  </button>
                  <button
                    className={activeStatus === 'Completed' ? 'active' : ''}
                    onClick={() => setActiveStatus('Completed')}
                  >
                    Completed
                  </button>
                  <button
                    className={activeStatus === 'On Hold' ? 'active' : ''}
                    onClick={() => setActiveStatus('On Hold')}
                  >
                    On Hold
                  </button>
                </div>
              </div>
            )}
          </div>
          <div className="projects-grid">
            {filteredProjects.map((project) => (
              <div key={project.id} className="project-card">
                <div className="card-header">
                  <h3>{project.title}</h3>
                  {project.isFavorite && <Star size={16} className="favorite-star" />}
                </div>
                <p className="project-description">{project.description}</p>
                <div className="card-stats">
                  <span className="task-count">∥ {project.tasks} tasks</span>
                  <span className="project-status">
                    {project.status === 'In Progress' && <AlertCircle size={16} className="status-icon in-progress" />}
                    {project.status === 'Planned' && <Clock size={16} className="status-icon planned" />}
                    {project.status === 'Completed' && <CheckCircle size={16} className="status-icon completed" />}
                    {project.status}
                  </span>
                </div>
                <Link to={`/project/${project.id}`} className="view-details-btn">
                  View Details →
                </Link>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Projects;