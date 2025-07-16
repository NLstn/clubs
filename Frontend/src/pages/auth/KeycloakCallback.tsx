import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useT } from '../../hooks/useTranslation';

// StrictMode is now disabled, so no need for global duplicate prevention
// const processedCodes = new Set<string>();

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

        // StrictMode is now disabled, so no need for duplicate prevention logic
        console.log('Processing callback for code:', code.substring(0, 8) + '...');

        console.log('Processing Keycloak callback with code');

        if (!isMounted) return;

        // Simplified approach: Send the authorization code directly to our backend
        // This avoids the complexity of managing OIDC client state and potential race conditions
        if (!state) {
          throw new Error('No state parameter - cannot validate callback');
        }

        // Create a timeout for the entire fetch operation
        const controller = new AbortController();
        const timeoutId = setTimeout(() => {
          console.log('Fetch timeout reached, aborting request');
          controller.abort();
        }, 15000); // 15 second timeout for the entire request

        const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/keycloak/callback`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ code, state }),
          signal: controller.signal,
        });

        clearTimeout(timeoutId);
        console.log('Backend response status:', response.status);

        if (!isMounted) return;

        if (!response.ok) {
          const errorText = await response.text();
          console.log('Backend error response:', errorText);
          throw new Error(`Backend authentication failed: ${errorText}`);
        }

        console.log('About to parse JSON response');
        
        // Implement a robust response reading mechanism
        let responseBody: string;
        
        try {
          console.log('Attempting to read response body...');
          
          // Try using the Response.body stream with a manual timeout
          if (response.body) {
            console.log('Using ReadableStream approach...');
            const reader = response.body.getReader();
            const chunks: Uint8Array[] = [];
            let totalLength = 0;
            
            // Set up a timeout for the stream reading
            const readTimeout = setTimeout(() => {
              console.error('Stream reading timeout - aborting');
              reader.cancel();
            }, 10000);
            
            try {
              while (true) {
                const { done, value } = await reader.read();
                if (done) break;
                
                chunks.push(value);
                totalLength += value.length;
                console.log('Read chunk, total bytes so far:', totalLength);
                
                // Safety check - prevent reading extremely large responses
                if (totalLength > 1024 * 1024) { // 1MB limit
                  throw new Error('Response too large');
                }
              }
              
              clearTimeout(readTimeout);
              
              // Combine all chunks into a single string
              const combined = new Uint8Array(totalLength);
              let offset = 0;
              for (const chunk of chunks) {
                combined.set(chunk, offset);
                offset += chunk.length;
              }
              
              responseBody = new TextDecoder().decode(combined);
              console.log('Stream reading completed, total length:', responseBody.length);
              
            } catch (streamError) {
              clearTimeout(readTimeout);
              console.error('Stream reading failed:', streamError);
              throw new Error('Failed to read response stream');
            }
          } else {
            // Fallback to text() method with manual timeout
            console.log('Falling back to response.text() method...');
            const textPromise = response.text();
            const timeoutPromise = new Promise<never>((_, reject) => {
              setTimeout(() => reject(new Error('text() method timeout')), 10000);
            });
            
            responseBody = await Promise.race([textPromise, timeoutPromise]);
          }
          
          console.log('Response body read successfully, length:', responseBody.length);
          console.log('Response preview:', responseBody.substring(0, 200));
          
        } catch (readError) {
          console.error('Failed to read response body:', readError);
          throw new Error(`Response reading failed: ${readError instanceof Error ? readError.message : 'Unknown error'}`);
        }
        
        // Now parse the JSON
        let data: { access: string; refresh: string; keycloakTokens?: { idToken: string } };
        try {
          console.log('Parsing JSON from response body...');
          data = JSON.parse(responseBody);
          console.log('JSON parse successful');
        } catch (jsonError) {
          console.error('Failed to parse JSON:', jsonError);
          console.error('Raw response body:', responseBody);
          throw new Error('Invalid JSON response from server');
        }
        
        console.log('Access token:', data.access ? 'Present' : 'Missing');
        console.log('Refresh token:', data.refresh ? 'Present' : 'Missing');
        
        if (!data.access || !data.refresh) {
          throw new Error('Missing access or refresh token in response');
        }
        
        // Store our application tokens
        console.log('Calling login function...');
        login(data.access, data.refresh);
        console.log('Login function completed successfully');

        // Store Keycloak tokens if available
        if (data.keycloakTokens && data.keycloakTokens.idToken) {
          localStorage.setItem('keycloak_id_token', data.keycloakTokens.idToken);
          console.log('Keycloak ID token stored');
        }

        // Clear the processing flag and redirect - simplified since no session storage
        const redirectPath = sessionStorage.getItem('auth_redirect_after_login') || '/';
        sessionStorage.removeItem('auth_redirect_after_login');
        
        console.log('Authentication completed successfully!');
        console.log('About to navigate to:', redirectPath);
        console.log('isMounted:', isMounted);
        
        if (isMounted) {
          console.log('Navigating now...');
          navigate(redirectPath, { replace: true });
          console.log('Navigate called');
        }
        
      } catch (error) {
        console.error('Keycloak callback error:', error);
        
        // No need to clear processing flags since StrictMode is disabled
        
        if (isMounted) {
          console.log('Setting error state due to:', error);
          
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
