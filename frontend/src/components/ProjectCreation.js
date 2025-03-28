// src/components/ProjectCreation.js
import React, { useState } from 'react';
import { X } from 'lucide-react';
import '../styles/ProjectCreation.css';

function ProjectCreation({ onAddProject, onClose }) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [teamMembers, setTeamMembers] = useState([]);
  const [status, setStatus] = useState('Planned');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!title || !category) {
      alert('Please fill in all required fields.');
      return;
    }
    onAddProject({ title, description, category, teamMembers, status });
    onClose();
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h2>Create New Project</h2>
          <button className="close-btn" onClick={onClose}>
            <X size={16} />
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>
          <div className="form-group">
            <label>Category</label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              required
            >
              <option value="">Select Category</option>
              <option value="Marketing">Marketing</option>
              <option value="Development">Development</option>
              <option value="Product">Product</option>
              <option value="Research">Research</option>
            </select>
          </div>
          <div className="form-group">
            <label>Team Members</label>
            <input
              type="text"
              value={teamMembers.join(', ')}
              onChange={(e) => setTeamMembers(e.target.value.split(', ').filter(Boolean))}
              placeholder="Enter names separated by commas"
            />
          </div>
          <div className="form-group">
            <label>Status</label>
            <select value={status} onChange={(e) => setStatus(e.target.value)}>
              <option value="Planned">Planned</option>
              <option value="In Progress">In Progress</option>
              <option value="Completed">Completed</option>
              <option value="On Hold">On Hold</option>
            </select>
          </div>
          <button type="submit" className="submit-btn">
            Create Project
          </button>
        </form>
      </div>
    </div>
  );
}

export default ProjectCreation;