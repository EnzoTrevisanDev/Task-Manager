import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom'; // For redirecting after login
import api from '../api/api.js'; // Import the API client
import '../styles/Login.css';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState(''); // For displaying error messages
  const navigate = useNavigate(); // Hook for navigation

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await api.post('/login', { email, password });
      const { access_token, refresh_token } = response.data;
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      navigate('/projects');
    } catch (err) {
      setError(err.response?.data?.error || 'Login failed. Please try again.');
      console.error('Login error:', err);
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>Optima</h1>
        <h3>Modern Work Management</h3>
        {error && <p className="error-message">{error}</p>} {/* Display error message */}
        <form onSubmit={handleLogin}>
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
          <button type="submit">Sign In</button>
          <h5 className="create-account">
            Don’t have an account? <a href="/register">Sign Up</a>
          </h5>
        </form>
      </div>
    </div>
  );
}

export default Login;