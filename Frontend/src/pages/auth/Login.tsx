import React, { useState, useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import CookieConsent from '../../components/CookieConsent';
import { useT } from '../../hooks/useTranslation';
import { Input, Button, Divider, ButtonState } from '@/components/ui';
import './Login.css';

const Login: React.FC = () => {
  const { t } = useT();
  const [email, setEmail] = useState('');
  const [magicLinkButtonState, setMagicLinkButtonState] = useState<ButtonState>('idle');
  const [keycloakButtonState, setKeycloakButtonState] = useState<ButtonState>('idle');
  const [message, setMessage] = useState('');
  const location = useLocation();
  const timeoutRefs = useRef<number[]>([]);

  useEffect(() => {
    // Check if this is a redirect from a protected page
    const params = new URLSearchParams(location.search);
    const redirectPath = params.get('redirect');
    if (redirectPath) {
      setMessage(t('auth.redirectMessage'));
    }
  }, [location, t]);

  useEffect(() => {
    // Cleanup timeouts on unmount
    const timeouts = timeoutRefs.current;
    return () => {
      timeouts.forEach(clearTimeout);
    };
  }, []);

  const handleKeycloakLogin = async () => {
    setKeycloakButtonState('loading');
    setMessage('');
    
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
      
      // Store PKCE code verifier for the callback
      if (data.codeVerifier) {
        sessionStorage.setItem('keycloak_code_verifier', data.codeVerifier);
      }
      
      // Redirect to Keycloak (don't set success state as we're redirecting)
      window.location.href = data.authURL;
    } catch {
      setKeycloakButtonState('error');
      setMessage(t('auth.keycloakError'));
      const timeoutId = window.setTimeout(() => setKeycloakButtonState('idle'), 3000);
      timeoutRefs.current.push(timeoutId);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMagicLinkButtonState('loading');
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
        setMagicLinkButtonState('success');
        setMessage(t('auth.checkEmail'));
        setEmail('');
        const timeoutId = window.setTimeout(() => setMagicLinkButtonState('idle'), 3000);
        timeoutRefs.current.push(timeoutId);
      } else {
        const error = await response.text();
        setMagicLinkButtonState('error');
        setMessage(`${t('common.error')}: ${error}`);
        const timeoutId = window.setTimeout(() => setMagicLinkButtonState('idle'), 3000);
        timeoutRefs.current.push(timeoutId);
      }
    } catch {
      setMagicLinkButtonState('error');
      setMessage(t('auth.networkError'));
      const timeoutId = window.setTimeout(() => setMagicLinkButtonState('idle'), 3000);
      timeoutRefs.current.push(timeoutId);
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
            disabled={magicLinkButtonState === 'loading'}
          />
          <Button 
            type="submit" 
            variant="primary" 
            fullWidth
            state={magicLinkButtonState}
            successMessage={t('auth.checkEmail')}
            errorMessage={t('auth.networkError')}
          >
            {t('auth.sendMagicLink')}
          </Button>
        </form>

        <Divider text={t('auth.or')} />

        <Button 
          type="button" 
          onClick={handleKeycloakLogin}
          variant="secondary"
          fullWidth
          state={keycloakButtonState}
          successMessage="Redirecting..."
          errorMessage={t('auth.keycloakError')}
        >
          {t('auth.loginWithKeycloak')}
        </Button>
      </div>
      <CookieConsent />
    </div>
  );
};

export default Login;