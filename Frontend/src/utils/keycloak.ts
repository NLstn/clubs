import { UserManager, User, UserManagerSettings, WebStorageStateStore } from 'oidc-client-ts';
import storage from './isomorphicStorage';

const keycloakConfig: UserManagerSettings = {
  authority: import.meta.env.VITE_KEYCLOAK_URL || 'https://auth.clubsstaging.dev/realms/clubs-dev',
  client_id: import.meta.env.VITE_KEYCLOAK_CLIENT_ID || 'clubs-frontend',
  response_type: 'code',
  scope: 'openid profile email',
  automaticSilentRenew: false,
  filterProtocolClaims: true,
  loadUserInfo: false,
  monitorSession: false,
};

if (typeof window !== 'undefined') {
  Object.assign(keycloakConfig, {
    redirect_uri: `${window.location.origin}/auth/callback`,
    post_logout_redirect_uri: `${window.location.origin}/login`,
    silent_redirect_uri: `${window.location.origin}/auth/silent-callback`,
    stateStore: new WebStorageStateStore({ store: window.localStorage }),
  });
}

class KeycloakService {
  private userManager: UserManager;

  constructor() {
    this.userManager = new UserManager(keycloakConfig);
    
    // Handle token expiration
    this.userManager.events.addUserSignedOut(() => {
      this.handleSignOut();
    });

    this.userManager.events.addSilentRenewError(() => {
      // Silent renew failed - user will need to re-authenticate
    });
  }

  async signinRedirect(forceLogin: boolean = false): Promise<void> {
    if (forceLogin) {
      await this.userManager.signinRedirect({ extraQueryParams: { prompt: 'login' } });
    } else {
      await this.userManager.signinRedirect();
    }
  }

  async signinRedirectCallback(): Promise<User> {
    return await this.userManager.signinRedirectCallback();
  }

  async signoutRedirect(): Promise<void> {
    await this.userManager.signoutRedirect();
  }

  async getUser(): Promise<User | null> {
    try {
      return await this.userManager.getUser();
    } catch {
      return null;
    }
  }

  async removeUser(): Promise<void> {
    await this.userManager.removeUser();
  }

  private handleSignOut(): void {
    // Clear local storage and redirect to login
    storage.removeItem('auth_token');
    storage.removeItem('refresh_token');
    storage.removeItem('keycloak_id_token');
    this.clearCallbackState();
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
  }

  async getAccessToken(): Promise<string | null> {
    const user = await this.getUser();
    return user?.access_token || null;
  }

  async isAuthenticated(): Promise<boolean> {
    const user = await this.getUser();
    return !!user && !user.expired;
  }

  // Clear callback processing state
  clearCallbackState(): void {
    if (typeof window !== 'undefined') {
      window.sessionStorage.removeItem('keycloak_callback_processed');
    }
    storage.removeItem('oidc.user');
  }

  // Clear all Keycloak-related data (for complete logout)
  clearAllKeycloakData(): void {
    this.clearCallbackState();
    // Clear any other OIDC client storage
    storage.removeItem('oidc.user:' + keycloakConfig.authority + ':' + keycloakConfig.client_id);
    // Clear the stored ID token
    storage.removeItem('keycloak_id_token');
    if (typeof window !== 'undefined') {
      Object.keys(window.localStorage).forEach(key => {
        if (key.startsWith('oidc.') || key.startsWith('keycloak')) {
          window.localStorage.removeItem(key);
        }
      });
    }
  }
}

export const keycloakService = new KeycloakService();
export default keycloakService;
