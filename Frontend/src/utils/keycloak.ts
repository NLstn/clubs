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
