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
            delete: "Delete"
          },
          navigation: {
            home: "Home",
            clubs: "Clubs",
            myClubs: "My Clubs",
            profile: "Profile",
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

export default testI18n;