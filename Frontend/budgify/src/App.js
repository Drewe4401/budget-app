import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import UsersPage from './components/UsersPage';
import BudgetsPage from './components/BudgetsPage';
import ChargesPage from './components/ChargesPage';
import SharesPage from './components/SharesPage';

function App() {
  // Get token/permissions from localStorage so the user stays logged in on refresh
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
      {/* Main container for sidebar layout */}
      <div className="flex min-h-screen bg-gray-100">
        {/* Sidebar */}
        <nav className="w-64 bg-gray-800 text-white p-4 flex flex-col">
          <h1 className="text-2xl font-bold mb-6">Budgify</h1>
          <ul className="flex-grow">
            <li className="mb-2">
              <Link
                to="/"
                className="block px-4 py-2 rounded hover:bg-gray-700 transition-colors"
              >
                Dashboard
              </Link>
            </li>
            {permissions === 'admin' && (
              <li className="mb-2">
                <Link
                  to="/users"
                  className="block px-4 py-2 rounded hover:bg-gray-700 transition-colors"
                >
                  Users
                </Link>
              </li>
            )}
            <li className="mb-2">
              <Link
                to="/budgets"
                className="block px-4 py-2 rounded hover:bg-gray-700 transition-colors"
              >
                Budgets
              </Link>
            </li>
            <li className="mb-2">
              <Link
                to="/charges"
                className="block px-4 py-2 rounded hover:bg-gray-700 transition-colors"
              >
                Charges
              </Link>
            </li>
            <li className="mb-2">
              <Link
                to="/shares"
                className="block px-4 py-2 rounded hover:bg-gray-700 transition-colors"
              >
                Shares
              </Link>
            </li>
          </ul>
          <button
            onClick={handleLogout}
            className="bg-red-500 hover:bg-red-600 text-white rounded px-4 py-2 transition-colors"
          >
            Logout
          </button>
        </nav>

        {/* Content area */}
        <div className="flex-1 p-6">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            {permissions === 'admin' && (
              <Route path="/users" element={<UsersPage token={token} />} />
            )}
            <Route path="/budgets" element={<BudgetsPage token={token} />} />
            <Route path="/charges" element={<ChargesPage token={token} />} />
            <Route path="/shares" element={<SharesPage token={token} />} />
            {/* Redirect unknown routes */}
            <Route path="*" element={<Navigate to="/" />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
