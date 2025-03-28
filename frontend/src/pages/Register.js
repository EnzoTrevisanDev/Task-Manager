import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom'; // For redirecting after registration
import api from '../api/api.js'; // Import the API client
import '../styles/Register.css';

function Register() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [fullName, setFullName] = useState('');
  const [error, setError] = useState(''); // For displaying error messages
  const navigate = useNavigate(); // Hook for navigation

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(''); // Clear any previous errors

    if (!fullName || !email || !password) {
      setError('Please fill in all fields.');
      return;
    }
    if (password.length < 6) {
      setError('Password must be at least 6 characters.');
      return;
    }

    try {
      await api.post('/users', {
        name: fullName,
        email,
        password,
      });
      // On successful registration, redirect to login page
      navigate('/login');
    } catch (err) {
      // Handle errors (e.g., email already exists)
      setError(err.response?.data?.error || 'Registration failed. Please try again.');
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>Optima</h1>
        <h3>Modern Work Management</h3>
        {error && <p className="error-message">{error}</p>} {/* Display error message */}
        <form onSubmit={handleSubmit}>
          <h4 className="title-form">Full Name</h4>
          <input
            type="text"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            placeholder="Enter your name"
          />
          <h4 className="title-form">Email</h4>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter your email"
          />
          <h4 className="title-form">Password</h4>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Enter your password"
          />
          <p className="alert-password">Password must be at least 6 characters</p>
          <button type="submit">Create Account</button>
          <h5 className="create-account">
            Already have an account? <a href="/login">Sign in</a>
          </h5>
        </form>
      </div>
    </div>
  );
}

export default Register;