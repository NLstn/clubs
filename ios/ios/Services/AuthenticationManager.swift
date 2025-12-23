import Foundation
import Combine

class AuthenticationManager: NSObject, ObservableObject {
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var errorMessage: String?
    @Published var loginSuccessMessage: String?
    
    private let keychainService = KeychainService()
    private let apiService = APIService()
    private var tokenRefreshTask: Task<Void, Never>?
    
    override init() {
        super.init()
        restoreAuthenticationState()
    }
    
    // MARK: - Authentication Methods
    
    func requestMagicLink(email: String) async {
        do {
            try await apiService.requestMagicLink(email: email)
            await MainActor.run {
                self.errorMessage = nil
            }
        } catch {
            await MainActor.run {
                self.errorMessage = error.localizedDescription
            }
        }
    }
    
    func verifyMagicLink(token: String) async {
        do {
            let authResponse = try await apiService.verifyMagicLink(token: token)
            
            // Store tokens in keychain
            try keychainService.storeTokens(
                access: authResponse.access,
                refresh: authResponse.refresh
            )
            
            // Fetch current user
            let user = try await apiService.getCurrentUser(accessToken: authResponse.access, userID: authResponse.userID)
            
            await MainActor.run {
                self.currentUser = user
                self.isAuthenticated = true
                self.errorMessage = nil
                
                // Schedule token refresh
                self.scheduleTokenRefresh()
                self.markLoginSuccess()
            }
        } catch {
            await MainActor.run {
                self.errorMessage = error.localizedDescription
            }
        }
    }

    func verifyMagicCode(code: String) async {
        do {
            let authResponse = try await apiService.verifyMagicCode(code: code)

            // Store tokens in keychain
            try keychainService.storeTokens(
                access: authResponse.access,
                refresh: authResponse.refresh
            )

            // Fetch current user
            let user = try await apiService.getCurrentUser(accessToken: authResponse.access, userID: authResponse.userID)

            await MainActor.run {
                self.currentUser = user
                self.isAuthenticated = true
                self.errorMessage = nil
                self.scheduleTokenRefresh()
                self.markLoginSuccess()
            }
        } catch {
            await MainActor.run {
                self.errorMessage = error.localizedDescription
            }
        }
    }

    private func markLoginSuccess() {
        loginSuccessMessage = "Signed in successfully"
        // Auto-dismiss after 2 seconds
        Task { @MainActor in
            try? await Task.sleep(nanoseconds: 2_000_000_000)
            self.loginSuccessMessage = nil
        }
    }
    
    func logout() async {
        // Invalidate refresh token on server
        if let refreshToken = keychainService.getRefreshToken() {
            _ = try? await apiService.logout(refreshToken: refreshToken)
        }
        
        // Clear tokens from keychain
        keychainService.deleteTokens()
        
        // Reset state
        await MainActor.run {
            self.isAuthenticated = false
            self.currentUser = nil
            self.errorMessage = nil
            self.tokenRefreshTask?.cancel()
        }
    }
    
    // MARK: - Private Methods
    
    private func restoreAuthenticationState() {
        if let accessToken = keychainService.getAccessToken(),
           !isTokenExpired(accessToken) {
            isAuthenticated = true
            scheduleTokenRefresh()
        } else if let refreshToken = keychainService.getRefreshToken() {
            // Try to refresh the token
            Task {
                await refreshAccessToken()
            }
        }
    }
    
    func refreshAccessToken() async {
        guard let refreshToken = keychainService.getRefreshToken() else {
            await logout()
            return
        }
        
        do {
            let authResponse = try await apiService.refreshToken(refreshToken: refreshToken)
            
            // Store new tokens
            try keychainService.storeTokens(
                access: authResponse.access,
                refresh: authResponse.refresh
            )
            
            // Update API service with new token
            apiService.setAccessToken(authResponse.access)
            
            await MainActor.run {
                self.isAuthenticated = true
            }
            
            scheduleTokenRefresh()
        } catch {
            await logout()
        }
    }
    
    private func scheduleTokenRefresh() {
        tokenRefreshTask?.cancel()
        
        guard let accessToken = keychainService.getAccessToken() else { return }
        
        // Decode token to get expiration time
        if let expirationDate = extractExpirationDate(from: accessToken) {
            // Refresh token 5 minutes before expiration
            let refreshInterval = expirationDate.timeIntervalSinceNow - 300
            
            if refreshInterval > 0 {
                tokenRefreshTask = Task {
                    try? await Task.sleep(nanoseconds: UInt64(refreshInterval * 1_000_000_000))
                    await self.refreshAccessToken()
                }
            }
        }
    }
    
    private func isTokenExpired(_ token: String) -> Bool {
        guard let expirationDate = extractExpirationDate(from: token) else {
            return true
        }
        return Date() >= expirationDate
    }
    
    private func extractExpirationDate(from token: String) -> Date? {
        let parts = token.split(separator: ".")
        guard parts.count == 3 else { return nil }
        
        var payload = String(parts[1])
        // Add padding if needed
        let remainder = payload.count % 4
        if remainder > 0 {
            payload += String(repeating: "=", count: 4 - remainder)
        }
        
        guard let data = Data(base64Encoded: payload),
              let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              let exp = json["exp"] as? TimeInterval else {
            return nil
        }
        
        return Date(timeIntervalSince1970: exp)
    }
    
    func getAccessToken() -> String? {
        keychainService.getAccessToken()
    }
}

