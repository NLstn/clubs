import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useT } from '../../hooks/useTranslation';

const KeycloakCallback: React.FC = () => {
  const { t } = useT();
  const navigate = useNavigate();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    
    const handleCallback = async () => {
      try {
        // Clear any existing OIDC client state that might interfere
        localStorage.removeItem('oidc.user:https://auth.clubsstaging.dev/realms/clubs-dev:clubs-frontend');
        Object.keys(localStorage).forEach(key => {
          if (key.startsWith('oidc.')) {
            localStorage.removeItem(key);
          }
        });

        // Clean up old callback processing flags
        Object.keys(sessionStorage).forEach(key => {
          if (key.startsWith('keycloak_callback_')) {
            sessionStorage.removeItem(key);
          }
        });

        // Check if we've already handled this callback to prevent double processing
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get('code');
        const state = urlParams.get('state');
        
        if (!code) {
          throw new Error('No authorization code received from Keycloak');
        }

        if (!state) {
          throw new Error('No state parameter - cannot validate callback');
        }

        if (!isMounted) return;

        // Retrieve the PKCE code verifier
        const codeVerifier = sessionStorage.getItem('keycloak_code_verifier');
        if (!codeVerifier) {
          throw new Error('Missing PKCE code verifier - authentication flow may have been interrupted');
        }

        // Create a timeout for the entire fetch operation
        const controller = new AbortController();
        const timeoutId = setTimeout(() => {
          controller.abort();
        }, 15000); // 15 second timeout for the entire request

        const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/keycloak/callback`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ code, state, codeVerifier }),
          signal: controller.signal,
        });

        clearTimeout(timeoutId);

        if (!isMounted) return;

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Backend authentication failed: ${errorText}`);
        }

        const data: { access: string; refresh: string; keycloakTokens?: { idToken: string } } = await response.json();
        
        if (!data.access || !data.refresh) {
          throw new Error('Missing access or refresh token in response');
        }
        
        // Store our application tokens
        login(data.access, data.refresh);

        // Store Keycloak tokens if available
        if (data.keycloakTokens && data.keycloakTokens.idToken) {
          localStorage.setItem('keycloak_id_token', data.keycloakTokens.idToken);
        }

        // Clear the processing flag and PKCE verifier, then redirect
        const redirectPath = sessionStorage.getItem('auth_redirect_after_login') || '/';
        sessionStorage.removeItem('auth_redirect_after_login');
        sessionStorage.removeItem('keycloak_code_verifier');
        
        if (isMounted) {
          navigate(redirectPath, { replace: true });
        }
        
      } catch (error) {
        if (isMounted) {
          // Handle specific error types
          let errorMessage = 'Authentication failed';
          if (error instanceof Error) {
            if (error.name === 'AbortError') {
              errorMessage = 'Authentication request timed out. Please try again.';
            } else if (error.message.includes('timeout')) {
              errorMessage = 'Authentication timed out. Please try again.';
            } else {
              errorMessage = error.message;
            }
          }
          
          setError(errorMessage);
          
          // Redirect to login page after a delay
          setTimeout(() => {
            if (isMounted) {
              navigate('/login', { replace: true });
            }
          }, 3000);
        }
      }
    };

    handleCallback();
    
    return () => {
      isMounted = false;
    };
  }, [navigate, login]);

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full bg-white rounded-lg shadow-md p-6">
          <div className="text-center">
            <div className="text-red-500 text-xl mb-4">⚠️</div>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              {t('auth.authenticationFailed')}
            </h2>
            <p className="text-gray-600 mb-4">{error}</p>
            <p className="text-sm text-gray-500">
              {t('auth.redirectingToLogin')}
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white rounded-lg shadow-md p-6">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            {t('auth.signingIn')}
          </h2>
          <p className="text-gray-600">
            {t('auth.pleaseWait')}
          </p>
        </div>
      </div>
    </div>
  );
};

export default KeycloakCallback;
