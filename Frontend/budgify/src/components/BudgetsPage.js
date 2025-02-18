import React, { useState, useEffect } from 'react';

// Comprehensive list of expense categories:
const categoryOptions = [
  'Rent',
  'Mortgage',
  'Utilities',
  'Internet',
  'Phone',
  'Groceries',
  'Dining Out',
  'Shopping',
  'Clothing',
  'Transportation',
  'Car Payment',
  'Fuel',
  'Insurance',
  'Medical',
  'Health & Fitness',
  'Entertainment',
  'Travel',
  'Education',
  'Childcare',
  'Savings',
  'Debt Payments',
  'Credit Card Payments',
  'Home Maintenance',
  'Pet Expenses',
  'Gifts & Donations',
  'Subscriptions',
  'Miscellaneous'
];

// Options for the period:
const periodOptions = [
  'Daily',
  'Weekly',
  'Monthly',
  'Yearly',
  'One-time'
];

function NewBudgetRow({ onCreate }) {
  const [name, setName] = useState('');
  const [amount, setAmount] = useState('');
  const [category, setCategory] = useState(categoryOptions[0]); // default to first category
  const [period, setPeriod] = useState(periodOptions[0]); // default to first period

  const handleCreate = () => {
    const newBudget = {
      name,
      amount: parseFloat(amount) || 0,
      category,
      period,
    };
    onCreate(newBudget);
    setName('');
    setAmount('');
    setCategory(categoryOptions[0]); // reset to default category
    setPeriod(periodOptions[0]); // reset to default period
  };

  return (
    <tr className="bg-gray-50 border-b">
      <td className="px-2 py-1 text-xs">
        <input 
          className="w-full border rounded p-1 text-xs" 
          value={name} 
          onChange={(e) => setName(e.target.value)} 
          placeholder="Name"
        />
      </td>
      <td className="px-2 py-1 text-xs">
        <input 
          className="w-full border rounded p-1 text-xs" 
          value={amount} 
          onChange={(e) => setAmount(e.target.value)} 
          placeholder="Amount"
          type="number"
        />
      </td>
      <td className="px-2 py-1 text-xs">
        <select
          className="w-full border rounded p-1 text-xs"
          value={category}
          onChange={(e) => setCategory(e.target.value)}
        >
          {categoryOptions.map((option) => (
            <option key={option} value={option}>{option}</option>
          ))}
        </select>
      </td>
      <td className="px-2 py-1 text-xs">
        <select
          className="w-full border rounded p-1 text-xs"
          value={period}
          onChange={(e) => setPeriod(e.target.value)}
        >
          {periodOptions.map((option) => (
            <option key={option} value={option}>{option}</option>
          ))}
        </select>
      </td>
      <td className="px-2 py-1 text-xs">
        <button 
          className="bg-green-500 hover:bg-green-600 text-white rounded px-2 py-1 text-xs w-full"
          onClick={handleCreate}
        >
          Create
        </button>
      </td>
    </tr>
  );
}

