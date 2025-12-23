import SwiftUI

// MARK: - App Info
// App Name: Clubs Manager
// Bundle Identifier: com.clubs.ios
// Minimum iOS: 15.0
// Target iOS: 17.0

// MARK: - Quick Reference

/*
 GETTING STARTED:
 ================
 
 1. Update API Base URL (if needed):
    - Edit APIService.swift
    - Change baseURL from "http://localhost:8080/api/v1" to your backend server
 
 2. Build & Run:
    - Open ios.xcodeproj in Xcode
    - Select iOS Simulator or Device
    - Press Cmd+R to run
 
 3. Test Login:
    - Use any email address
    - Check backend console for magic link token
    - Use that token to verify login in the app
 
 ARCHITECTURE OVERVIEW:
 ======================
 
 Models (APIModels.swift):
   - User, Club, Event, Fine, ClubMember
   - AuthResponse, AuthToken
   - Custom Codable implementations
 
 Services:
   - AuthenticationManager: Handles auth state, token refresh
   - APIService: All API communication, JWT injection
   - KeychainService: Secure token storage
 
 Views:
   - LoginView: Magic link authentication
   - MainTabView: Tab bar navigation
   - ClubsListView: Club listing and details
   - EventsListView: Event listing
   - FinesListView: Fine listing with filtering
   - ProfileView: User profile and logout
 
 COMMON TASKS:
 =============
 
 Adding a new API endpoint:
   1. Add method to APIService.swift
   2. Use makeAuthorizedRequest() helper for auth
   3. Return proper Codable type
 
 Adding a new view:
   1. Create in appropriate Views/{Category} folder
   2. Use @EnvironmentObject var authManager for auth state
   3. Use APIService.shared for API calls
   4. Follow existing component patterns
 
 Debugging:
   1. Check API base URL configuration
   2. Verify backend is running on correct port
   3. Check Console output for errors
   4. Use debugger breakpoints
   5. Test with Simulator first
 
 STYLING GUIDE:
 ==============
 
 Colors:
   - Use Color.clubsGreen for primary actions
   - Use Color.clubsBlue for secondary actions
   - Use Color.clubsRed for destructive actions
   - Use system colors for UI elements
 
 Spacing:
   - Small: 8px (padding 8)
   - Medium: 16px (padding 16)
   - Large: 24px (padding 24)
 
 Corner Radius:
   - Buttons/Inputs: 8
   - Cards: 12
   - Avatars: 50% (via .clipShape(Circle()))
 
 DEPLOYMENT CHECKLIST:
 ======================
 
 [ ] Update baseURL in APIService for production
 [ ] Set DEVELOPMENT_TEAM in build settings
 [ ] Configure App Icon in Assets.xcassets
 [ ] Update app version in build settings
 [ ] Run on physical device for testing
 [ ] Test all authentication flows
 [ ] Test with real backend API
 [ ] Check for memory leaks with Instruments
 [ ] Test with iOS 15+ devices
 [ ] Configure push notifications (if needed)
 [ ] Create App Store listing
 
 SUPPORT:
 ========
 
 For issues:
   1. Check iOS Simulator console
   2. Review error messages in app
   3. Check backend API logs
   4. Verify network connectivity
   5. Test with fresh login
 */
