import React, { useState } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar, Line, Pie } from 'react-chartjs-2';
import Sidebar from '../components/Sidebar';
import TopBar from '../components/TopBar';
import '../styles/AnalyticsDashboard.css';

// Register ChartJS components
ChartJS.register(CategoryScale, LinearScale, BarElement, LineElement, ArcElement, Title, Tooltip, Legend);

function AnalyticsDashboard() {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  // Mock data for metrics
  const metrics = {
    totalProjects: 6,
    totalProjectsChange: '+12% vs last month',
    totalTasks: 105,
    totalTasksChange: '-8% vs last month',
    averageProgress: '46.7%',
    averageProgressChange: '+1% vs last month',
    upcomingDeadlines: 2,
    upcomingDeadlinesChange: '-1 vs last month',
  };

  // Data for Projects by Status (Pie Chart)
  const projectsByStatusData = {
    labels: ['Planned', 'In Progress', 'Completed', 'On Hold'],
    datasets: [
      {
        label: 'Projects by Status',
        data: [33, 33, 17, 17], // Percentages
        backgroundColor: ['#9b87f6', '#f1c40f', '#2ecc71', '#e74c3c'],
        borderWidth: 0,
      },
    ],
  };

  // Data for Projects by Category (Bar Chart)
  const projectsByCategoryData = {
    labels: ['Design', 'Development', 'Marketing', 'HR'],
    datasets: [
      {
        label: 'Projects by Category',
        data: [2, 3, 3, 1],
        backgroundColor: '#9b87f6',
        borderWidth: 0,
      },
    ],
  };

  // Data for Project Progress (Line Chart)
  const projectProgressData = {
    labels: [
      'Web Redesign',
      'Mobile App',
      'Content Marketing Campaign',
      'CRM Integration',
      'Product Launch',
      'Employee Training Program',
    ],
    datasets: [
      {
        label: 'Project Progress',
        data: [80, 60, 40, 20, 10, 0],
        fill: false,
        borderColor: '#9b87f6',
        tension: 0.4,
      },
    ],
  };

  // Budget Overview Data
  const budget = {
    utilization: 49,
    total: '$163,000',
    spent: '$89,300',
    remaining: '$93,700',
  };

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: { position: 'bottom' },
      title: { display: false },
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 100, // For progress chart
      },
    },
  };

  const barChartOptions = {
    ...chartOptions,
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  return (
    <div className="dashboard">
      <Sidebar isSidebarCollapsed={isSidebarCollapsed} setIsSidebarCollapsed={setIsSidebarCollapsed} />
      <div className={`dashboard-content ${isSidebarCollapsed ? 'sidebar-collapsed' : ''}`}>
        <TopBar isSidebarCollapsed={isSidebarCollapsed} />
        <div className="analytics-dashboard-content">
          <div className="dashboard-header">
            <h1>Analytics Dashboard</h1>
            <p>Track your project performance and key metrics.</p>
          </div>

          {/* Metrics Cards */}
          <div className="metrics-cards">
            <div className="metric-card">
              <h3>Total Projects</h3>
              <p className="metric-value">{metrics.totalProjects}</p>
              <p className="metric-change positive">{metrics.totalProjectsChange}</p>
            </div>
            <div className="metric-card">
              <h3>Tasks</h3>
              <p className="metric-value">{metrics.totalTasks}</p>
              <p className="metric-change negative">{metrics.totalTasksChange}</p>
            </div>
            <div className="metric-card">
              <h3>Average Progress</h3>
              <p className="metric-value">{metrics.averageProgress}</p>
              <p className="metric-change positive">{metrics.averageProgressChange}</p>
            </div>
            <div className="metric-card">
              <h3>Upcoming Deadlines</h3>
              <p className="metric-value">{metrics.upcomingDeadlines}</p>
              <p className="metric-change negative">{metrics.upcomingDeadlinesChange}</p>
            </div>
          </div>

          {/* Charts Section */}
          <div className="charts-section">
            {/* Projects by Status (Pie Chart) */}
            <div className="chart-container">
              <h2>Projects by Status</h2>
              <div className="pie-chart-wrapper">
                <Pie data={projectsByStatusData} options={chartOptions} />
              </div>
            </div>

            {/* Projects by Category (Bar Chart) */}
            <div className="chart-container">
              <h2>Projects by Category</h2>
              <Bar data={projectsByCategoryData} options={barChartOptions} />
            </div>

            {/* Project Progress (Line Chart) */}
            <div className="chart-container full-width">
              <h2>Project Progress</h2>
              <Line data={projectProgressData} options={chartOptions} />
            </div>

            {/* Budget Overview */}
            <div className="budget-overview">
              <h2>Budget Overview</h2>
              <div className="budget-details">
                <div className="budget-utilization">
                  <p>Budget Utilization</p>
                  <div className="progress-bar">
                    <div
                      className="progress-fill"
                      style={{ width: `${budget.utilization}%` }}
                    ></div>
                  </div>
                  <p>{budget.utilization}%</p>
                </div>
                <div className="budget-stats">
                  <div>
                    <p>Total</p>
                    <p className="budget-value">{budget.total}</p>
                  </div>
                  <div>
                    <p>Spent</p>
                    <p className="budget-value">{budget.spent}</p>
                  </div>
                  <div>
                    <p>Remaining</p>
                    <p className="budget-value">{budget.remaining}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AnalyticsDashboard;