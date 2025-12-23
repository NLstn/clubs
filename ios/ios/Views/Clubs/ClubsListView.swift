import SwiftUI

struct ClubsListView: View {
    @EnvironmentObject var authManager: AuthenticationManager
    @State private var clubs: [Club] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var selectedClub: Club?
    
    var body: some View {
        NavigationStack {
            Group {
                if isLoading {
                    ProgressView()
                } else if clubs.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "building.2")
                            .font(.system(size: 48))
                            .foregroundColor(.gray)
                        Text("No Clubs")
                            .font(.headline)
                        Text("You haven't joined any clubs yet")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .frame(maxHeight: .infinity, alignment: .center)
                } else {
                    List(clubs) { club in
                        NavigationLink(destination: ClubDetailView(club: club)) {
                            VStack(alignment: .leading, spacing: 8) {
                                Text(club.name)
                                    .font(.headline)
                                
                                if let description = club.description, !description.isEmpty {
                                    Text(description)
                                        .font(.subheadline)
                                        .foregroundColor(.secondary)
                                        .lineLimit(2)
                                }
                                
                                HStack(spacing: 16) {
                                    Label("\(club.memberCount) members", systemImage: "person.2.fill")
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                            }
                            .padding(.vertical, 8)
                        }
                    }
                }
            }
            .navigationTitle("Clubs")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button(action: loadClubs) {
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
                loadClubs()
            }
        }
    }
    
    private func loadClubs() {
        isLoading = true
        errorMessage = nil
        
        Task {
            do {
                let apiService = APIService.shared
                self.clubs = try await apiService.getClubs()
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

// MARK: - Club Detail View

struct ClubDetailView: View {
    let club: Club
    @State private var members: [ClubMember] = []
    @State private var events: [Event] = []
    @State private var fines: [Fine] = []
    @State private var isLoading = false
    @State private var selectedTab = 0
    
    var body: some View {
        ZStack {
            if isLoading {
                ProgressView()
            } else {
                TabView(selection: $selectedTab) {
                    // Overview Tab
                    VStack(alignment: .leading, spacing: 16) {
                        ScrollView {
                            VStack(alignment: .leading, spacing: 16) {
                                if let description = club.description {
                                    VStack(alignment: .leading, spacing: 8) {
                                        Text("About")
                                            .font(.headline)
                                        Text(description)
                                            .font(.body)
                                            .foregroundColor(.secondary)
                                    }
                                    .padding(12)
                                    .background(Color(.systemGray6))
                                    .cornerRadius(8)
                                }
                                
                                VStack(alignment: .leading, spacing: 8) {
                                    Text("Statistics")
                                        .font(.headline)
                                    
                                    HStack(spacing: 16) {
                                        StatCard(value: "\(club.memberCount)", label: "Members")
                                        StatCard(value: "\(events.count)", label: "Events")
                                        StatCard(value: "\(fines.count)", label: "Fines")
                                    }
                                }
                            }
                            .padding(16)
                        }
                    }
                    .tag(0)
                    
                    // Members Tab
                    VStack {
                        if members.isEmpty {
                            Text("No members")
                                .foregroundColor(.secondary)
                        } else {
                            List(members) { member in
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(member.userName)
                                        .font(.headline)
                                    Text(member.email)
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                    Text(member.role)
                                        .font(.caption2)
                                        .padding(.vertical, 2)
                                        .padding(.horizontal, 6)
                                        .background(Color.blue.opacity(0.2))
                                        .cornerRadius(4)
                                }
                            }
                        }
                    }
                    .tag(1)
                    
                    // Events Tab
                    VStack {
                        if events.isEmpty {
                            Text("No events")
                                .foregroundColor(.secondary)
                        } else {
                            List(events) { event in
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(event.title)
                                        .font(.headline)
                                    Text(event.formattedDate)
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                            }
                        }
                    }
                    .tag(2)
                }
                .tabViewStyle(.page(indexDisplayMode: .never))
            }
        }
        .navigationTitle(club.name)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .bottomBar) {
                HStack(spacing: 24) {
                    TabBarButton(icon: "info.circle", label: "Overview", isSelected: selectedTab == 0) {
                        selectedTab = 0
                    }
                    TabBarButton(icon: "person.2", label: "Members", isSelected: selectedTab == 1) {
                        selectedTab = 1
                    }
                    TabBarButton(icon: "calendar", label: "Events", isSelected: selectedTab == 2) {
                        selectedTab = 2
                    }
                }
            }
        }
        .onAppear {
            loadClubData()
        }
    }
    
    private func loadClubData() {
        isLoading = true
        let apiService = APIService.shared
        
        Task {
            do {
                async let membersTask = apiService.getClubMembers(clubId: club.id)
                async let eventsTask = apiService.getClubEvents(clubId: club.id)
                async let finesTask = apiService.getClubFines(clubId: club.id)
                
                let (membersResult, eventsResult, finesResult) = await (
                    try membersTask,
                    try eventsTask,
                    try finesTask
                )
                
                await MainActor.run {
                    self.members = membersResult
                    self.events = eventsResult
                    self.fines = finesResult
                    self.isLoading = false
                }
            } catch {
                await MainActor.run {
                    isLoading = false
                }
            }
        }
    }
}

// MARK: - Supporting Views

struct StatCard: View {
    let value: String
    let label: String
    
    var body: some View {
        VStack(spacing: 8) {
            Text(value)
                .font(.title2)
                .fontWeight(.bold)
                .foregroundColor(.green)
            Text(label)
                .font(.caption)
                .foregroundColor(.secondary)
        }
        .frame(maxWidth: .infinity)
        .padding(12)
        .background(Color(.systemGray6))
        .cornerRadius(8)
    }
}

struct TabBarButton: View {
    let icon: String
    let label: String
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            VStack(spacing: 4) {
                Image(systemName: icon)
                    .font(.system(size: 20))
                Text(label)
                    .font(.caption2)
            }
        }
        .frame(maxWidth: .infinity)
        .foregroundColor(isSelected ? .green : .secondary)
    }
}

#Preview {
    NavigationStack {
        ClubDetailView(club: Club(
            id: "1",
            name: "Tech Club",
            description: "A club for tech enthusiasts",
            memberCount: 42,
            createdAt: nil,
            imageUrl: nil
        ))
    }
    .preferredColorScheme(.dark)
}
