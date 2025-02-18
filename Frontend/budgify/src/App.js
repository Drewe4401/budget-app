import React, { useState } from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
  Link,
} from 'react-router-dom';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import UsersPage from './components/UsersPage';
import BudgetsPage from './components/BudgetsPage';
import ChargesPage from './components/ChargesPage';
import SharesPage from './components/SharesPage';

// Extract the routes into a separate component to avoid duplication
function AppRoutes({ token, permissions }) {
  return (
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
  );
}

function App() {
  // Get token/permissions from localStorage so the user stays logged in on refresh
  const [token, setToken] = useState(localStorage.getItem('token') || '');
  const [permissions, setPermissions] = useState(
    localStorage.getItem('permissions') || ''
  );

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
      {/* Desktop Layout (visible on md and up) */}
      <div className="hidden md:flex min-h-screen bg-gray-100">
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
          <AppRoutes token={token} permissions={permissions} />
        </div>
      </div>

      {/* Mobile Layout (visible below md) */}
      <div className="md:hidden flex flex-col min-h-screen">
        {/* Fixed header */}
        <header className="fixed top-0 left-0 right-0 bg-gray-800 text-white p-4 flex justify-between items-center z-10">
          <h1 className="text-lg font-bold">Budgify</h1>
          <button
            onClick={handleLogout}
            className="bg-red-500 hover:bg-red-600 text-white rounded px-3 py-1"
          >
            Logout
          </button>
        </header>

        {/* Main content with top and bottom padding to avoid overlap */}
        <main className="flex-1 p-4 mt-16 mb-16">
          <AppRoutes token={token} permissions={permissions} />
        </main>

        {/* Fixed bottom navigation */}
        <nav className="fixed bottom-0 left-0 right-0 bg-gray-800 text-white flex justify-around items-center p-2">
          <Link to="/" className="flex flex-col items-center text-xs">
            Dashboard
          </Link>
          {permissions === 'admin' && (
            <Link to="/users" className="flex flex-col items-center text-xs">
              Users
            </Link>
          )}
          <Link to="/budgets" className="flex flex-col items-center text-xs">
            Budgets
          </Link>
          <Link to="/charges" className="flex flex-col items-center text-xs">
            Charges
          </Link>
          <Link to="/shares" className="flex flex-col items-center text-xs">
            Shares
          </Link>
        </nav>
      </div>
    </Router>
  );
}

export default App;
