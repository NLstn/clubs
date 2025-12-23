import Foundation

class KeychainService {
    private let service = "com.clubs.ios"
    private let accessTokenKey = "accessToken"
    private let refreshTokenKey = "refreshToken"
    
    func storeTokens(access: String, refresh: String) throws {
        try storeToken(access, key: accessTokenKey)
        try storeToken(refresh, key: refreshTokenKey)
    }
    
    func getAccessToken() -> String? {
        retrieveToken(key: accessTokenKey)
    }
    
    func getRefreshToken() -> String? {
        retrieveToken(key: refreshTokenKey)
    }
    
    func deleteTokens() {
        deleteToken(key: accessTokenKey)
        deleteToken(key: refreshTokenKey)
    }
    
    // MARK: - Private Methods
    
    private func storeToken(_ token: String, key: String) throws {
        let data = token.data(using: .utf8)!
        
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecValueData as String: data
        ]
        
        SecItemDelete(query as CFDictionary)
        
        let status = SecItemAdd(query as CFDictionary, nil)
        guard status == errSecSuccess else {
            throw KeychainError.storeFailure(status)
        }
    }
    
    private func retrieveToken(key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true
        ]
        
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        
        guard status == errSecSuccess,
              let data = result as? Data,
              let token = String(data: data, encoding: .utf8) else {
            return nil
        }
        
        return token
    }
    
    private func deleteToken(key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key
        ]
        
        SecItemDelete(query as CFDictionary)
    }
}

enum KeychainError: LocalizedError {
    case storeFailure(OSStatus)
    
    var errorDescription: String? {
        switch self {
        case .storeFailure(let status):
            return "Keychain store failed with status: \(status)"
        }
    }
}
