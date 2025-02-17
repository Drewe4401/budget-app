import React, { useState, useEffect } from 'react';

function BudgetsPage({ token }) {
  const [budgets, setBudgets] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchBudgets = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/budgets', {
          headers: { 'Authorization': `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          setBudgets(data);
        } else {
          setError('Failed to fetch budgets');
        }
      } catch (err) {
        setError('Error fetching budgets');
      }
    };

    fetchBudgets();
  }, [token]);

  return (
    <div>
      <h2>Your Budgets</h2>
      {error && <p style={{ color:'red' }}>{error}</p>}
      <ul>
        {budgets.map(budget => (
          <li key={budget.id}>
            {budget.name}: ${budget.amount}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default BudgetsPage;
