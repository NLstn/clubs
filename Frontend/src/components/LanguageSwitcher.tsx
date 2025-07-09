import React from 'react';
import { useT } from '../hooks/useTranslation';

const LanguageSwitcher: React.FC = () => {
  const { changeLanguage, getCurrentLanguage } = useT();
  
  const languages = [
    { code: 'en', name: 'English' },
    { code: 'de', name: 'Deutsch' }
  ];
  
  const handleLanguageChange = (languageCode: string) => {
    changeLanguage(languageCode);
  };
  
  return (
    <div className="language-switcher">
      <select 
        id="language-select"
        value={getCurrentLanguage()}
        onChange={(e) => handleLanguageChange(e.target.value)}
      >
        {languages.map((lang) => (
          <option key={lang.code} value={lang.code}>
            {lang.name}
          </option>
        ))}
      </select>
    </div>
  );
};

export default LanguageSwitcher;