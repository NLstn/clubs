import { useTranslation } from 'react-i18next';

export const useT = () => {
  const { t, i18n } = useTranslation();
  
  const changeLanguage = (language: string) => {
    i18n.changeLanguage(language);
  };
  
  const getCurrentLanguage = () => {
    return i18n.language;
  };
  
  return {
    t,
    changeLanguage,
    getCurrentLanguage,
    language: i18n.language
  };
};