import SwiftUI

// MARK: - Colors

extension Color {
    static let clubsGreen = Color(red: 0.298, green: 0.686, blue: 0.314)
    static let clubsBlue = Color(red: 0.4, green: 0.424, blue: 1.0)
    static let clubsRed = Color(red: 0.956, green: 0.263, blue: 0.212)
    
    static let background = Color(red: 0.141, green: 0.141, blue: 0.141)
    static let backgroundLight = Color(red: 0.2, green: 0.2, blue: 0.2)
    static let textPrimary = Color(white: 0.87)
    static let textSecondary = Color(white: 0.533)
    static let borderColor = Color(white: 0.867)
}

// MARK: - Buttons

struct PrimaryButton: View {
    let title: String
    let action: () -> Void
    var isLoading = false
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 8) {
                if isLoading {
                    ProgressView()
                        .tint(.white)
                }
                Text(title)
                    .fontWeight(.semibold)
            }
            .frame(maxWidth: .infinity)
            .padding(.vertical, 12)
            .background(Color.clubsGreen)
            .foregroundColor(.white)
            .cornerRadius(8)
        }
        .disabled(isLoading)
    }
}

struct SecondaryButton: View {
    let title: String
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            Text(title)
                .fontWeight(.semibold)
                .frame(maxWidth: .infinity)
                .padding(.vertical, 12)
                .background(Color(.systemGray5))
                .foregroundColor(.primary)
                .cornerRadius(8)
        }
    }
}

// MARK: - Cards

struct CardContainer<Content: View>: View {
    let content: Content
    
    init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }
    
    var body: some View {
        VStack {
            content
        }
        .padding(16)
        .background(Color(.systemGray6))
        .cornerRadius(12)
    }
}

// MARK: - Empty State

struct EmptyState: View {
    let icon: String
    let title: String
    let description: String
    
    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: icon)
                .font(.system(size: 48))
                .foregroundColor(.gray)
            Text(title)
                .font(.headline)
            Text(description)
                .font(.subheadline)
                .foregroundColor(.secondary)
        }
        .frame(maxHeight: .infinity, alignment: .center)
        .multilineTextAlignment(.center)
    }
}

// MARK: - Error Banner

struct ErrorBanner: View {
    let message: String
    
    var body: some View {
        HStack(spacing: 12) {
            Image(systemName: "exclamationmark.circle.fill")
                .foregroundColor(.red)
            Text(message)
                .font(.subheadline)
            Spacer()
        }
        .padding(12)
        .background(Color.red.opacity(0.1))
        .cornerRadius(8)
    }
}

// MARK: - Loading Skeleton

struct LoadingSkeleton: View {
    var body: some View {
        VStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color(.systemGray5))
                .frame(height: 20)
            
            RoundedRectangle(cornerRadius: 8)
                .fill(Color(.systemGray5))
                .frame(height: 16)
            
            RoundedRectangle(cornerRadius: 8)
                .fill(Color(.systemGray5))
                .frame(height: 16)
                .frame(maxWidth: 200, alignment: .leading)
        }
        .padding(12)
        .background(Color(.systemGray6))
        .cornerRadius(8)
        .shimmer()
    }
}

// MARK: - Shimmer Effect

extension View {
    func shimmer() -> some View {
        modifier(ShimmerModifier())
    }
}

struct ShimmerModifier: ViewModifier {
    @State private var isShimmering = false
    
    func body(content: Content) -> some View {
        content
            .opacity(isShimmering ? 0.6 : 1.0)
            .onAppear {
                withAnimation(.easeInOut(duration: 1).repeatForever(autoreverses: true)) {
                    isShimmering = true
                }
            }
    }
}