function BudgetRow({ budget, onUpdate, onDelete }) {
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState(budget.name);
  const [editAmount, setEditAmount] = useState(budget.amount);
  const [editCategory, setEditCategory] = useState(budget.category);
  const [editPeriod, setEditPeriod] = useState(budget.period);

  const handleSave = () => {
    const updatedBudget = {
      ...budget,
      name: editName,
      amount: parseFloat(editAmount) || 0,
      category: editCategory,
      period: editPeriod,
    };
    onUpdate(updatedBudget);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setIsEditing(false);
    setEditName(budget.name);
    setEditAmount(budget.amount);
    setEditCategory(budget.category);
    setEditPeriod(budget.period);
  };

  const handleDelete = () => {
    onDelete(budget.id);
  };

  if (isEditing) {
    return (
      <tr className="bg-white border-b">
        <td className="px-2 py-1 text-xs">
          <input 
            className="w-full border rounded p-1 text-xs" 
            value={editName}
            onChange={(e) => setEditName(e.target.value)}
          />
        </td>
        <td className="px-2 py-1 text-xs">
          <input 
            className="w-full border rounded p-1 text-xs"
            type="number"
            value={editAmount}
            onChange={(e) => setEditAmount(e.target.value)}
          />
        </td>
        <td className="px-2 py-1 text-xs">
          <select
            className="w-full border rounded p-1 text-xs"
            value={editCategory}
            onChange={(e) => setEditCategory(e.target.value)}
          >
            {categoryOptions.map((option) => (
              <option key={option} value={option}>{option}</option>
            ))}
          </select>
        </td>
        <td className="px-2 py-1 text-xs">
          <select
            className="w-full border rounded p-1 text-xs"
            value={editPeriod}
            onChange={(e) => setEditPeriod(e.target.value)}
          >
            {periodOptions.map((option) => (
              <option key={option} value={option}>{option}</option>
            ))}
          </select>
        </td>
        <td className="px-2 py-1 text-xs">
          <div className="flex items-center justify-center space-x-1">
            <button 
              className="bg-blue-500 hover:bg-blue-600 text-white rounded px-2 py-1 text-xs"
              onClick={handleSave}
            >
              Save
            </button>
            <button 
              className="bg-gray-500 hover:bg-gray-600 text-white rounded px-2 py-1 text-xs"
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </td>
      </tr>
    );
  } else {
    return (
      <tr className="bg-white border-b">
        <td className="px-2 py-1 text-xs">{budget.name}</td>
        <td className="px-2 py-1 text-xs">{budget.amount}</td>
        <td className="px-2 py-1 text-xs">{budget.category}</td>
        <td className="px-2 py-1 text-xs">{budget.period}</td>
        <td className="px-2 py-1 text-xs">
          <div className="flex items-center justify-center space-x-1">
            <button 
              className="bg-blue-500 hover:bg-blue-600 text-white rounded px-2 py-1 text-xs"
              onClick={() => setIsEditing(true)}
            >
              Edit
            </button>
            <button 
              className="bg-red-500 hover:bg-red-600 text-white rounded px-2 py-1 text-xs"
              onClick={handleDelete}
            >
              Delete
            </button>
          </div>
        </td>
      </tr>
    );
  }
}

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
          setBudgets(data || []);
        } else {
          setError('Failed to fetch budgets');
        }
      } catch (err) {
        setError('Error fetching budgets');
      }
    };

    fetchBudgets();
  }, [token]);

  const handleCreateBudget = async (newBudget) => {
    try {
      const response = await fetch('http://localhost:8080/api/budgets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(newBudget),
      });

      if (response.ok) {
        const created = await response.json();
        setBudgets((prev) => [...prev, created]);
      } else {
        setError('Failed to create budget');
      }
    } catch (err) {
      setError('Error creating budget');
    }
  };

  const handleUpdateBudget = async (updatedBudget) => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/budgets/${updatedBudget.id}`, 
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify(updatedBudget),
        }
      );

      if (response.ok) {
        setBudgets((prev) =>
          prev.map((b) => (b.id === updatedBudget.id ? updatedBudget : b))
        );
      } else {
        setError('Failed to update budget');
      }
    } catch (err) {
      setError('Error updating budget');
    }
  };

  const handleDeleteBudget = async (budgetId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/budgets/${budgetId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        setBudgets((prev) => prev.filter((b) => b.id !== budgetId));
      } else {
        setError('Failed to delete budget');
      }
    } catch (err) {
      setError('Error deleting budget');
    }
  };

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold text-center mb-4">Your Budgets</h2>
      {error && <p className="text-red-500 text-center mb-2">{error}</p>}
      <div className="shadow-md sm:rounded-lg">
        <table className="min-w-full table-fixed text-sm text-center text-gray-500">
          <colgroup>
            <col className="w-1/5" />
            <col className="w-1/5" />
            <col className="w-1/5" />
            <col className="w-1/5" />
            <col className="w-1/5" />
          </colgroup>
          <thead className="text-xs text-white uppercase bg-blue-500">
            <tr>
              <th className="px-2 py-1">Name</th>
              <th className="px-2 py-1">Amount</th>
              <th className="px-2 py-1">Category</th>
              <th className="px-2 py-1">Period</th>
              <th className="px-2 py-1">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white">
            {budgets.map((budget) => (
              <BudgetRow
                key={budget.id}
                budget={budget}
                onUpdate={handleUpdateBudget}
                onDelete={handleDeleteBudget}
              />
            ))}
            <NewBudgetRow onCreate={handleCreateBudget} />
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default BudgetsPage;
