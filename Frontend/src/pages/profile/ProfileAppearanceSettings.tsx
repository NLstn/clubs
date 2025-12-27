import { useState } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useTheme } from "../../hooks/useTheme";
import { ThemeMode } from "../../context/ThemeContext";
import { useT } from '../../hooks/useTranslation';
import { SettingsList, SettingsListSection, SettingsListItem } from '@/components/ui';
import './Profile.css';
import './ProfilePreferences.css';

const ProfileAppearanceSettings = () => {
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
        <Layout title={t('preferences.appearance')}>
            <ProfileContentLayout title={t('preferences.appearance')}>
                {message && (
                    <div className="success-message">
                        {message}
                    </div>
                )}

                <SettingsList>
                    <SettingsListSection 
                        description={t('preferences.selectTheme')}
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
            </ProfileContentLayout>
        </Layout>
    );
};

export default ProfileAppearanceSettings;
