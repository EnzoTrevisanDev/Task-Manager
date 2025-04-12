// src/components/ProjectCreation.js
import React, { useState, useEffect } from 'react';
import { X } from 'lucide-react';
import '../styles/ProjectCreation.css';
import api from '../api/api';
import UserSelect from './UserSelect';

function ProjectCreation({ onAddProject, onClose }) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [teamMembers, setTeamMembers] = useState([]);
  const [status, setStatus] = useState('Planned');
  const [users, setUsers] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await api.get('/users');
        console.log('Raw users data from API:', response.data);
        
        // Transform the users data to match the expected format
        const transformedUsers = response.data.map(user => {
          console.log('Processing user:', user);
          // Check if the user data is nested under User property
          const userData = user.User || user;
          return {
            id: userData.id || userData.ID,
            name: userData.name || userData.Name || 'Unknown User'
          };
        }).filter(user => user.id); // Filter out any invalid users
        
        console.log('Transformed users:', transformedUsers);
        setUsers(transformedUsers);
      } catch (err) {
        console.error('Error fetching users:', err);
        setError('Failed to fetch users. Please try again.');
      }
    };
    fetchUsers();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title || !category) {
      setError('Please fill in all required fields.');
      return;
    }

    try {
      // Create the project
      const projectResponse = await api.post('/projects', {
        name: title,
        description,
        category,
        status,
      });

      console.log('Project created:', projectResponse.data);

      // Add team members to the project
      if (teamMembers.length > 0) {
        try {
          await Promise.all(teamMembers.map(async (userId) => {
            console.log('Adding user to project:', userId);
            await api.post(`/projects/${projectResponse.data.id}/users`, {
              user_id: parseInt(userId),
              role: 'viewer',
            });
          }));
        } catch (err) {
          console.error('Error adding team members:', err);
          // Don't throw here, just log the error
          // The project was created successfully, so we can still close the modal
        }
      }

      onAddProject(projectResponse.data);
      onClose();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create project. Please try again.');
      console.error('Error creating project:', err);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h2>Create New Project</h2>
          <button className="close-button" onClick={onClose}>
            <X size={16} />
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          {error && <div className="error-message">{error}</div>}
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
              <option value="">Select a category</option>
              <option value="Development">Development</option>
              <option value="Design">Design</option>
              <option value="Marketing">Marketing</option>
              <option value="Sales">Sales</option>
              <option value="Other">Other</option>
            </select>
          </div>
          <div className="form-group">
            <label>Team Members</label>
            {users.length === 0 ? (
              <div style={{ color: '#9ca3af', fontSize: '14px' }}>Loading users...</div>
            ) : (
              <>
                <UserSelect
                  users={users}
                  value={teamMembers}
                  onChange={(e) => {
                    console.log('Selected options:', Array.from(e.target.selectedOptions, option => option.value));
                    setTeamMembers(Array.from(e.target.selectedOptions, option => option.value));
                  }}
                  placeholder="Select team members"
                  multiple={true}
                />
                <small style={{ color: '#9ca3af', fontSize: '12px', display: 'block', marginTop: '5px' }}>
                  Hold Ctrl/Cmd to select multiple members
                </small>
              </>
            )}
          </div>
          <div className="form-group">
            <label>Status</label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value)}
            >
              <option value="Planned">Planned</option>
              <option value="In Progress">In Progress</option>
              <option value="Completed">Completed</option>
              <option value="On Hold">On Hold</option>
            </select>
          </div>
          <div className="modal-footer">
            <button type="submit" className="submit-button">
              Create Project
            </button>
            <button type="button" onClick={onClose} className="cancel-button">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default ProjectCreation;