import SwiftUI

struct EventsListView: View {
    @State private var events: [Event] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    
    var body: some View {
        NavigationStack {
            Group {
                if isLoading {
                    ProgressView()
                } else if events.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "calendar")
                            .font(.system(size: 48))
                            .foregroundColor(.gray)
                        Text("No Events")
                            .font(.headline)
                        Text("No events scheduled yet")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .frame(maxHeight: .infinity, alignment: .center)
                } else {
                    List(events) { event in
                        VStack(alignment: .leading, spacing: 8) {
                            HStack {
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(event.title)
                                        .font(.headline)
                                    
                                    HStack(spacing: 12) {
                                        Label(event.formattedDate, systemImage: "calendar")
                                            .font(.caption)
                                            .foregroundColor(.secondary)
                                        
                                        if event.isRecurring {
                                            Label("Recurring", systemImage: "repeat")
                                                .font(.caption)
                                                .foregroundColor(.secondary)
                                        }
                                    }
                                }
                                
                                Spacer()
                                
                                if let attendeeCount = event.attendeeCount {
                                    VStack(spacing: 4) {
                                        Text(String(attendeeCount))
                                            .font(.headline)
                                            .foregroundColor(.green)
                                        Text("Attending")
                                            .font(.caption2)
                                            .foregroundColor(.secondary)
                                    }
                                    .padding(8)
                                    .background(Color.green.opacity(0.1))
                                    .cornerRadius(6)
                                }
                            }
                            
                            if let description = event.description, !description.isEmpty {
                                Text(description)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                                    .lineLimit(2)
                            }
                            
                            if let location = event.location, !location.isEmpty {
                                Label(location, systemImage: "mappin")
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                        }
                        .padding(.vertical, 8)
                    }
                }
            }
            .navigationTitle("Events")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button(action: loadEvents) {
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
                loadEvents()
            }
        }
    }
    
    private func loadEvents() {
        isLoading = true
        errorMessage = nil
        
        Task {
            do {
                let apiService = APIService.shared
                self.events = try await apiService.getEvents()
                
                // Sort by start time
                self.events.sort { $0.startTime < $1.startTime }
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

#Preview {
    EventsListView()
        .preferredColorScheme(.dark)
}
