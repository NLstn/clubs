import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import CookieConsent from '../../components/CookieConsent';
import { useT } from '../../hooks/useTranslation';
import { Input } from '@/components/ui';

const Login: React.FC = () => {
  const { t } = useT();
  const [email, setEmail] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState('');
  const location = useLocation();

  useEffect(() => {
    // Check if this is a redirect from a protected page
    const params = new URLSearchParams(location.search);
    const redirectPath = params.get('redirect');
    if (redirectPath) {
      setMessage(t('auth.redirectMessage'));
    }
  }, [location, t]);

  const handleKeycloakLogin = async () => {
    try {
      // Store redirect path for after login
      const params = new URLSearchParams(location.search);
      const redirectPath = params.get('redirect') || '/';
      sessionStorage.setItem('auth_redirect_after_login', redirectPath);
      
      // Get the Keycloak auth URL from our backend
      const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/keycloak/login`);
      if (!response.ok) {
        throw new Error('Failed to get auth URL');
      }
      
      const data = await response.json();
      
      // Redirect to Keycloak
      window.location.href = data.authURL;
    } catch {
      setMessage(t('auth.keycloakError'));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setMessage('');

    try {
      const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/requestMagicLink`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      if (response.ok) {
        setMessage(t('auth.checkEmail'));
        setEmail('');
      } else {
        const error = await response.text();
        setMessage(`${t('common.error')}: ${error}`);
      }
    } catch {
      setMessage(t('auth.networkError'));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>{t('auth.login')}</h1>
        <p>{t('auth.loginInstruction')}</p>

        {message && (
          <div className={`message ${message.includes('Error') || message.includes('error') ? 'error' : 'success'}`}>
            {message}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <Input
            label={t('auth.email')}
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={isSubmitting}
          />
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? t('auth.sending') : t('auth.sendMagicLink')}
          </button>
        </form>

        <div className="divider">
          <span>{t('auth.or')}</span>
        </div>

        <button 
          type="button" 
          className="keycloak-login-btn"
          onClick={handleKeycloakLogin}
          disabled={isSubmitting}
        >
          {t('auth.loginWithKeycloak')}
        </button>
      </div>
      <CookieConsent />
    </div>
  );
};

export default Login;