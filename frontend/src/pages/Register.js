import React, { useState } from 'react';
import '../styles/Register.css';

function Register() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [fullName, setFullName] = useState('');


  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('FullName:', fullName,'Email:', email, 'Password:', password);
    // Later, we’ll connect this to backend
  };

  return (
    <div className="login-container">
        <div className="login-box">
            <h1>Optima</h1>
            <h3>Modern Work Management</h3>
            <form onSubmit={handleSubmit}>
                <h4 className="title-form">Full Name</h4>
                <input
                    type='text'
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                    placeholder="Enter your name"
                />
                <h4 className="title-form">Email</h4>
                <input
                    type='email'
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Enter your email"
                />
                <h4 className="title-form">Password</h4>
                <input
                    type='password'
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Enter your name"
                />
                <p className="alert-password">Password must be at least 6 characters</p>
              <button type="submit">Create Account</button>
              <h5 className='create-account'>Already have an account? <a href="/Login">Sign in</a></h5>
            </form>
        </div>
    </div>
  );
}
export default Register;