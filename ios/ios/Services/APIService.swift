import Foundation

class APIService {
    static let shared = APIService()
    
    private let baseURL = "http://localhost:8080/api/v1"
    private let session: URLSession
    private var accessToken: String?
    
    init() {
        let config = URLSessionConfiguration.default
        config.waitsForConnectivity = true
        config.timeoutIntervalForRequest = 30
        config.timeoutIntervalForResource = 300
        self.session = URLSession(configuration: config)
    }
    
    func setAccessToken(_ token: String) {
        self.accessToken = token
    }
    
    // MARK: - Authentication Endpoints
    
    func requestMagicLink(email: String) async throws {
        let endpoint = "/auth/requestMagicLink"
        let url = URL(string: baseURL + endpoint)!
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let payload = ["email": email]
        request.httpBody = try JSONSerialization.data(withJSONObject: payload)
        
        let (_, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw APIError(message: "Failed to request magic link", statusCode: nil)
        }
    }
    
    func verifyMagicLink(token: String) async throws -> AuthResponse {
        let endpoint = "/auth/verifyMagicLink"
        let urlString = baseURL + endpoint + "?token=\(token)"
        let url = URL(string: urlString)!
        
        let (data, response) = try await session.data(from: url)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError(message: "Invalid response", statusCode: nil)
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError(
                message: "Invalid or expired token",
                statusCode: httpResponse.statusCode
            )
        }
        
        let authResponse = try JSONDecoder.iso8601.decode(AuthResponse.self, from: data)
        return authResponse
    }

    func verifyMagicCode(code: String) async throws -> AuthResponse {
        let endpoint = "/auth/verifyMagicCode"
        let url = URL(string: baseURL + endpoint)!

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        let payload = ["code": code]
        request.httpBody = try JSONSerialization.data(withJSONObject: payload)

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError(message: "Invalid response", statusCode: nil)
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError(message: "Invalid or expired code", statusCode: httpResponse.statusCode)
        }
        let authResponse = try JSONDecoder.iso8601.decode(AuthResponse.self, from: data)
        return authResponse
    }
    
    func refreshToken(refreshToken: String) async throws -> AuthResponse {
        let endpoint = "/auth/refreshToken"
        let url = URL(string: baseURL + endpoint)!
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue(refreshToken, forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw APIError(message: "Failed to refresh token", statusCode: nil)
        }
        
        let authResponse = try JSONDecoder.iso8601.decode(AuthResponse.self, from: data)
        return authResponse
    }
    
    func logout(refreshToken: String) async throws {
        let endpoint = "/auth/logout"
        let url = URL(string: baseURL + endpoint)!
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue(refreshToken, forHTTPHeaderField: "Authorization")
        
        _ = try await session.data(for: request)
    }
    
    // MARK: - User Endpoints
    
    func getCurrentUser(accessToken: String, userID: String) async throws -> User {
        // Use OData v2 endpoint to get user by ID
        let url = URL(string: "http://localhost:8080/api/v2/Users('\(userID)')")!
        
        var request = URLRequest(url: url)
        request.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError(message: "Invalid response", statusCode: nil)
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError(
                message: "Failed to fetch user",
                statusCode: httpResponse.statusCode
            )
        }
        
        let user = try JSONDecoder.iso8601.decode(User.self, from: data)
        return user
    }
    
    // MARK: - Club Endpoints
    
    func getClubs() async throws -> [Club] {
        let endpoint = "/clubs"
        let clubs: [Club] = try await makeAuthorizedRequest(endpoint: endpoint)
        return clubs
    }
    
    func getClub(id: String) async throws -> Club {
        let endpoint = "/clubs/\(id)"
        let club: Club = try await makeAuthorizedRequest(endpoint: endpoint)
        return club
    }
    
    func getClubMembers(clubId: String) async throws -> [ClubMember] {
        let endpoint = "/clubs/\(clubId)/members"
        let members: [ClubMember] = try await makeAuthorizedRequest(endpoint: endpoint)
        return members
    }
    
    // MARK: - Event Endpoints
    
    func getEvents() async throws -> [Event] {
        let endpoint = "/events"
        let events: [Event] = try await makeAuthorizedRequest(endpoint: endpoint)
        return events
    }
    
    func getClubEvents(clubId: String) async throws -> [Event] {
        let endpoint = "/clubs/\(clubId)/events"
        let events: [Event] = try await makeAuthorizedRequest(endpoint: endpoint)
        return events
    }
    
    // MARK: - Fine Endpoints
    
    func getFines() async throws -> [Fine] {
        let endpoint = "/fines"
        let fines: [Fine] = try await makeAuthorizedRequest(endpoint: endpoint)
        return fines
    }
    
    func getClubFines(clubId: String) async throws -> [Fine] {
        let endpoint = "/clubs/\(clubId)/fines"
        let fines: [Fine] = try await makeAuthorizedRequest(endpoint: endpoint)
        return fines
    }
    
    // MARK: - Helper Methods
    
    private func makeAuthorizedRequest<T: Decodable>(
        accessToken: String? = nil,
        endpoint: String
    ) async throws -> T {
        let token = accessToken ?? self.accessToken
        guard let token = token else {
            throw APIError(message: "Not authenticated", statusCode: 401)
        }
        
        let url = URL(string: baseURL + endpoint)!
        var request = URLRequest(url: url)
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError(message: "Invalid response", statusCode: nil)
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError(
                message: "Request failed",
                statusCode: httpResponse.statusCode
            )
        }
        
        return try JSONDecoder.iso8601.decode(T.self, from: data)
    }
}

