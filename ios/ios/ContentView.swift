//
//  ContentView.swift
//  ios
//
//  Created by Niklas on 23.12.25.
//

import SwiftUI
import Combine

struct MainTabView: View {
    @EnvironmentObject var authManager: AuthenticationManager
    
    var body: some View {
        ZStack(alignment: .top) {
            TabView {
            ClubsListView()
                .tabItem {
                    Label("Clubs", systemImage: "building.2.fill")
                }
            
            EventsListView()
                .tabItem {
                    Label("Events", systemImage: "calendar")
                }
            
            FinesListView()
                .tabItem {
                    Label("Fines", systemImage: "exclamationmark.circle.fill")
                }
            
            ProfileView()
                .tabItem {
                    Label("Profile", systemImage: "person.fill")
                }
            }
            .environmentObject(authManager)
            
            if let message = authManager.loginSuccessMessage {
                Text(message)
                    .font(.footnote)
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                    .background(.ultraThinMaterial)
                    .cornerRadius(8)
                    .padding(.top, 8)
            }
        }
    }
}

#Preview {
    MainTabView()
        .environmentObject(AuthenticationManager())
}

