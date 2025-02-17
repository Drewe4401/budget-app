import React, { useState, useEffect } from 'react';

function ChargesPage({ token }) {
  const [charges, setCharges] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchCharges = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/charges', {
          headers: { 'Authorization': `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          setCharges(data);
        } else {
          setError('Failed to fetch charges');
        }
      } catch (err) {
        setError('Error fetching charges');
      }
    };

    fetchCharges();
  }, [token]);

  return (
    <div>
      <h2>Your Charges</h2>
      {error && <p style={{ color:'red' }}>{error}</p>}
      <ul>
        {charges.map(charge => (
          <li key={charge.id}>
            {charge.name}: ${charge.amount} ({charge.charge_type})
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ChargesPage;
