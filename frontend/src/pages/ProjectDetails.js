import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { AlertCircle, Clock, CheckCircle, User, X, GripVertical, Trash2 } from 'lucide-react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar, Line } from 'react-chartjs-2';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import Sidebar from '../components/Sidebar';
import TopBar from '../components/TopBar';
import api from '../api/api';
import '../styles/ProjectDetails.css';

// Register ChartJS components
ChartJS.register(CategoryScale, LinearScale, BarElement, LineElement, PointElement, Title, Tooltip, Legend);

function ProjectDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeTab, setActiveTab] = useState('Tasks');
  const [viewMode, setViewMode] = useState('Kanban');
  const [columns, setColumns] = useState(['To Do', 'In Progress', 'Completed']);
  const [tasks, setTasks] = useState([]);
  const [project, setProject] = useState(null);
  const [undoFeedback, setUndoFeedback] = useState({ show: false, taskId: null });
  const [showAddTaskModal, setShowAddTaskModal] = useState(false);
  const [showAddColumnModal, setShowAddColumnModal] = useState(false);
  const [showChangeOwnerModal, setShowChangeOwnerModal] = useState(false);
  const [newTask, setNewTask] = useState({ title: '', due: '', status: 'To Do', assignee: '' });
  const [newColumnName, setNewColumnName] = useState('');
  const [newColumnPosition, setNewColumnPosition] = useState('end');
  const [newOwnerID, setNewOwnerID] = useState('');
  const [error, setError] = useState('');
  const [analytics, setAnalytics] = useState(null);
  const [recentActivity, setRecentActivity] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  // Fetch project and tasks on mount
  useEffect(() => {
    const fetchProjectData = async () => {
      setIsLoading(true);
      try {
        // Fetch project
        const projectResponse = await api.get(`/projects/${id}`);
        const projectData = projectResponse.data;
        setProject({
          id: projectData.id,
          title: projectData.name,
          description: projectData.description || 'No description provided',
          category: projectData.category || 'Uncategorized',
          status: projectData.status || 'Planned',
          teamMembers: projectData.users.map(user => user.User.name),
          progress: `${projectData.tasks.filter(task => task.status === 'Completed').length} of ${projectData.tasks.length} tasks completed`,
          creator: projectData.creator ? projectData.creator.name : 'Unknown',
        });

        // Fetch tasks
        const tasksResponse = await api.get(`/projects/${id}/tasks`);
        const tasksData = tasksResponse.data.map(task => ({
          id: task.id,
          title: task.title,
          due: new Date(task.due_date).toLocaleDateString(),
          status: task.status,
          assignee: task.user ? task.user.name : 'Unassigned',
          previousStatus: null,
        }));
        setTasks(tasksData);

        // Fetch analytics
        const analyticsResponse = await api.get(`/projects/${id}/analytics`);
        setAnalytics(analyticsResponse.data);

        // Fetch recent activity
        const activityResponse = await api.get(`/projects/${id}/activities`);
        setRecentActivity(activityResponse.data.map(activity => ({
          user: activity.user.name,
          action: activity.action,
          time: new Date(activity.timestamp).toLocaleString(),
        })));

        setError('');
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to fetch project data. Please try again.');
        console.error('Error fetching project data:', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchProjectData();
  }, [id]);

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: { position: 'top' },
      title: { display: true, text: '' },
    },
    scales: {
      y: { beginAtZero: true },
    },
  };

  // Handle drag-and-drop for both tasks and columns
  const onDragEnd = async (result) => {
    const { source, destination, type } = result;

    if (!destination) return;

    if (source.droppableId === destination.droppableId && source.index === destination.index) return;

    if (type === 'column') {
      const newColumns = [...columns];
      const [movedColumn] = newColumns.splice(source.index, 1);
      newColumns.splice(destination.index, 0, movedColumn);
      setColumns(newColumns);
    } else {
      const newTasks = [...tasks];
      const [movedTask] = newTasks.splice(source.index, 1);
      const oldStatus = movedTask.status;
      movedTask.status = destination.droppableId;
      movedTask.previousStatus = oldStatus;
      newTasks.splice(destination.index, 0, movedTask);
      setTasks(newTasks);

      // Update task status in the backend
      try {
        const taskToUpdate = tasks.find(task => task.id === movedTask.id);
        await api.put(`/tasks/${movedTask.id}`, {
          title: taskToUpdate.title,
          description: taskToUpdate.description || '',
          project_id: parseInt(id),
          user_id: project.users?.find(user => user.User.name === taskToUpdate.assignee)?.UserID || 0,
          status: movedTask.status,
          due_date: new Date(taskToUpdate.due).toISOString(),
        });
      } catch (err) {
        setError('Failed to update task status. Please try again.');
        console.error('Error updating task status:', err);
        // Revert the UI change
        movedTask.status = oldStatus;
        movedTask.previousStatus = null;
        setTasks([...newTasks]);
      }
    }
  };

  // Handle marking task as completed
  const markAsCompleted = async (taskId) => {
    const taskToUpdate = tasks.find(task => task.id === taskId);
    if (!taskToUpdate || taskToUpdate.status === 'Completed') return;

    const newTasks = tasks.map(task =>
      task.id === taskId
        ? { ...task, status: 'Completed', previousStatus: task.status }
        : task
    );
    setTasks(newTasks);

    try {
      await api.put(`/tasks/${taskId}`, {
        title: taskToUpdate.title,
        description: taskToUpdate.description || '',
        project_id: parseInt(id),
        user_id: project.users?.find(user => user.User.name === taskToUpdate.assignee)?.UserID || 0,
        status: 'Completed',
        due_date: new Date(taskToUpdate.due).toISOString(),
      });
    } catch (err) {
      setError('Failed to mark task as completed. Please try again.');
      console.error('Error marking task as completed:', err);
      // Revert the UI change
      const revertedTasks = tasks.map(task =>
        task.id === taskId
          ? { ...task, status: taskToUpdate.status, previousStatus: null }
          : task
      );
      setTasks(revertedTasks);
    }
  };

  // Handle undoing completion
  const undoCompletion = async (taskId) => {
    const taskToUndo = tasks.find(task => task.id === taskId);
    if (!taskToUndo || taskToUndo.status !== 'Completed' || !taskToUndo.previousStatus) return;

    if (window.confirm('Are you sure you want to undo this completion?')) {
      const newTasks = tasks.map(task =>
        task.id === taskId
          ? { ...task, status: taskToUndo.previousStatus, previousStatus: null }
          : task
      );
      setTasks(newTasks);

      try {
        await api.put(`/tasks/${taskId}`, {
          title: taskToUndo.title,
          description: taskToUndo.description || '',
          project_id: parseInt(id),
          user_id: project.users?.find(user => user.User.name === taskToUndo.assignee)?.UserID || 0,
          status: taskToUndo.previousStatus,
          due_date: new Date(taskToUndo.due).toISOString(),
        });
        setUndoFeedback({ show: true, taskId });
        setTimeout(() => setUndoFeedback({ show: false, taskId: null }), 2000);
      } catch (err) {
        setError('Failed to undo task completion. Please try again.');
        console.error('Error undoing task completion:', err);
        // Revert the UI change
        const revertedTasks = tasks.map(task =>
          task.id === taskId
            ? { ...task, status: 'Completed', previousStatus: taskToUndo.previousStatus }
            : task
        );
        setTasks(revertedTasks);
      }
    }
  };

  // Handle adding a new task
  const handleAddTask = async (e) => {
    e.preventDefault();
    if (!newTask.title || !newTask.due || !newTask.status || !newTask.assignee) {
      alert('Please fill in all fields.');
      return;
    }

    try {
      const response = await api.post('/tasks', {
        title: newTask.title,
        description: '',
        project_id: parseInt(id),
        user_id: project.users?.find(user => user.User.name === newTask.assignee)?.UserID || 0,
        status: newTask.status,
        due_date: new Date(newTask.due).toISOString(),
      });

      const newTaskData = {
        id: response.data.id,
        title: newTask.title,
        due: new Date(newTask.due).toLocaleDateString(),
        status: newTask.status,
        assignee: newTask.assignee,
        previousStatus: null,
      };

      setTasks([...tasks, newTaskData]);
      setNewTask({ title: '', due: '', status: 'To Do', assignee: '' });
      setShowAddTaskModal(false);
    } catch (err) {
      setError('Failed to add task. Please try again.');
      console.error('Error adding task:', err);
    }
  };

  // Handle deleting a task
  const handleDeleteTask = async (taskId) => {
    if (!window.confirm('Are you sure you want to delete this task?')) return;

    try {
      await api.delete(`/tasks/${taskId}`);
      setTasks(tasks.filter(task => task.id !== taskId));
    } catch (err) {
      setError('Failed to delete task. Please try again.');
      console.error('Error deleting task:', err);
    }
  };

  // Handle adding a new column
  const handleAddColumn = (e) => {
    e.preventDefault();
    if (!newColumnName) {
      alert('Please enter a column name.');
      return;
    }

    if (columns.includes(newColumnName)) {
      alert('Column name already exists.');
      return;
    }

    const newColumns = [...columns];
    if (newColumnPosition === 'end') {
      newColumns.push(newColumnName);
    } else {
      const positionIndex = columns.indexOf(newColumnPosition);
      newColumns.splice(positionIndex + 1, 0, newColumnName);
    }

    setColumns(newColumns);
    setNewColumnName('');
    setNewColumnPosition('end');
    setShowAddColumnModal(false);
  };

  // Handle deleting the project
  const handleDeleteProject = async () => {
    if (!window.confirm('Are you sure you want to delete this project? This action cannot be undone.')) return;

    try {
      await api.delete(`/projects/${id}`);
      navigate('/projects');
    } catch (err) {
      setError('Failed to delete project. Please try again.');
      console.error('Error deleting project:', err);
    }
  };

  // Handle changing the project owner
  const handleChangeOwner = async (e) => {
    e.preventDefault();
    if (!newOwnerID) {
      alert('Please select a new owner.');
      return;
    }

    try {
      await api.put(`/projects/${id}/owner`, { new_owner_id: parseInt(newOwnerID) });
      const projectResponse = await api.get(`/projects/${id}`);
      const projectData = projectResponse.data;
      setProject(prev => ({
        ...prev,
        creator: projectData.creator ? projectData.creator.name : 'Unknown',
      }));
      setShowChangeOwnerModal(false);
      setNewOwnerID('');
    } catch (err) {
      setError('Failed to change project owner. Please try again.');
      console.error('Error changing project owner:', err);
    }
  };

  if (isLoading) {
    return <div className="dashboard">Loading...</div>;
  }

  if (!project) {
    return <div className="dashboard">Project not found.</div>;
  }

  return (
    <div className="dashboard">
      <Sidebar isSidebarCollapsed={isSidebarCollapsed} setIsSidebarCollapsed={setIsSidebarCollapsed} />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="project-details-content">
          <div className="project-header">
            <h1>
              {project.title} <span className="status">{project.status.toLowerCase()}</span>
            </h1>
            <p>{project.description}</p>
          </div>
          <div className="tabs">
            <button
              className={activeTab === 'Tasks' ? 'active' : ''}
              onClick={() => setActiveTab('Tasks')}
            >
              Tasks
            </button>
            <button
              className={activeTab === 'Analytics' ? 'active' : ''}
              onClick={() => setActiveTab('Analytics')}
            >
              Analytics
            </button>
            <button
              className={activeTab === 'Activity' ? 'active' : ''}
              onClick={() => setActiveTab('Activity')}
            >
              Activity
            </button>
            <button
              className={activeTab === 'Files' ? 'active' : ''}
              onClick={() => setActiveTab('Files')}
            >
              Files
            </button>
          </div>
          {error && <p className="error-message">{error}</p>}
          <div className="main-container">
            <div className="main-content">
              {activeTab === 'Tasks' && (
                <div className="tasks-section">
                  <div className="task-header">
                    <div className="task-view-toggle">
                      <button
                        className={viewMode === 'List' ? 'active' : ''}
                        onClick={() => setViewMode('List')}
                      >
                        List
                      </button>
                      <button
                        className={viewMode === 'Kanban' ? 'active' : ''}
                        onClick={() => setViewMode('Kanban')}
                      >
                        Kanban
                      </button>
                      <button onClick={() => setShowAddTaskModal(true)}>+ Add Task</button>
                    </div>
                    <p>{tasks.length} tasks in {columns.length} columns</p>
                  </div>
                  {viewMode === 'Kanban' ? (
                    <DragDropContext onDragEnd={onDragEnd}>
                      <Droppable droppableId="kanban-board" direction="horizontal" type="column">
                        {(provided) => (
                          <div
                            className="kanban-board"
                            ref={provided.innerRef}
                            {...provided.droppableProps}
                          >
                            {columns.map((column, index) => (
                              <Draggable key={column} draggableId={`column-${column}`} index={index}>
                                {(provided) => (
                                  <div
                                    className="column"
                                    ref={provided.innerRef}
                                    {...provided.draggableProps}
                                  >
                                    <h3>
                                      <span className="drag-handle" {...provided.dragHandleProps}>
                                        <GripVertical size={16} />
                                      </span>
                                      <span className={`status-icon ${column.toLowerCase().replace(' ', '-')}`}>
                                        {column === 'To Do' && <Clock size={16} />}
                                        {column === 'In Progress' && <AlertCircle size={16} />}
                                        {column === 'Completed' && <CheckCircle size={16} />}
                                        {column !== 'To Do' && column !== 'In Progress' && column !== 'Completed' && <Clock size={16} />}
                                      </span>
                                      {column}
                                    </h3>
                                    <Droppable droppableId={column} type="task">
                                      {(provided) => (
                                        <div
                                          className="task-list"
                                          ref={provided.innerRef}
                                          {...provided.droppableProps}
                                        >
                                          {tasks
                                            .filter((task) => task.status === column)
                                            .map((task, index) => (
                                              <Draggable
                                                key={task.id}
                                                draggableId={`task-${task.id}`}
                                                index={index}
                                              >
                                                {(provided) => (
                                                  <div
                                                    className={`task-card ${undoFeedback.show && undoFeedback.taskId === task.id ? 'undo-success' : ''}`}
                                                    ref={provided.innerRef}
                                                    {...provided.draggableProps}
                                                    {...provided.dragHandleProps}
                                                  >
                                                    <h4>
                                                      <span className={`status-icon ${column.toLowerCase().replace(' ', '-')}`}>
                                                        {column === 'To Do' && <Clock size={16} />}
                                                        {column === 'In Progress' && <AlertCircle size={16} />}
                                                        {column === 'Completed' && <CheckCircle size={16} />}
                                                        {column !== 'To Do' && column !== 'In Progress' && column !== 'Completed' && <Clock size={16} />}
                                                      </span>
                                                      {task.title}
                                                    </h4>
                                                    <p>Due {task.due}</p>
                                                    <div className="task-assignee">
                                                      <User size={16} />
                                                      <span>{task.assignee}</span>
                                                    </div>
                                                    <div className="task-actions">
                                                      {task.status === 'Completed' ? (
                                                        <button
                                                          className="undo-btn"
                                                          onClick={() => undoCompletion(task.id)}
                                                        >
                                                          Undo
                                                        </button>
                                                      ) : (
                                                        <button className="add-card-btn">+ Add Card</button>
                                                      )}
                                                      <button
                                                        className="delete-btn"
                                                        onClick={() => handleDeleteTask(task.id)}
                                                      >
                                                        <Trash2 size={16} />
                                                      </button>
                                                    </div>
                                                  </div>
                                                )}
                                              </Draggable>
                                            ))}
                                          {provided.placeholder}
                                        </div>
                                      )}
                                    </Droppable>
                                  </div>
                                )}
                              </Draggable>
                            ))}
                            {provided.placeholder}
                          </div>
                        )}
                      </Droppable>
                    </DragDropContext>
                  ) : (
                    <div className="list-view">
                      <table>
                        <thead>
                          <tr>
                            <th>Task</th>
                            <th>Status</th>
                            <th>Due Date</th>
                            <th>Assignee</th>
                            <th>Action</th>
                          </tr>
                        </thead>
                        <tbody>
                          {tasks.map((task) => (
                            <tr key={task.id}>
                              <td>{task.title}</td>
                              <td>
                                <span className={`status-icon ${task.status.toLowerCase().replace(' ', '-')}`}>
                                  {task.status === 'To Do' && <Clock size={16} />}
                                  {task.status === 'In Progress' && <AlertCircle size={16} />}
                                  {task.status === 'Completed' && <CheckCircle size={16} />}
                                  {task.status !== 'To Do' && task.status !== 'In Progress' && task.status !== 'Completed' && <Clock size={16} />}
                                </span>
                                {task.status}
                              </td>
                              <td>{task.due}</td>
                              <td>
                                <div className="task-assignee">
                                  <User size={16} />
                                  <span>{task.assignee}</span>
                                </div>
                              </td>
                              <td>
                                <div className="task-actions">
                                  {task.status !== 'Completed' ? (
                                    <button
                                      className="complete-btn"
                                      onClick={() => markAsCompleted(task.id)}
                                    >
                                      Mark as Completed
                                    </button>
                                  ) : (
                                    <button
                                      className="undo-btn"
                                      onClick={() => undoCompletion(task.id)}
                                    >
                                      Undo
                                    </button>
                                  )}
                                  <button
                                    className="delete-btn"
                                    onClick={() => handleDeleteTask(task.id)}
                                  >
                                    <Trash2 size={16} />
                                  </button>
                                </div>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                      {undoFeedback.show && (
                        <div className="undo-feedback">
                          Task undone successfully!
                        </div>
                      )}
                    </div>
                  )}
                  {viewMode === 'Kanban' && (
                    <button className="add-column-btn" onClick={() => setShowAddColumnModal(true)}>
                      + Add Column
                    </button>
                  )}
                </div>
              )}
              {activeTab === 'Analytics' && analytics && (
                <div className="analytics-section">
                  <div className="chart-row">
                    <div className="chart-container">
                      <Bar data={analytics.cycleTime} options={chartOptions} />
                    </div>
                    <div className="chart-container">
                      <Bar data={analytics.velocity} options={chartOptions} />
                    </div>
                  </div>
                  <div className="chart-row">
                    <div className="chart-container">
                      <Line data={analytics.burndown} options={chartOptions} />
                    </div>
                    <div className="chart-container">
                      <Bar data={analytics.cumulativeFlow} options={chartOptions} />
                    </div>
                  </div>
                </div>
              )}
              {activeTab === 'Activity' && (
                <div className="activity-section">
                  <h2>Recent Activity</h2>
                  {recentActivity.map((activity, index) => (
                    <div key={index} className="activity-item">
                      <div className="activity-avatar">
                        <User size={24} />
                      </div>
                      <div className="activity-details">
                        <p>{activity.user} {activity.action}</p>
                        <span>{activity.time}</span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
              {activeTab === 'Files' && (
                <div className="files-section">
                  <p>File management coming soon</p>
                </div>
              )}
            </div>
            <div className="project-details-sidebar">
              <h3>Project Details</h3>
              <div className="detail-item">
                <span>Category</span>
                <span>{project.category}</span>
              </div>
              <div className="detail-item">
                <span>Owner</span>
                <span>{project.creator}</span>
              </div>
              <div className="detail-item">
                <span>Team Members</span>
                <div className="team-members">
                  {project.teamMembers.map((member, index) => (
                    <span key={index} className="team-member">
                      <User size={16} /> {member}
                    </span>
                  ))}
                </div>
              </div>
              <div className="detail-item">
                <span>Progress</span>
                <span>{project.progress}</span>
              </div>
              {analytics && (
                <div className="detail-item">
                  <span>Key Metrics</span>
                  <div className="metrics">
                    <p>Cycle Time: {analytics.metrics.cycleTime}</p>
                    <p>Velocity: {analytics.metrics.velocity}</p>
                    <p>Defects: {analytics.metrics.defects}</p>
                    <p>Code Coverage: {analytics.metrics.codeCoverage}</p>
                  </div>
                </div>
              )}
              <div className="project-actions">
                <button className="change-owner-btn" onClick={() => setShowChangeOwnerModal(true)}>
                  Change Owner
                </button>
                <button className="delete-project-btn" onClick={handleDeleteProject}>
                  Delete Project
                </button>
              </div>
            </div>
          </div>

          {/* Add Task Modal */}
          {showAddTaskModal && (
            <div className="modal-overlay">
              <div className="modal">
                <div className="modal-header">
                  <h2>Add New Task</h2>
                  <button className="close-btn" onClick={() => setShowAddTaskModal(false)}>
                    <X size={16} />
                  </button>
                </div>
                <form onSubmit={handleAddTask}>
                  <div className="form-group">
                    <label>Title</label>
                    <input
                      type="text"
                      value={newTask.title}
                      onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Due Date</label>
                    <input
                      type="date"
                      value={newTask.due}
                      onChange={(e) => setNewTask({ ...newTask, due: e.target.value })}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Status</label>
                    <select
                      value={newTask.status}
                      onChange={(e) => setNewTask({ ...newTask, status: e.target.value })}
                      required
                    >
                      {columns.map((column) => (
                        <option key={column} value={column}>
                          {column}
                        </option>
                      ))}
                    </select>
                  </div>
                  <div className="form-group">
                    <label>Assignee</label>
                    <select
                      value={newTask.assignee}
                      onChange={(e) => setNewTask({ ...newTask, assignee: e.target.value })}
                      required
                    >
                      <option value="">Select Assignee</option>
                      {project.teamMembers.map((member) => (
                        <option key={member} value={member}>
                          {member}
                        </option>
                      ))}
                    </select>
                  </div>
                  <button type="submit" className="submit-btn">
                    Add Task
                  </button>
                </form>
              </div>
            </div>
          )}

          {/* Add Column Modal */}
          {showAddColumnModal && (
            <div className="modal-overlay">
              <div className="modal">
                <div className="modal-header">
                  <h2>Add New Column</h2>
                  <button className="close-btn" onClick={() => setShowAddColumnModal(false)}>
                    <X size={16} />
                  </button>
                </div>
                <form onSubmit={handleAddColumn}>
                  <div className="form-group">
                    <label>Column Name</label>
                    <input
                      type="text"
                      value={newColumnName}
                      onChange={(e) => setNewColumnName(e.target.value)}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Position After</label>
                    <select
                      value={newColumnPosition}
                      onChange={(e) => setNewColumnPosition(e.target.value)}
                    >
                      <option value="end">End of Board</option>
                      {columns.map((column) => (
                        <option key={column} value={column}>
                          After {column}
                        </option>
                      ))}
                    </select>
                  </div>
                  <button type="submit" className="submit-btn">
                    Add Column
                  </button>
                </form>
              </div>
            </div>
          )}

          {/* Change Owner Modal */}
          {showChangeOwnerModal && (
            <div className="modal-overlay">
              <div className="modal">
                <div className="modal-header">
                  <h2>Change Project Owner</h2>
                  <button className="close-btn" onClick={() => setShowChangeOwnerModal(false)}>
                    <X size={16} />
                  </button>
                </div>
                <form onSubmit={handleChangeOwner}>
                  <div className="form-group">
                    <label>New Owner</label>
                    <select
                      value={newOwnerID}
                      onChange={(e) => setNewOwnerID(e.target.value)}
                      required
                    >
                      <option value="">Select New Owner</option>
                      {project.teamMembers.map((member, index) => (
                        <option key={index} value={project.users[index].UserID}>
                          {member}
                        </option>
                      ))}
                    </select>
                  </div>
                  <button type="submit" className="submit-btn">
                    Change Owner
                  </button>
                </form>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ProjectDetails;