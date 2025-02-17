import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import UsersPage from './components/UsersPage';
import BudgetsPage from './components/BudgetsPage';
import ChargesPage from './components/ChargesPage';
import SharesPage from './components/SharesPage';

function App() {
  // Get token/permissions from localStorage so the user stays logged in on refresh.
  const [token, setToken] = useState(localStorage.getItem('token') || '');
  const [permissions, setPermissions] = useState(localStorage.getItem('permissions') || '');

  const handleLogin = (token, permissions) => {
    setToken(token);
    setPermissions(permissions);
    localStorage.setItem('token', token);
    localStorage.setItem('permissions', permissions);
  };

  const handleLogout = () => {
    setToken('');
    setPermissions('');
    localStorage.removeItem('token');
    localStorage.removeItem('permissions');
  };

  // If not logged in, show the Login component
  if (!token) {
    return <Login onLogin={handleLogin} />;
  }

  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li><Link to="/">Dashboard</Link></li>
            {permissions === 'admin' && <li><Link to="/users">Users</Link></li>}
            <li><Link to="/budgets">Budgets</Link></li>
            <li><Link to="/charges">Charges</Link></li>
            <li><Link to="/shares">Shares</Link></li>
            <li><button onClick={handleLogout}>Logout</button></li>
          </ul>
        </nav>
        <hr />
        <Routes>
          <Route path="/" element={<Dashboard />} />
          {permissions === 'admin' && <Route path="/users" element={<UsersPage token={token} />} />}
          <Route path="/budgets" element={<BudgetsPage token={token} />} />
          <Route path="/charges" element={<ChargesPage token={token} />} />
          <Route path="/shares" element={<SharesPage token={token} />} />
          {/* Redirect unknown routes */}
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
