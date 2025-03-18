import React from 'react';

function Navbar(){
    return (
        <nav>
        <ul>
            <li><a href="/">Home</a></li>
            <li><a href="/login">Login</a></li>
            <li><a href="/dashboard">Dashboard</a></li>
        </ul>
    </nav>
  );
}
export default Navbar;