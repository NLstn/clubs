import { UserManager, User, UserManagerSettings, WebStorageStateStore } from 'oidc-client-ts';

const keycloakConfig: UserManagerSettings = {
  authority: import.meta.env.VITE_KEYCLOAK_URL || 'https://auth.clubsstaging.dev/realms/clubs-dev',
  client_id: import.meta.env.VITE_KEYCLOAK_CLIENT_ID || 'clubs-frontend',
  redirect_uri: `${window.location.origin}/auth/callback`,
  post_logout_redirect_uri: `${window.location.origin}/login`,
  response_type: 'code',
  scope: 'openid profile email',
  automaticSilentRenew: false,
  silent_redirect_uri: `${window.location.origin}/auth/silent-callback`,
  filterProtocolClaims: true,
  loadUserInfo: false,
  monitorSession: false,
  // Add state validation
  stateStore: new WebStorageStateStore({ store: window.localStorage }),
};

class KeycloakService {
  private userManager: UserManager;

  constructor() {
    this.userManager = new UserManager(keycloakConfig);
    
    // Handle token expiration
    this.userManager.events.addUserSignedOut(() => {
      this.handleSignOut();
    });

    this.userManager.events.addSilentRenewError((error) => {
      console.error('Silent renew error:', error);
    });
  }

  async signinRedirect(forceLogin: boolean = false): Promise<void> {
    try {
      console.log('Keycloak signinRedirect called with forceLogin:', forceLogin);
      if (forceLogin) {
        console.log('Using prompt=login to force fresh authentication');
        await this.userManager.signinRedirect({ extraQueryParams: { prompt: 'login' } });
      } else {
        console.log('Using normal signin redirect (may use existing session)');
        await this.userManager.signinRedirect();
      }
    } catch (error) {
      console.error('Signin redirect error:', error);
      throw error;
    }
  }

  async signinRedirectCallback(): Promise<User> {
    try {
      const user = await this.userManager.signinRedirectCallback();
      return user;
    } catch (error) {
      console.error('Signin redirect callback error:', error);
      throw error;
    }
  }

  async signoutRedirect(): Promise<void> {
    try {
      await this.userManager.signoutRedirect();
    } catch (error) {
      console.error('Signout redirect error:', error);
      throw error;
    }
  }

  async getUser(): Promise<User | null> {
    try {
      return await this.userManager.getUser();
    } catch (error) {
      console.error('Get user error:', error);
      return null;
    }
  }

  async removeUser(): Promise<void> {
    try {
      await this.userManager.removeUser();
    } catch (error) {
      console.error('Remove user error:', error);
    }
  }

  private handleSignOut(): void {
    // Clear local storage and redirect to login
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('keycloak_id_token');
    this.clearCallbackState();
    window.location.href = '/login';
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
    sessionStorage.removeItem('keycloak_callback_processed');
    localStorage.removeItem('oidc.user');
  }

  // Clear all Keycloak-related data (for complete logout)
  clearAllKeycloakData(): void {
    this.clearCallbackState();
    // Clear any other OIDC client storage
    localStorage.removeItem('oidc.user:' + keycloakConfig.authority + ':' + keycloakConfig.client_id);
    // Clear the stored ID token
    localStorage.removeItem('keycloak_id_token');
    // Clear any state store data
    Object.keys(localStorage).forEach(key => {
      if (key.startsWith('oidc.') || key.startsWith('keycloak')) {
        localStorage.removeItem(key);
      }
    });
  }
}

export const keycloakService = new KeycloakService();
export default keycloakService;
