import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { ToggleSwitch, SettingsSection, SettingItem } from '@/components/ui';
import './AdminClubSettings.css';

interface ClubSettings {
    ID: string;
    ClubID: string;
    FinesEnabled: boolean;
    ShiftsEnabled: boolean;
    TeamsEnabled: boolean;
    NewsEnabled: boolean;
    MembersListVisible: boolean;
    DiscoverableByNonMembers: boolean;
    CreatedAt: string;
    CreatedBy: string;
    UpdatedAt: string;
    UpdatedBy: string;
}

/** Settings that can be toggled on/off in the admin settings page */
type ToggleableSettings = Pick<ClubSettings, 'FinesEnabled' | 'ShiftsEnabled' | 'TeamsEnabled' | 'NewsEnabled' | 'MembersListVisible' | 'DiscoverableByNonMembers'>;

interface AdminClubSettingsProps {
    onSettingsUpdate?: () => void;
}

const AdminClubSettings = ({ onSettingsUpdate }: AdminClubSettingsProps) => {
    const { t } = useT();
    const { id } = useParams();
    const [settings, setSettings] = useState<ClubSettings | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        const fetchSettings = async () => {
            try {
                setLoading(true);
                const response = await api.get<ClubSettings>(`/api/v2/Clubs('${id}')/Settings`);
                if (response.data) {
                    setSettings(response.data);
                }
                setError(null);
            } catch (err: unknown) {
                console.error('Error fetching club settings:', err);
                setError(t('clubs.errors.failedToLoadSettings'));
            } finally {
                setLoading(false);
            }
        };

        fetchSettings();
    }, [id, t]);

    const updateSettings = async (newSettings: Partial<ToggleableSettings>) => {
        if (!settings) return;

        try {
            setSaving(true);
            const completeSettings = {
                FinesEnabled: newSettings.FinesEnabled ?? settings.FinesEnabled,
                ShiftsEnabled: newSettings.ShiftsEnabled ?? settings.ShiftsEnabled,
                TeamsEnabled: newSettings.TeamsEnabled ?? settings.TeamsEnabled,
                NewsEnabled: newSettings.NewsEnabled ?? settings.NewsEnabled,
                MembersListVisible: newSettings.MembersListVisible ?? settings.MembersListVisible,
                DiscoverableByNonMembers: newSettings.DiscoverableByNonMembers ?? settings.DiscoverableByNonMembers,
            };
            
            if (settings.ID) {
                // Update existing settings
                await api.patch(`/api/v2/ClubSettings('${settings.ID}')`, completeSettings);
            } else {
                // Create new settings
                await api.post(`/api/v2/ClubSettings`, {
                    ...completeSettings,
                    ClubID: id
                });
            }
            
            // Update local state with PascalCase field names
            setSettings({
                ...settings,
                FinesEnabled: newSettings.FinesEnabled ?? settings.FinesEnabled,
                ShiftsEnabled: newSettings.ShiftsEnabled ?? settings.ShiftsEnabled,
                TeamsEnabled: newSettings.TeamsEnabled ?? settings.TeamsEnabled,
                NewsEnabled: newSettings.NewsEnabled ?? settings.NewsEnabled,
                MembersListVisible: newSettings.MembersListVisible ?? settings.MembersListVisible,
                DiscoverableByNonMembers: newSettings.DiscoverableByNonMembers ?? settings.DiscoverableByNonMembers,
            });
            setError(null);
            // Notify parent component that settings have been updated
            onSettingsUpdate?.();
        } catch (err: unknown) {
            console.error('Error updating club settings:', err);
            setError(t('clubs.errors.failedToUpdateSettings'));
        } finally {
            setSaving(false);
        }
    };

    const handleFinesToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ FinesEnabled: checked });
    };

    const handleShiftsToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ ShiftsEnabled: checked });
    };

    const handleTeamsToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ TeamsEnabled: checked });
    };

    const handleNewsToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ NewsEnabled: checked });
    };

    const handleMembersListToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ MembersListVisible: checked });
    };

    const handleDiscoverableToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ DiscoverableByNonMembers: checked });
    };

    if (loading) return <div>{t('clubs.loading.settings')}</div>;
    if (error) return <div className="error">{error}</div>;
    if (!settings) return <div>{t('clubs.errors.settingsNotFound')}</div>;

    return (
        <div className="club-settings">
            <h2 className="club-settings-title">{t('clubs.clubSettings')}</h2>
            <p className="club-settings-description">{t('clubs.configureFeatures')}</p>
            
            {error && <div className="error">{error}</div>}
            
            {/* Features Section */}
            <SettingsSection
                title={t('clubs.settings.features')}
                description={t('clubs.settings.featuresDescription')}
            >
                <SettingItem
                    title={t('clubs.fines')}
                    description={t('clubs.finesDescription')}
                >
                    <ToggleSwitch
                        checked={settings.FinesEnabled}
                        onChange={handleFinesToggle}
                        disabled={saving}
                    />
                </SettingItem>

                <SettingItem
                    title={t('clubs.shifts')}
                    description={t('clubs.shiftsDescription')}
                >
                    <ToggleSwitch
                        checked={settings.ShiftsEnabled}
                        onChange={handleShiftsToggle}
                        disabled={saving}
                    />
                </SettingItem>

                <SettingItem
                    title={t('clubs.teams')}
                    description={t('clubs.teamsDescription')}
                >
                    <ToggleSwitch
                        checked={settings.TeamsEnabled}
                        onChange={handleTeamsToggle}
                        disabled={saving}
                    />
                </SettingItem>

                <SettingItem
                    title={t('clubs.newsFeature')}
                    description={t('clubs.newsDescription')}
                >
                    <ToggleSwitch
                        checked={settings.NewsEnabled}
                        onChange={handleNewsToggle}
                        disabled={saving}
                    />
                </SettingItem>
            </SettingsSection>

            {/* Privacy Section */}
            <SettingsSection
                title={t('clubs.settings.privacy')}
                description={t('clubs.settings.privacyDescription')}
            >
                <SettingItem
                    title={t('clubs.membersList')}
                    description={t('clubs.membersListDescription')}
                >
                    <ToggleSwitch
                        checked={settings.MembersListVisible}
                        onChange={handleMembersListToggle}
                        disabled={saving}
                    />
                </SettingItem>

                <SettingItem
                    title={t('clubs.discoverable')}
                    description={t('clubs.discoverableDescription')}
                >
                    <ToggleSwitch
                        checked={settings.DiscoverableByNonMembers}
                        onChange={handleDiscoverableToggle}
                        disabled={saving}
                    />
                </SettingItem>
            </SettingsSection>

            {saving && <div className="saving-indicator">{t('clubs.saving')}</div>}
        </div>
    );
};

export default AdminClubSettings;