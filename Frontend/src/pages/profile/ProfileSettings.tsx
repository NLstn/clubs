import { useNavigate } from 'react-router-dom';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useT } from '../../hooks/useTranslation';
import { SettingsList, SettingsListSection, SettingsListItem } from '@/components/ui';
import './Profile.css';

const ProfileSettings = () => {
    const { t } = useT();
    const navigate = useNavigate();

    return (
        <Layout title={t('preferences.title')}>
            <ProfileContentLayout title={t('preferences.title')}>
                <SettingsList>
                    <SettingsListSection>
                        <SettingsListItem
                            title={t('preferences.language')}
                            subtitle={t('preferences.preferredLanguage')}
                            icon={<span aria-hidden="true">🌐</span>}
                            onClick={() => navigate('/settings/language')}
                        />
                        <SettingsListItem
                            title={t('preferences.appearance')}
                            subtitle={t('preferences.themeDescription')}
                            icon={<span aria-hidden="true">🎨</span>}
                            onClick={() => navigate('/settings/appearance')}
                        />
                    </SettingsListSection>
                </SettingsList>
            </ProfileContentLayout>
        </Layout>
    );
};

export default ProfileSettings;
