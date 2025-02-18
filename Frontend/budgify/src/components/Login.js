import React, { useState } from 'react';
import './Login.css'; // Ensure the path is correct

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [passwordVisible, setPasswordVisible] = useState(false);

  const togglePasswordVisibility = () => {
    setPasswordVisible(!passwordVisible);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      const data = await response.json();
      if (response.ok) {
        onLogin(data.token, data.permissions);
      } else {
        setError(data.message || 'Login failed');
      }
    } catch (err) {
      setError('Login error');
    }
  };

  return (
    <div className="login-overlay">
      <form className="login-form" onSubmit={handleSubmit}>
        <div className="login-container">
          <header className="login-header">
            <h2 className="login-title">Log In</h2>
            <p className="login-subtitle">
              login here using your username and password
            </p>
          </header>
          {error && <p className="login-error">{error}</p>}
          <div className="login-fieldset">
            <span className="login-input-item">
              <i className="fa fa-user-circle"></i>
            </span>
            <input
              className="login-input"
              id="login-txt-input"
              type="text"
              placeholder="@UserName"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <br />
            <span className="login-input-item">
              <i className="fa fa-key"></i>
            </span>
            <input
              className="login-input"
              type={passwordVisible ? 'text' : 'password'}
              placeholder="Password"
              id="login-pwd"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <span className="login-eye-container">
              <i
                className="fa fa-eye login-eye"
                aria-hidden="true"
                onClick={togglePasswordVisibility}
              ></i>
            </span>
            <br />
            <button className="login-button login-submit" type="submit">
              Log In
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}

export default Login;
