import SwiftUI

struct FinesListView: View {
    @State private var fines: [Fine] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var filterStatus = "All"
    
    var filteredFines: [Fine] {
        if filterStatus == "All" {
            return fines
        }
        return fines.filter { $0.status.lowercased() == filterStatus.lowercased() }
    }
    
    var body: some View {
        NavigationStack {
            Group {
                if isLoading {
                    ProgressView()
                } else if filteredFines.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "checkmark.circle.fill")
                            .font(.system(size: 48))
                            .foregroundColor(.green)
                        Text("No Fines")
                            .font(.headline)
                        Text("You don't have any fines")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .frame(maxHeight: .infinity, alignment: .center)
                } else {
                    VStack {
                        Picker("Status", selection: $filterStatus) {
                            Text("All").tag("All")
                            Text("Pending").tag("Pending")
                            Text("Paid").tag("Paid")
                            Text("Overdue").tag("Overdue")
                        }
                        .pickerStyle(.segmented)
                        .padding(12)
                        
                        List(filteredFines) { fine in
                            FineListItem(fine: fine)
                        }
                    }
                }
            }
            .navigationTitle("Fines")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button(action: loadFines) {
                        Image(systemName: "arrow.clockwise")
                    }
                }
            }
            .alert("Error", isPresented: .constant(errorMessage != nil), actions: {
                Button("OK") { errorMessage = nil }
            }, message: {
                if let error = errorMessage {
                    Text(error)
                }
            })
            .onAppear {
                loadFines()
            }
        }
    }
    
    private func loadFines() {
        isLoading = true
        errorMessage = nil
        
        Task {
            do {
                let apiService = APIService.shared
                self.fines = try await apiService.getFines()
                isLoading = false
            } catch {
                await MainActor.run {
                    errorMessage = error.localizedDescription
                    isLoading = false
                }
            }
        }
    }
}

struct FineListItem: View {
    let fine: Fine
    
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text(fine.reason)
                        .font(.headline)
                    
                    if let dueDate = fine.dueDate {
                        Text("Due: \(formattedDate(dueDate))")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
                
                Spacer()
                
                VStack(alignment: .trailing, spacing: 4) {
                    Text(String(format: "%.2f€", fine.amount))
                        .font(.headline)
                        .foregroundColor(.red)
                    
                    Text(fine.status)
                        .font(.caption)
                        .padding(.vertical, 2)
                        .padding(.horizontal, 6)
                        .background(fine.statusColor.opacity(0.2))
                        .foregroundColor(fine.statusColor)
                        .cornerRadius(4)
                }
            }
        }
        .padding(.vertical, 8)
    }
    
    private func formattedDate(_ date: Date) -> String {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        return formatter.string(from: date)
    }
}

#Preview {
    FinesListView()
        .preferredColorScheme(.dark)
}
