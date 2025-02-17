import React, { useState, useEffect } from 'react';

function SharesPage({ token }) {
  const [shares, setShares] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchShares = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/shares', {
          headers: { 'Authorization': `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          setShares(data);
        } else {
          setError('Failed to fetch shares');
        }
      } catch (err) {
        setError('Error fetching shares');
      }
    };

    fetchShares();
  }, [token]);

  return (
    <div>
      <h2>Your Shares</h2>
      {error && <p style={{ color:'red' }}>{error}</p>}
      <ul>
        {shares.map(share => (
          <li key={share.id}>
            Share ID: {share.id} — Access: {share.access}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default SharesPage;
