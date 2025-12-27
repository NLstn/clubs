import { useNavigate } from 'react-router-dom';
import Layout from "../../components/layout/Layout";
import { useT } from '../../hooks/useTranslation';
import { SettingsList, SettingsListSection, SettingsListItem } from '@/components/ui';
import './Profile.css';

const SettingsIndex = () => {
    const { t } = useT();
    const navigate = useNavigate();

    return (
        <Layout title="Settings">
            <div style={{ maxWidth: '800px', margin: '0 auto', padding: 'var(--space-lg)' }}>
                <h1 style={{ marginBottom: 'var(--space-lg)', color: 'var(--color-text)' }}>Settings</h1>
                
                <SettingsList>
                    <SettingsListSection>
                        <SettingsListItem
                            title="Profile"
                            subtitle="Personal information and account details"
                            icon={<span aria-hidden="true">👤</span>}
                            onClick={() => navigate('/settings/profile')}
                        />
                        <SettingsListItem
                            title={t('preferences.title')}
                            subtitle="Language and appearance settings"
                            icon={<span aria-hidden="true">⚙️</span>}
                            onClick={() => navigate('/settings/preferences')}
                        />
                        <SettingsListItem
                            title="Privacy"
                            subtitle="Manage your privacy settings"
                            icon={<span aria-hidden="true">🔒</span>}
                            onClick={() => navigate('/settings/privacy')}
                        />
                        <SettingsListItem
                            title="Invites"
                            subtitle="View and manage club invitations"
                            icon={<span aria-hidden="true">✉️</span>}
                            onClick={() => navigate('/settings/invites')}
                        />
                        <SettingsListItem
                            title="Fines"
                            subtitle="View your fines and payment history"
                            icon={<span aria-hidden="true">💰</span>}
                            onClick={() => navigate('/settings/fines')}
                        />
                        <SettingsListItem
                            title="Shifts"
                            subtitle="Manage your shift assignments"
                            icon={<span aria-hidden="true">📅</span>}
                            onClick={() => navigate('/settings/shifts')}
                        />
                        <SettingsListItem
                            title="Sessions"
                            subtitle="Active sessions and security"
                            icon={<span aria-hidden="true">🔐</span>}
                            onClick={() => navigate('/settings/sessions')}
                        />
                        <SettingsListItem
                            title="API Keys"
                            subtitle="Manage API keys for integrations"
                            icon={<span aria-hidden="true">🔑</span>}
                            onClick={() => navigate('/settings/api-keys')}
                        />
                        <SettingsListItem
                            title="Notifications"
                            subtitle="Configure notification preferences"
                            icon={<span aria-hidden="true">🔔</span>}
                            onClick={() => navigate('/settings/notifications')}
                        />
                    </SettingsListSection>
                </SettingsList>
            </div>
        </Layout>
    );
};

export default SettingsIndex;
