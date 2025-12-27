import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import LanguageSwitcher from "../../components/LanguageSwitcher";
import { useT } from '../../hooks/useTranslation';
import { SettingsList, SettingsListSection } from '@/components/ui';
import './Profile.css';

const ProfileLanguageSettings = () => {
    const { t } = useT();

    return (
        <Layout title={t('preferences.language')}>
            <ProfileContentLayout title={t('preferences.language')}>
                <SettingsList>
                    <SettingsListSection 
                        description={t('preferences.preferredLanguage')}
                    >
                        <div style={{ padding: 'var(--space-md)' }}>
                            <LanguageSwitcher />
                        </div>
                    </SettingsListSection>
                </SettingsList>
            </ProfileContentLayout>
        </Layout>
    );
};

export default ProfileLanguageSettings;
