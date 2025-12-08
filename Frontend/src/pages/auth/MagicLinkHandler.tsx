import React, { useEffect, useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import CookieConsent from '../../components/CookieConsent';
import './MagicLinkHandler.css';

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
          
          // Redirect based on profile completion status
          setTimeout(() => {
            if (!data.profileComplete) {
              navigate('/signup');
            } else {
              const redirectPath = localStorage.getItem('loginRedirect');
              if (redirectPath) {
                localStorage.removeItem('loginRedirect');
                navigate(redirectPath);
              } else {
                navigate('/');
              }
            }
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
      <div className="magic-link-box">
        {status === 'verifying' && (
          <div>
            <h1>Verifying your login...</h1>
            <p>Please wait while we authenticate you.</p>
          </div>
        )}

        {status === 'success' && (
          <div>
            <h1>Login Successful!</h1>
            <p>You are now logged in. Redirecting...</p>
            <div className="message success">
              Authentication successful! Redirecting to dashboard...
            </div>
          </div>
        )}

        {status === 'error' && (
          <div>
            <h1>Login Failed</h1>
            <div className="message error">
              {errorMessage}
            </div>
            <div className="magic-link-actions">
              <a href="/login" className="return-link">Return to login page</a>
            </div>
          </div>
        )}
      </div>
      <CookieConsent />
    </div>
  );
};

export default MagicLinkHandler;