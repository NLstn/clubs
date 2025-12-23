import SwiftUI

struct LoginView: View {
    @EnvironmentObject var authManager: AuthenticationManager
    @State private var email = ""
    @State private var isLoading = false
    @State private var showMagicLinkInput = false
    @State private var code = ""
    @State private var isVerifyingCode = false
    
    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {
                // Header
                VStack(spacing: 8) {
                    Image(systemName: "building.2.fill")
                        .font(.system(size: 48))
                        .foregroundColor(.green)
                    
                    Text("Clubs Manager")
                        .font(.system(size: 28, weight: .bold))
                    
                    Text("Manage your club activities")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                .padding(.top, 60)
                
                Spacer()
                
                // Login Form
                VStack(spacing: 16) {
                    if !showMagicLinkInput {
                        VStack(spacing: 12) {
                            Text("Sign in to get started")
                                .font(.headline)
                                .frame(maxWidth: .infinity, alignment: .leading)
                            
                            TextField("Enter your email", text: $email)
                                .textFieldStyle(.roundedBorder)
                                .textContentType(.emailAddress)
                                .keyboardType(.emailAddress)
                                .autocapitalization(.none)
                            
                            Button(action: requestMagicLink) {
                                HStack {
                                    if isLoading {
                                        ProgressView()
                                            .tint(.white)
                                    } else {
                                        Text("Send Magic Link")
                                    }
                                }
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 12)
                                .background(Color.green)
                                .foregroundColor(.white)
                                .cornerRadius(8)
                            }
                            .disabled(email.isEmpty || isLoading)
                        }
                    } else {
                        VStack(spacing: 12) {
                            Text("Check your email")
                                .font(.headline)
                                .frame(maxWidth: .infinity, alignment: .leading)
                            
                            HStack {
                                Image(systemName: "checkmark.circle.fill")
                                    .foregroundColor(.green)
                                
                                VStack(alignment: .leading, spacing: 4) {
                                    Text("Link sent to \(email)")
                                        .font(.caption)
                                    Text("Click the link in your email to sign in")
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                            }
                            .padding(12)
                            .background(Color(.systemGray6))
                            .cornerRadius(8)
                            
                            Button(action: { showMagicLinkInput = false; email = "" }) {
                                Text("Use Different Email")
                                    .frame(maxWidth: .infinity)
                                    .padding(.vertical, 12)
                                    .background(Color(.systemGray5))
                                    .foregroundColor(.primary)
                                    .cornerRadius(8)
                            }
                            
                            Text("Open email app or enter your 6-digit code")
                                .font(.caption)
                                .foregroundColor(.secondary)
                                .padding(.top, 8)

                            VStack(spacing: 8) {
                                TextField("Enter 6-digit code", text: $code)
                                    .textFieldStyle(.roundedBorder)
                                    .keyboardType(.numberPad)
                                    .onChange(of: code) { newValue in
                                        // Keep only digits and max length 6
                                        let filtered = newValue.filter { $0.isNumber }
                                        if filtered.count > 6 {
                                            code = String(filtered.prefix(6))
                                        } else {
                                            code = filtered
                                        }
                                    }
                                Button(action: submitCode) {
                                    HStack {
                                        if isVerifyingCode {
                                            ProgressView().tint(.white)
                                        } else {
                                            Text("Verify Code")
                                        }
                                    }
                                    .frame(maxWidth: .infinity)
                                    .padding(.vertical, 12)
                                    .background(code.count == 6 ? Color.green : Color.gray)
                                    .foregroundColor(.white)
                                    .cornerRadius(8)
                                }
                                .disabled(code.count != 6 || isVerifyingCode)
                            }
                        }
                    }
                    
                    if let error = authManager.errorMessage {
                        HStack {
                            Image(systemName: "exclamationmark.circle.fill")
                                .foregroundColor(.red)
                            Text(error)
                                .font(.caption)
                        }
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .padding(12)
                        .background(Color.red.opacity(0.1))
                        .cornerRadius(8)
                    }
                }
                .padding(24)
                .background(Color(.systemGray6))
                .cornerRadius(12)
                
                Spacer()
                
                // Footer
                Text("We'll send you a secure link to sign in")
                    .font(.caption)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
            }
            .padding(20)
            .navigationBarTitleDisplayMode(.inline)
        }
    }
    
    private func requestMagicLink() {
        guard !email.isEmpty else { return }
        
        isLoading = true
        Task {
            await authManager.requestMagicLink(email: email)
            await MainActor.run {
                isLoading = false
                if authManager.errorMessage == nil {
                    showMagicLinkInput = true
                }
            }
        }
    }

    private func submitCode() {
        guard code.count == 6 else { return }
        isVerifyingCode = true
        Task {
            await authManager.verifyMagicCode(code: code)
            await MainActor.run {
                isVerifyingCode = false
            }
        }
    }
}

#Preview {
    LoginView()
        .environmentObject(AuthenticationManager())
        .preferredColorScheme(.dark)
}
