import { useState } from 'react';
import Layout from "../../components/layout/Layout";
import SimpleSettingsLayout from '../../components/layout/SimpleSettingsLayout';
import LanguageSwitcher from "../../components/LanguageSwitcher";
import { useTheme } from "../../hooks/useTheme";
import { ThemeMode } from "../../context/ThemeContext";
import { useT } from '../../hooks/useTranslation';
import { SettingsList, SettingsListSection, SettingsListItem } from '@/components/ui';
import './Profile.css';
import './ProfilePreferences.css';

const ProfileSettings = () => {
    const { t } = useT();
    const { theme, setTheme, effectiveTheme } = useTheme();
    const [message, setMessage] = useState('');

    const handleThemeChange = (newTheme: ThemeMode) => {
        setTheme(newTheme);
        setMessage(t('preferences.themeSaved'));
        setTimeout(() => setMessage(''), 3000);
    };

    const themeOptions: { value: ThemeMode; label: string; description: string }[] = [
        {
            value: 'light',
            label: t('preferences.lightMode'),
            description: t('preferences.lightModeDescription')
        },
        {
            value: 'dark',
            label: t('preferences.darkMode'),
            description: t('preferences.darkModeDescription')
        },
        {
            value: 'system',
            label: t('preferences.systemSetting'),
            description: t('preferences.systemSettingDescription')
        }
    ];

    return (
        <Layout title={t('preferences.title')}>
            <SimpleSettingsLayout title={t('preferences.title')}>
                {message && (
                    <div className="success-message" style={{ marginBottom: 'var(--space-md)' }}>
                        {message}
                    </div>
                )}

                <SettingsList>
                    <SettingsListSection 
                        title={t('preferences.language').toUpperCase()}
                        description={t('preferences.preferredLanguage')}
                    >
                        <div style={{ padding: 'var(--space-md)' }}>
                            <LanguageSwitcher />
                        </div>
                    </SettingsListSection>

                    <SettingsListSection 
                        title={t('preferences.appearance').toUpperCase()}
                        description={t('preferences.themeDescription')}
                    >
                        {themeOptions.map((option) => (
                            <SettingsListItem
                                key={option.value}
                                title={option.label}
                                subtitle={option.description}
                                value={theme === option.value ? '✓' : ''}
                                onClick={() => handleThemeChange(option.value)}
                                showChevron={false}
                                icon={
                                    <span aria-hidden="true">
                                        {option.value === 'light' ? '☀️' :
                                         option.value === 'dark' ? '🌙' : 'ℹ️'}
                                    </span>
                                }
                            />
                        ))}
                    </SettingsListSection>

                    <div className="preferences-current-theme-indicator">
                        <strong>{t('preferences.currentlyActive')}</strong>{' '}
                        {effectiveTheme === 'light' ? 
                            <><span aria-hidden="true">☀️</span> {t('preferences.lightTheme')}</> : 
                            <><span aria-hidden="true">🌙</span> {t('preferences.darkTheme')}</>
                        }
                        {theme === 'system' && <> ({t('preferences.automaticallySet')})</>}
                    </div>
                </SettingsList>
            </SimpleSettingsLayout>
        </Layout>
    );
};

export default ProfileSettings;
