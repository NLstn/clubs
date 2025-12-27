import { I18nextProvider } from 'react-i18next';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// Create a test i18n instance with minimal configuration
const testI18n = i18n.createInstance();
testI18n
  .use(initReactI18next)
  .init({
    lng: 'en',
    fallbackLng: 'en',
    debug: false,
    interpolation: {
      escapeValue: false,
    },
    resources: {
      en: {
        translation: {
          common: {
            loading: "Loading...",
            error: "Error",
            save: "Save",
            cancel: "Cancel",
            edit: "Edit",
            delete: "Delete",
            yes: "Yes",
            no: "No"
          },
          navigation: {
            home: "Home",
            clubs: "Clubs",
            myClubs: "My Clubs",
            profile: "Profile",
            settings: "Settings",
            adminPanel: "Admin Panel",
            createNewClub: "Create New Club",
            logout: "Logout",
            login: "Login"
          },
          auth: {
            login: "Login",
            loginInstruction: "Enter your email to receive a magic link for logging in.",
            email: "Email",
            sendMagicLink: "Send Magic Link",
            sending: "Sending...",
            checkEmail: "Check your email for a login link!",
            networkError: "Network error. Please try again.",
            redirectMessage: "Please log in to continue to your requested page."
          },
          profile: {
            profile: "Profile",
            language: "Language",
            selectLanguage: "Select Language"
          },
          dashboard: {
            title: "Dashboard",
            activityFeed: "Activity Feed",
            loadingDashboard: "Loading dashboard...",
            loadingMore: "Loading more activities...",
            failedToLoad: "Failed to load dashboard data",
            noActivities: "No recent activities from your clubs.",
            postedOn: "Posted on",
            createdBy: "Created by",
            by: "by",
            yourRsvp: "Your RSVP"
          },
          recentClubs: {
            title: "Recent Clubs",
            recentClubs: "Recent clubs",
            noRecentClubs: "No recent clubs",
            viewAllClubs: "View All Clubs",
            removeFromRecent: "Remove from recent clubs"
          },
          clubList: {
            myClubs: "My Clubs",
            clubsIManage: "Clubs I Manage",
            clubsImMemberOf: "Clubs I'm a Member Of",
            myTeams: "My Teams",
            deleted: "Deleted",
            noClubsYet: "No Clubs Yet",
            notMemberYet: "You're not a member of any clubs yet.",
            createFirstClub: "Create Your First Club",
            loadingClubs: "Loading clubs...",
            failedToFetch: "Failed to fetch clubs"
          },
          createClub: {
            title: "Create New Club",
            clubName: "Club Name:",
            description: "Description:",
            createClub: "Create Club",
            successMessage: "Club created successfully!",
            errorMessage: "Error creating club"
          },
          shifts: {
            myFutureShifts: "My Future Shifts",
            loadingShifts: "Loading shifts...",
            failedToLoad: "Failed to load shifts",
            noUpcomingShifts: "No upcoming shifts found.",
            checkBackLater: "Check back later or contact your club administrators if you expect to see shifts here.",
            time: "Time",
            location: "Location",
            teamMembers: "Team Members",
            noOtherMembers: "No other members assigned",
            moreMembers: "+{{count}} more"
          }
        }
      }
    }
  });

export const TestI18nProvider = ({ children }: { children: React.ReactNode }) => (
  <I18nextProvider i18n={testI18n}>
    {children}
  </I18nextProvider>
);