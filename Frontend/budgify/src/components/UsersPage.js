import React, { useState, useEffect } from 'react';

function UsersPage({ token }) {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/users', {
          headers: { 'Authorization': `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          setUsers(data);
        } else {
          setError('Failed to fetch users');
        }
      } catch (err) {
        setError('Error fetching users');
      }
    };

    fetchUsers();
  }, [token]);

  return (
    <div>
      <h2>Users (Admin Only)</h2>
      {error && <p style={{ color:'red' }}>{error}</p>}
      <ul>
        {users.map(user => (
          <li key={user.id}>
            {user.username} - {user.permissions}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default UsersPage;
