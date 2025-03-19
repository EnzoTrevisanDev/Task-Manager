import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { AlertCircle, Clock, CheckCircle, User } from 'lucide-react';
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
import '../styles/ProjectDetails.css';

// Register ChartJS components
ChartJS.register(CategoryScale, LinearScale, BarElement, LineElement, PointElement, Title, Tooltip, Legend);

function ProjectDetails() {
  const { id } = useParams();
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeTab, setActiveTab] = useState('Tasks');
  const [viewMode, setViewMode] = useState('Kanban');
  const [tasks, setTasks] = useState([
    { id: 1, title: 'Component documentation', due: '3/24/2024', status: 'To Do', assignee: 'Emily Brown', previousStatus: null },
    { id: 2, title: 'Design system implementation', due: '3/19/2024', status: 'In Progress', assignee: 'Sarah Wilson', previousStatus: null },
    { id: 3, title: 'User research interviews', due: '3/14/2024', status: 'Completed', assignee: 'James Miller', previousStatus: 'In Progress' },
  ]);
  const [undoFeedback, setUndoFeedback] = useState({ show: false, taskId: null }); // State for visual feedback

  // Mock project data (excluding tasks, handled separately)
  const project = {
    id: 1,
    title: 'Marketing Campaign',
    description: 'Q1 2024 Digital Marketing Initiative',
    category: 'Marketing',
    status: 'In Progress',
    teamMembers: ['Sarah Wilson', 'James Miller', 'Michael Chen'],
    progress: '3 of 3 tasks completed',
    metrics: {
      cycleTime: '4.2 days',
      velocity: '28 points',
      defects: 3,
      codeCoverage: '87%',
    },
  };

  // Chart data
  const cycleTimeData = {
    labels: ['Task 1', 'Task 2', 'Task 3', 'Task 4', 'Task 5'],
    datasets: [
      {
        label: 'Cycle Time (Days)',
        data: [5, 7, 3, 6, 4],
        backgroundColor: '#9b87f6',
      },
    ],
  };

  const velocityData = {
    labels: ['Sprint 1', 'Sprint 2', 'Sprint 3', 'Sprint 4'],
    datasets: [
      {
        label: 'Velocity (Story Points)',
        data: [18, 27, 36, 30],
        backgroundColor: '#2ecc71',
      },
    ],
  };

  const burndownData = {
    labels: ['Day 1', 'Day 2', 'Day 3', 'Day 4', 'Day 5', 'Day 6', 'Day 7'],
    datasets: [
      {
        label: 'Burndown',
        data: [12, 10, 8, 6, 4, 2, 0],
        fill: false,
        borderColor: '#9b87f6',
        tension: 0.1,
      },
    ],
  };

  const cumulativeFlowData = {
    labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5'],
    datasets: [
      {
        label: 'To Do',
        data: [15, 12, 10, 8, 5],
        backgroundColor: '#e74c3c',
      },
      {
        label: 'In Progress',
        data: [0, 3, 5, 7, 8],
        backgroundColor: '#f1c40f',
      },
      {
        label: 'Completed',
        data: [0, 0, 0, 0, 2],
        backgroundColor: '#2ecc71',
      },
    ],
  };

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

  // Mock recent activity
  const recentActivity = [
    { user: 'Sarah Wilson', action: 'created Marketing Campaign', time: '30 minutes ago' },
    { user: 'James Miller', action: 'completed User research interviews', time: '2 hours ago' },
    { user: 'Sarah Wilson', action: 'updated Project deadline', time: '1 day ago' },
    { user: 'Michael Chen', action: 'started Design system implementation', time: '3 days ago' },
  ];

  // Handle drag-and-drop
  const onDragEnd = (result) => {
    const { source, destination } = result;

    if (!destination) return;

    if (source.droppableId === destination.droppableId && source.index === destination.index) return;

    const newTasks = [...tasks];
    const [movedTask] = newTasks.splice(source.index, 1);
    const oldStatus = movedTask.status;
    movedTask.status = destination.droppableId; // Update status based on the column
    movedTask.previousStatus = oldStatus; // Store the previous status
    newTasks.splice(destination.index, 0, movedTask);

    setTasks(newTasks);
  };

  // Handle marking task as completed in List view
  const markAsCompleted = (taskId) => {
    const newTasks = tasks.map((task) =>
      task.id === taskId && task.status !== 'Completed'
        ? { ...task, status: 'Completed', previousStatus: task.status }
        : task
    );
    setTasks(newTasks);
  };

  // Handle undoing completion with restrictions and confirmation
  const undoCompletion = (taskId) => {
    const taskToUndo = tasks.find((task) => task.id === taskId);
    if (!taskToUndo || taskToUndo.status !== 'Completed' || !taskToUndo.previousStatus) {
      return; // Undo restriction: only allow if completed and has previous status
    }

    if (window.confirm('Are you sure you want to undo this completion?')) {
      const newTasks = tasks.map((task) =>
        task.id === taskId
          ? { ...task, status: task.previousStatus, previousStatus: null }
          : task
      );
      setTasks(newTasks);
      setUndoFeedback({ show: true, taskId }); // Trigger visual feedback
      setTimeout(() => setUndoFeedback({ show: false, taskId: null }), 2000); // Hide after 2 seconds
    }
  };

  return (
    <div className="dashboard">
      <Sidebar isSidebarCollapsed={isSidebarCollapsed} setIsSidebarCollapsed={setIsSidebarCollapsed} />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="project-details-content">
          <div className="project-header">
            <h1>
              {project.title} <span className="status">in-progress</span>
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
                      <button>+ Add Task</button>
                    </div>
                    <p>{tasks.length} tasks in 3 columns</p>
                  </div>
                  {viewMode === 'Kanban' ? (
                    <DragDropContext onDragEnd={onDragEnd}>
                      <div className="kanban-board">
                        <Droppable droppableId="To Do">
                          {(provided) => (
                            <div
                              className="column"
                              ref={provided.innerRef}
                              {...provided.droppableProps}
                            >
                              <h3>
                                <span className="status-icon todo">
                                  <Clock size={16} />
                                </span>
                                To Do
                              </h3>
                              {tasks
                                .filter((task) => task.status === 'To Do')
                                .map((task, index) => (
                                  <Draggable
                                    key={task.id}
                                    draggableId={`task-${task.id}`}
                                    index={index}
                                  >
                                    {(provided) => (
                                      <div
                                        className="task-card"
                                        ref={provided.innerRef}
                                        {...provided.draggableProps}
                                        {...provided.dragHandleProps}
                                      >
                                        <h4>
                                          <span className="status-icon todo">
                                            <Clock size={16} />
                                          </span>
                                          {task.title}
                                        </h4>
                                        <p>Due {task.due}</p>
                                        <div className="task-assignee">
                                          <User size={16} />
                                          <span>{task.assignee}</span>
                                        </div>
                                        <button className="add-card-btn">+ Add Card</button>
                                      </div>
                                    )}
                                  </Draggable>
                                ))}
                              {provided.placeholder}
                            </div>
                          )}
                        </Droppable>
                        <Droppable droppableId="In Progress">
                          {(provided) => (
                            <div
                              className="column"
                              ref={provided.innerRef}
                              {...provided.droppableProps}
                            >
                              <h3>
                                <span className="status-icon in-progress">
                                  <AlertCircle size={16} />
                                </span>
                                In Progress
                              </h3>
                              {tasks
                                .filter((task) => task.status === 'In Progress')
                                .map((task, index) => (
                                  <Draggable
                                    key={task.id}
                                    draggableId={`task-${task.id}`}
                                    index={index}
                                  >
                                    {(provided) => (
                                      <div
                                        className="task-card"
                                        ref={provided.innerRef}
                                        {...provided.draggableProps}
                                        {...provided.dragHandleProps}
                                      >
                                        <h4>
                                          <span className="status-icon in-progress">
                                            <AlertCircle size={16} />
                                          </span>
                                          {task.title}
                                        </h4>
                                        <p>Due {task.due}</p>
                                        <div className="task-assignee">
                                          <User size={16} />
                                          <span>{task.assignee}</span>
                                        </div>
                                        <button className="add-card-btn">+ Add Card</button>
                                      </div>
                                    )}
                                  </Draggable>
                                ))}
                              {provided.placeholder}
                            </div>
                          )}
                        </Droppable>
                        <Droppable droppableId="Completed">
                          {(provided) => (
                            <div
                              className="column"
                              ref={provided.innerRef}
                              {...provided.droppableProps}
                            >
                              <h3>
                                <span className="status-icon completed">
                                  <CheckCircle size={16} />
                                </span>
                                Completed
                              </h3>
                              {tasks
                                .filter((task) => task.status === 'Completed')
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
                                          <span className="status-icon completed">
                                            <CheckCircle size={16} />
                                          </span>
                                          {task.title}
                                        </h4>
                                        <p>Due {task.due}</p>
                                        <div className="task-assignee">
                                          <User size={16} />
                                          <span>{task.assignee}</span>
                                        </div>
                                        <button
                                          className="undo-btn"
                                          onClick={() => undoCompletion(task.id)}
                                        >
                                          Undo
                                        </button>
                                      </div>
                                    )}
                                  </Draggable>
                                ))}
                              {provided.placeholder}
                            </div>
                          )}
                        </Droppable>
                      </div>
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
                    <button className="add-column-btn">+ Add Column</button>
                  )}
                </div>
              )}
              {activeTab === 'Analytics' && (
                <div className="analytics-section">
                  <div className="chart-row">
                    <div className="chart-container">
                      <Bar data={cycleTimeData} options={chartOptions} />
                    </div>
                    <div className="chart-container">
                      <Bar data={velocityData} options={chartOptions} />
                    </div>
                  </div>
                  <div className="chart-row">
                    <div className="chart-container">
                      <Line data={burndownData} options={chartOptions} />
                    </div>
                    <div className="chart-container">
                      <Bar data={cumulativeFlowData} options={chartOptions} />
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
              <div className="detail-item">
                <span>Key Metrics</span>
                <div className="metrics">
                  <p>Cycle Time: {project.metrics.cycleTime}</p>
                  <p>Velocity: {project.metrics.velocity}</p>
                  <p>Defects: {project.metrics.defects}</p>
                  <p>Code Coverage: {project.metrics.codeCoverage}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default ProjectDetails;