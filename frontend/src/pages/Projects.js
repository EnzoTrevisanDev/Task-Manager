import React, { useState, useEffect } from 'react';
import { Star, Clock, AlertCircle, CheckCircle } from 'lucide-react';
import { Link } from 'react-router-dom';
import '../styles/Projects.css';
import Sidebar from '../components/Sidebar';
import TopBar from '../components/TopBar';
import ProjectCreation from '../components/ProjectCreation';
import api from '../api/api';

function Projects() {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeTab, setActiveTab] = useState('All Projects');
  const [activeCategory, setActiveCategory] = useState('All');
  const [activeStatus, setActiveStatus] = useState('All');
  const [projectsState, setProjectsState] = useState([]);
  const [showProjectCreation, setShowProjectCreation] = useState(false);
  const [error, setError] = useState('');
  const [refreshProjects, setRefreshProjects] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchProjects = async () => {
    setIsLoading(true);
    try {
      const response = await api.get('/projects', {
        params: { _t: new Date().getTime() }, // Cache-busting parameter
      });
      console.log('Fetched projects:', response.data);
      const mappedProjects = response.data.map((project) => ({
        id: project.id,
        title: project.name,
        description: project.description || 'No description provided',
        tasks: project.tasks ? project.tasks.length : 0,
        category: project.category ? project.category.charAt(0).toUpperCase() + project.category.slice(1).toLowerCase() : 'Uncategorized',
        status: project.status ? project.status.charAt(0).toUpperCase() + project.status.slice(1).toLowerCase() : 'Planned',
        isFavorite: project.is_favorite || false,
        creator: project.creator ? project.creator.name : 'Unknown',
      }));
      console.log('Mapped projects:', mappedProjects);
      setProjectsState(mappedProjects);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch projects. Please try again.');
      console.error('Error fetching projects:', err);
      // If the error is a 401, the interceptor should have redirected to login
      if (err.response?.status === 401) {
        console.log('401 error after retry - redirecting to login');
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, [refreshProjects]);

  const filteredProjects = projectsState.filter((project) => {
    if (activeTab === 'Favorites') return project.isFavorite;
    const categoryMatch = activeCategory === 'All' || project.category === activeCategory;
    const statusMatch = activeStatus === 'All' || project.status === activeStatus;
    return categoryMatch && statusMatch;
  });

  console.log('Active filters:', { activeTab, activeCategory, activeStatus });
  console.log('Filtered projects:', filteredProjects);

  const toggleFavorite = async (projectId) => {
    try {
      const project = projectsState.find((p) => p.id === projectId);
      const newFavoriteStatus = !project.isFavorite;
      await api.put(`/projects/${projectId}/favorite`, { is_favorite: newFavoriteStatus });
      setProjectsState((prev) =>
        prev.map((project) =>
          project.id === projectId ? { ...project, isFavorite: newFavoriteStatus } : project
        )
      );
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update favorite status. Please try again.');
      console.error('Error toggling favorite:', err);
    }
  };

  const handleNewProject = () => {
    setShowProjectCreation(true);
  };

  const handleAddProject = async (newProject) => {
    try {
      console.log('Creating new project:', newProject);
      const response = await api.post('/projects', {
        name: newProject.title,
        description: newProject.description,
        category: newProject.category,
        status: newProject.status,
      });
      console.log('Project creation response:', response.data);
      setShowProjectCreation(false);
      setActiveTab('All Projects');
      setActiveCategory('All');
      setActiveStatus('All');
      await new Promise((resolve) => setTimeout(resolve, 500));
      setRefreshProjects((prev) => !prev);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create project. Please try again.');
      console.error('Error creating project:', err);
    }
  };

  return (
    <div className="dashboard">
      <Sidebar
        isSidebarCollapsed={isSidebarCollapsed}
        setIsSidebarCollapsed={setIsSidebarCollapsed}
      />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="projects-content">
          <button className="new-project-btn" onClick={handleNewProject}>
            + New Project
          </button>
          {error && <p className="error-message">{error}</p>}
          {isLoading ? (
            <p className="loading-message">Loading projects...</p>
          ) : (
            <>
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
                {filteredProjects.length === 0 ? (
                  <p className="no-projects-message">No projects found.</p>
                ) : (
                  <>
                    {console.log('Rendering projects:', filteredProjects)}
                    {filteredProjects.map((project) => (
                      <div key={project.id} className="project-card">
                        <div className="card-header">
                          <h3>{project.title}</h3>
                          <span className="project-tag">{project.category}</span>
                          <Star
                            size={16}
                            className={project.isFavorite ? 'favorite-star' : 'unfavorite-star'}
                            onClick={() => toggleFavorite(project.id)}
                          />
                        </div>
                        <p className="project-description">{project.description}</p>
                        <div className="card-stats">
                          <span className="task-count">∥ {project.tasks} tasks</span>
                          <span className="project-status">
                            {project.status === 'In Progress' && (
                              <AlertCircle size={16} className="status-icon in-progress" />
                            )}
                            {project.status === 'Planned' && (
                              <Clock size={16} className="status-icon planned" />
                            )}
                            {project.status === 'Completed' && (
                              <CheckCircle size={16} className="status-icon completed" />
                            )}
                            {project.status}
                          </span>
                        </div>
                        <Link to={`/project/${project.id}`} className="view-details-btn">
                          View Details →
                        </Link>
                      </div>
                    ))}
                  </>
                )}
              </div>
            </>
          )}
        </div>
      </div>
      {showProjectCreation && (
        <ProjectCreation
          onAddProject={handleAddProject}
          onClose={() => setShowProjectCreation(false)}
        />
      )}
    </div>
  );
}

export default Projects;