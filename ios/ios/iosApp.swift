//
//  iosApp.swift
//  ios
//
//  Created by Niklas on 23.12.25.
//

import SwiftUI

@main
struct iosApp: App {
    @StateObject private var authManager = AuthenticationManager()
    
    var body: some Scene {
        WindowGroup {
            Group {
                if authManager.isAuthenticated {
                    MainTabView()
                        .environmentObject(authManager)
                } else {
                    LoginView()
                        .environmentObject(authManager)
                }
            }
            .preferredColorScheme(.dark)
            .onContinueUserActivity(NSUserActivityTypeBrowsingWeb) { activity in
                guard let url = activity.webpageURL else { return }
                if url.path.hasPrefix("/auth/magic") {
                    if let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
                       let token = components.queryItems?.first(where: { $0.name == "token" })?.value,
                       !token.isEmpty {
                        Task { await authManager.verifyMagicLink(token: token) }
                    }
                }
            }
        }
    }
}
