import React, { useState } from 'react';
import '../styles/Login.css';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('Email:', email, 'Password:', password);
    // Later, we’ll connect this to backend
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>Optima</h1>
        <h3>Modern Work Management</h3>
        <form onSubmit={handleSubmit}>
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
          <h5 className='create-account'>Don’t have an account? <a href="/register">Sign Up</a></h5>
        </form>
      </div>
    </div>
  );
}

export default Login;