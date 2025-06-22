import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import CookieConsent from '../../components/CookieConsent';
import { useT } from '../../hooks/useTranslation';

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
          <div className="form-group">
            <label htmlFor="email">{t('auth.email')}</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isSubmitting}
            />
          </div>
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? t('auth.sending') : t('auth.sendMagicLink')}
          </button>
        </form>
      </div>
      <CookieConsent />
    </div>
  );
};

export default Login;