import SwiftUI

struct ProfileView: View {
    @EnvironmentObject var authManager: AuthenticationManager
    @State private var showLogoutAlert = false
    
    var body: some View {
        NavigationStack {
            VStack(spacing: 16) {
                // Profile Header
                VStack(spacing: 12) {
                    // Avatar
                    Circle()
                        .fill(Color.green.opacity(0.2))
                        .frame(width: 80, height: 80)
                        .overlay(
                            Image(systemName: "person.fill")
                                .font(.system(size: 40))
                                .foregroundColor(.green)
                        )
                    
                    if let user = authManager.currentUser {
                        Text(user.displayName)
                            .font(.headline)
                        Text(user.email)
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                }
                .frame(maxWidth: .infinity)
                .padding(24)
                .background(Color(.systemGray6))
                .cornerRadius(12)
                
                // User Info
                if let user = authManager.currentUser {
                    VStack(spacing: 12) {
                        InfoRow(label: "First Name", value: user.firstName)
                        Divider()
                        InfoRow(label: "Last Name", value: user.lastName)
                        Divider()
                        InfoRow(label: "Email", value: user.email)
                        Divider()
                        InfoRow(label: "Role", value: user.role)
                    }
                    .padding(16)
                    .background(Color(.systemGray6))
                    .cornerRadius(12)
                }
                
                Spacer()
                
                // Logout Button
                Button(role: .destructive, action: {
                    showLogoutAlert = true
                }) {
                    HStack {
                        Image(systemName: "arrowtriang.left.fill")
                        Text("Logout")
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 12)
                    .background(Color.red.opacity(0.1))
                    .foregroundColor(.red)
                    .cornerRadius(8)
                }
            }
            .padding(16)
            .navigationTitle("Profile")
            .navigationBarTitleDisplayMode(.inline)
            .alert("Logout", isPresented: $showLogoutAlert, actions: {
                Button("Cancel", role: .cancel) { }
                Button("Logout", role: .destructive) {
                    Task {
                        await authManager.logout()
                    }
                }
            }, message: {
                Text("Are you sure you want to logout?")
            })
        }
    }
}

struct InfoRow: View {
    let label: String
    let value: String
    
    var body: some View {
        HStack {
            Text(label)
                .font(.subheadline)
                .foregroundColor(.secondary)
            
            Spacer()
            
            Text(value)
                .font(.subheadline)
                .fontWeight(.medium)
        }
    }
}

#Preview {
    ProfileView()
        .environmentObject(AuthenticationManager())
        .preferredColorScheme(.dark)
}
