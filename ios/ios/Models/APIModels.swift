import Foundation

// MARK: - Authentication Models

struct AuthResponse: Codable {
    let access: String
    let refresh: String
    let userID: String
    
    enum CodingKeys: String, CodingKey {
        case access
        case refresh
        case userID
    }
}

struct AuthToken {
    let access: String
    let refresh: String
    let expiresAt: Date
}

// MARK: - User Model

struct User: Codable, Identifiable {
    let id: String
    let email: String
    let firstName: String
    let lastName: String
    let role: String
    let profileImageUrl: String?
    
    enum CodingKeys: String, CodingKey {
        case id = "Id"
        case email = "Email"
        case firstName = "FirstName"
        case lastName = "LastName"
        case role = "Role"
        case profileImageUrl = "ProfileImageUrl"
    }
    
    var displayName: String {
        "\(firstName) \(lastName)"
    }
}

// MARK: - Club Models

struct Club: Codable, Identifiable {
    let id: String
    let name: String
    let description: String?
    let memberCount: Int
    let createdAt: Date?
    let imageUrl: String?
    
    enum CodingKeys: String, CodingKey {
        case id = "Id"
        case name = "Name"
        case description = "Description"
        case memberCount = "MemberCount"
        case createdAt = "CreatedAt"
        case imageUrl = "ImageUrl"
    }
}

struct ClubMember: Codable, Identifiable {
    let id: String
    let userId: String
    let userName: String
    let email: String
    let role: String
    let joinedAt: Date?
    
    enum CodingKeys: String, CodingKey {
        case id = "Id"
        case userId = "UserId"
        case userName = "UserName"
        case email = "Email"
        case role = "Role"
        case joinedAt = "JoinedAt"
    }
}

// MARK: - Event Models

struct Event: Codable, Identifiable {
    let id: String
    let clubId: String
    let title: String
    let description: String?
    let startTime: Date
    let endTime: Date?
    let location: String?
    let isRecurring: Bool
    let attendeeCount: Int?
    
    enum CodingKeys: String, CodingKey {
        case id = "Id"
        case clubId = "ClubId"
        case title = "Title"
        case description = "Description"
        case startTime = "StartTime"
        case endTime = "EndTime"
        case location = "Location"
        case isRecurring = "IsRecurring"
        case attendeeCount = "AttendeeCount"
    }
    
    var formattedDate: String {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        formatter.timeStyle = .short
        return formatter.string(from: startTime)
    }
}

// MARK: - Fine Models

struct Fine: Codable, Identifiable {
    let id: String
    let clubId: String
    let userId: String
    let amount: Double
    let reason: String
    let status: String
    let createdAt: Date?
    let dueDate: Date?
    
    enum CodingKeys: String, CodingKey {
        case id = "Id"
        case clubId = "ClubId"
        case userId = "UserId"
        case amount = "Amount"
        case reason = "Reason"
        case status = "Status"
        case createdAt = "CreatedAt"
        case dueDate = "DueDate"
    }
    
    var statusCategory: String {
        switch status.lowercased() {
        case "paid":
            return "paid"
        case "overdue":
            return "overdue"
        case "pending":
            return "pending"
        default:
            return "other"
        }
    }
}

#if canImport(SwiftUI)
import SwiftUI

extension Fine {
    var statusColor: Color {
        switch statusCategory {
        case "paid":
            return .green
        case "overdue":
            return .red
        case "pending":
            return .yellow
        default:
            return .gray
        }
    }
}
#endif

// MARK: - API Error

struct APIError: LocalizedError {
    let message: String
    let statusCode: Int?
    
    var errorDescription: String? {
        message
    }
}

// MARK: - Date Coding Strategy

extension JSONDecoder {
    static var iso8601: JSONDecoder {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        return decoder
    }
}

extension JSONEncoder {
    static var iso8601: JSONEncoder {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        return encoder
    }
}
