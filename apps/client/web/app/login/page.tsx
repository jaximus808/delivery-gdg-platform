"use client"
import React, { useState } from 'react'
import './LoginSignup.css'
import Image from 'next/image'
import user_icon from '../../Assets/user.png'
import email_icon from '../../Assets/mail.png'
import password_icon from '../../Assets/Password.png'

const LoginSignup: React.FC = () => {
  const [action, setAction] = useState<"Login" | "Sign Up">("Login");
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
    setError('');
    setSuccess('');
  };

  const handleSubmit = async () => {
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const endpoint = action === "Sign Up" ? '/api/signup' : '/api/signin';
      const payload = action === "Sign Up" 
        ? { name: formData.name, email: formData.email, password: formData.password }
        : { email: formData.email, password: formData.password };

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Something went wrong');
      }

      setSuccess(action === "Sign Up" ? 'Account created successfully!' : 'Logged in successfully!');
      
      // Handle successful response (e.g., store token, redirect)
      if (data.token) {
        // Store token in localStorage or cookie
        // localStorage.setItem('authToken', data.token);
        // Redirect to dashboard or home page
        window.location.href = '/dashboard';
      }

      // Reset form
      setFormData({ name: '', email: '', password: '' });

    } catch (err: unknown) {
        if (err instanceof Error) {
             setError(err.message); 
        }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className='container'>
      <div className="header">
        <div className="text">{action}</div>
        <div className="underline"></div>
      </div>

      {error && <div className="error-message" style={{ color: 'red', textAlign: 'center', marginBottom: '10px' }}>{error}</div>}
      {success && <div className="success-message" style={{ color: 'green', textAlign: 'center', marginBottom: '10px' }}>{success}</div>}

      <div className="inputs">
        {action === "Login" ? <div></div> : 
          <div className="input">
            <div className="input-icon">
              <Image src={user_icon} alt="User icon" width={24} height={24} />
            </div>
            <input 
              type="text" 
              name="name"
              placeholder="Name" 
              value={formData.name}
              onChange={handleInputChange}
            />
          </div>
        }
        
        <div className="input">
          <div className="input-icon">
            <Image src={email_icon} alt="Email icon" width={24} height={24} />
          </div>
          <input 
            type="email" 
            name="email"
            placeholder="Email Id" 
            value={formData.email}
            onChange={handleInputChange}
          />
        </div>
        
        <div className="input">
          <div className="input-icon">
            <Image src={password_icon} alt="Password icon" width={24} height={24} />
          </div>
          <input 
            type="password" 
            name="password"
            placeholder="Password" 
            value={formData.password}
            onChange={handleInputChange}
          />
        </div>
      </div>

      {action === "Sign Up" ? <div></div> : 
        <div className="forgot-password">Lost Password? <span>Click Here!</span></div>
      }

      <div className="mode-toggle" style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '10px',
        marginTop: '20px',
        marginBottom: '10px',
        fontSize: '14px',
        color: '#666'
      }}>
        <span>{action === "Login" ? "Don't have an account?" : "Already have an account?"}</span>
        <span 
          onClick={() => { 
            setAction(action === "Login" ? "Sign Up" : "Login"); 
            setError(''); 
            setSuccess(''); 
          }}
          style={{
            color: '#4A90E2',
            cursor: 'pointer',
            fontWeight: 'bold',
            textDecoration: 'underline'
          }}
        >
          {action === "Login" ? "Sign Up" : "Login"}
        </span>
      </div>

      <button 
        className="submit-button"
        onClick={handleSubmit}
        disabled={loading}
        style={{
          width: '100%',
          padding: '15px',
          marginTop: '10px',
          backgroundColor: loading ? '#ccc' : '#4A90E2',
          color: 'white',
          border: 'none',
          borderRadius: '50px',
          cursor: loading ? 'not-allowed' : 'pointer',
          fontSize: '16px',
          fontWeight: 'bold'
        }}
      >
        {loading ? 'Processing...' : action === "Sign Up" ? 'Create Account' : 'Sign In'}
      </button>
    </div>
  )
}

export default LoginSignup