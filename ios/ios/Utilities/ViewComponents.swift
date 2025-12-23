import SwiftUI

// MARK: - Reusable List Components

struct ListSectionHeader: View {
    let title: String
    
    var body: some View {
        Text(title)
            .font(.headline)
            .foregroundColor(.primary)
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
    }
}

struct StatusBadge: View {
    let status: String
    let color: Color
    
    var body: some View {
        Text(status)
            .font(.caption2)
            .fontWeight(.semibold)
            .padding(.vertical, 4)
            .padding(.horizontal, 8)
            .background(color.opacity(0.2))
            .foregroundColor(color)
            .cornerRadius(4)
    }
}

// MARK: - Navigation Helpers

struct NavigationButton: View {
    let icon: String
    let label: String
    
    var body: some View {
        HStack {
            Image(systemName: icon)
                .font(.system(size: 16))
            Text(label)
            Spacer()
            Image(systemName: "chevron.right")
                .font(.system(size: 14))
                .foregroundColor(.secondary)
        }
        .foregroundColor(.primary)
    }
}

// MARK: - Input Fields

struct EmailTextField: View {
    @Binding var text: String
    let placeholder: String
    
    var body: some View {
        TextField(placeholder, text: $text)
            .textFieldStyle(.roundedBorder)
            .textContentType(.emailAddress)
            .keyboardType(.emailAddress)
            .autocapitalization(.none)
            .autocorrectionDisabled()
    }
}

// MARK: - Data Display

struct DataRow: View {
    let label: String
    let value: String
    let icon: String?
    
    var body: some View {
        HStack(spacing: 12) {
            if let icon = icon {
                Image(systemName: icon)
                    .foregroundColor(.secondary)
                    .frame(width: 24)
            }
            
            VStack(alignment: .leading, spacing: 2) {
                Text(label)
                    .font(.caption)
                    .foregroundColor(.secondary)
                Text(value)
                    .font(.body)
                    .fontWeight(.semibold)
            }
            
            Spacer()
        }
    }
}

// MARK: - Animations

struct PulseEffect: ViewModifier {
    @State private var isAnimating = false
    
    func body(content: Content) -> some View {
        content
            .opacity(isAnimating ? 1.0 : 0.6)
            .onAppear {
                withAnimation(.easeInOut(duration: 1.5).repeatForever(autoreverses: true)) {
                    isAnimating = true
                }
            }
    }
}

extension View {
    func pulse() -> some View {
        modifier(PulseEffect())
    }
}

// MARK: - Custom Modifiers

struct HorizontalPadding: ViewModifier {
    func body(content: Content) -> some View {
        content
            .padding(.horizontal, 16)
    }
}

struct VerticalPadding: ViewModifier {
    func body(content: Content) -> some View {
        content
            .padding(.vertical, 12)
    }
}

extension View {
    func horizontalPadding() -> some View {
        modifier(HorizontalPadding())
    }
    
    func verticalPadding() -> some View {
        modifier(VerticalPadding())
    }
}

// MARK: - Dividers

struct CustomDivider: View {
    var body: some View {
        Divider()
            .overlay(Color(.systemGray4))
    }
}

// MARK: - Progress Indicators

struct CircleProgress: View {
    let progress: Double
    let color: Color
    
    var body: some View {
        ZStack {
            Circle()
                .stroke(Color(.systemGray5), lineWidth: 4)
            
            Circle()
                .trim(from: 0, to: progress)
                .stroke(color, style: StrokeStyle(lineWidth: 4, lineCap: .round))
                .rotationEffect(.degrees(-90))
                .animation(.easeInOut, value: progress)
            
            Text(String(format: "%.0f%%", progress * 100))
                .font(.caption)
                .fontWeight(.semibold)
        }
    }
}
