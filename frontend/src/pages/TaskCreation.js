// TaskCreation.js
import React, { useState } from 'react';
import { X } from 'lucide-react';
import '../styles/TaskCreation.css';

function TaskCreation({ onAddTask, columns, teamMembers, onClose }) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [dueDate, setDueDate] = useState('');
  const [assignee, setAssignee] = useState('');
  const [status, setStatus] = useState(columns[0]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!title || !dueDate || !assignee || !status) {
      alert('Please fill in all required fields.');
      return;
    }
    onAddTask({ title, description, dueDate, assignee, status });
    onClose();
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h2>Add New Task</h2>
          <button className="close-btn" onClick={onClose}><X size={16} /></button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Title</label>
            <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <div className="form-group">
            <label>Due Date</label>
            <input type="date" value={dueDate} onChange={(e) => setDueDate(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Assignee</label>
            <select value={assignee} onChange={(e) => setAssignee(e.target.value)} required>
              <option value="">Select Assignee</option>
              {teamMembers.map((member) => (
                <option key={member} value={member}>{member}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Status</label>
            <select value={status} onChange={(e) => setStatus(e.target.value)} required>
              {columns.map((column) => (
                <option key={column} value={column}>{column}</option>
              ))}
            </select>
          </div>
          <button type="submit" className="submit-btn">Add Task</button>
        </form>
      </div>
    </div>
  );
}

export default TaskCreation;