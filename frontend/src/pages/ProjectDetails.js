import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { AlertCircle, Clock, CheckCircle, User, X, GripVertical } from 'lucide-react';
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
  const [columns, setColumns] = useState(['To Do', 'In Progress', 'Completed']);
  const [tasks, setTasks] = useState([
    { id: 1, title: 'Component documentation', due: '3/24/2024', status: 'To Do', assignee: 'Emily Brown', previousStatus: null },
    { id: 2, title: 'Design system implementation', due: '3/19/2024', status: 'In Progress', assignee: 'Sarah Wilson', previousStatus: null },
    { id: 3, title: 'User research interviews', due: '3/14/2024', status: 'Completed', assignee: 'James Miller', previousStatus: 'In Progress' },
  ]);
  const [undoFeedback, setUndoFeedback] = useState({ show: false, taskId: null });
  const [showAddTaskModal, setShowAddTaskModal] = useState(false);
  const [showAddColumnModal, setShowAddColumnModal] = useState(false);
  const [newTask, setNewTask] = useState({ title: '', due: '', status: 'To Do', assignee: '' });
  const [newColumnName, setNewColumnName] = useState('');
  const [newColumnPosition, setNewColumnPosition] = useState('end'); // Default to adding at the end

  // Mock project data
  const project = {
    id: 1,
    title: 'Marketing Campaign',
    description: 'Q1 2024 Digital Marketing Initiative',
    category: 'Marketing',
    status: 'In Progress',
    teamMembers: ['Sarah Wilson', 'James Miller', 'Michael Chen', 'Emily Brown'],
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

  // Handle drag-and-drop for both tasks and columns
  const onDragEnd = (result) => {
    const { source, destination, type } = result;

    if (!destination) return;

    if (source.droppableId === destination.droppableId && source.index === destination.index) return;

    if (type === 'column') {
      // Handle column reordering
      const newColumns = [...columns];
      const [movedColumn] = newColumns.splice(source.index, 1);
      newColumns.splice(destination.index, 0, movedColumn);
      setColumns(newColumns);
    } else {
      // Handle task reordering
      const newTasks = [...tasks];
      const [movedTask] = newTasks.splice(source.index, 1);
      const oldStatus = movedTask.status;
      movedTask.status = destination.droppableId;
      movedTask.previousStatus = oldStatus;
      newTasks.splice(destination.index, 0, movedTask);
      setTasks(newTasks);
    }
  };

  // Handle marking task as completed
  const markAsCompleted = (taskId) => {
    const newTasks = tasks.map((task) =>
      task.id === taskId && task.status !== 'Completed'
        ? { ...task, status: 'Completed', previousStatus: task.status }
        : task
    );
    setTasks(newTasks);
  };

  // Handle undoing completion
  const undoCompletion = (taskId) => {
    const taskToUndo = tasks.find((task) => task.id === taskId);
    if (!taskToUndo || taskToUndo.status !== 'Completed' || !taskToUndo.previousStatus) {
      return;
    }

    if (window.confirm('Are you sure you want to undo this completion?')) {
      const newTasks = tasks.map((task) =>
        task.id === taskId
          ? { ...task, status: task.previousStatus, previousStatus: null }
          : task
      );
      setTasks(newTasks);
      setUndoFeedback({ show: true, taskId });
      setTimeout(() => setUndoFeedback({ show: false, taskId: null }), 2000);
    }
  };

  // Handle adding a new task
  const handleAddTask = (e) => {
    e.preventDefault();
    if (!newTask.title || !newTask.due || !newTask.status || !newTask.assignee) {
      alert('Please fill in all fields.');
      return;
    }

    const newTaskData = {
      id: tasks.length + 1,
      title: newTask.title,
      due: newTask.due,
      status: newTask.status,
      assignee: newTask.assignee,
      previousStatus: null,
    };

    setTasks([...tasks, newTaskData]);
    setNewTask({ title: '', due: '', status: 'To Do', assignee: '' });
    setShowAddTaskModal(false);
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
                    <button className="add-column-btn" onClick={() => setShowAddColumnModal(true)}>
                      + Add Column
                    </button>
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
        </div>
      </div>
    </div>
  );
}

export default ProjectDetails;