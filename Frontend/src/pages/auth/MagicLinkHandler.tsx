import React, { useEffect, useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

const MagicLinkHandler: React.FC = () => {
  const [status, setStatus] = useState<'verifying' | 'success' | 'error'>('verifying');
  const [errorMessage, setErrorMessage] = useState('');
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const hasVerified = useRef(false);

  useEffect(() => {
    const verifyToken = async () => {
      // Only proceed if we haven't already verified
      if (hasVerified.current) return;
      hasVerified.current = true;

      const params = new URLSearchParams(location.search);
      const token = params.get('token');

      if (!token) {
        setStatus('error');
        setErrorMessage('Invalid link. No token found.');
        return;
      }

      try {
        const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/verifyMagicLink?token=${token}`, {
          method: 'GET',
        });

        if (response.ok) {
          const data = await response.json();
          login(data.access, data.refresh);
          setStatus('success');
          
          // Redirect after short delay
          setTimeout(() => {
            navigate('/');
          }, 2000);
        } else {
          const error = await response.text();
          setStatus('error');
          setErrorMessage(`Authentication failed: ${error}`);
        }
      } catch {
        setStatus('error');
        setErrorMessage('Network error. Please try again.');
      }
    };

    verifyToken();
  }, [location, login, navigate]);

  return (
    <div className="magic-link-container">
      {status === 'verifying' && (
        <div>
          <h2>Verifying your login...</h2>
          <p>Please wait while we authenticate you.</p>
        </div>
      )}

      {status === 'success' && (
        <div>
          <h2>Login Successful!</h2>
          <p>You are now logged in. Redirecting...</p>
        </div>
      )}

      {status === 'error' && (
        <div>
          <h2>Login Failed</h2>
          <p>{errorMessage}</p>
          <p>
            <a href="/login">Return to login page</a>
          </p>
        </div>
      )}
    </div>
  );
};

export default MagicLinkHandler;