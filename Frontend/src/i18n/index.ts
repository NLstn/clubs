import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import en from './locales/en.json';
import de from './locales/de.json';

const resources = {
  en: {
    translation: en
  },
  de: {
    translation: de
  }
};

if (typeof window !== 'undefined') {
  i18n
    .use(LanguageDetector)
    .use(initReactI18next)
    .init({
      resources,
      fallbackLng: 'en',
      debug: false,

      detection: {
        order: ['localStorage', 'navigator', 'htmlTag'],
        caches: ['localStorage'],
        lookupLocalStorage: 'i18nextLng',
      },

      interpolation: {
        escapeValue: false,
      },
    });
} else {
  i18n.use(initReactI18next).init({
    resources,
    fallbackLng: 'en',
    debug: false,
    interpolation: {
      escapeValue: false,
    },
  });
}

export default i18n;