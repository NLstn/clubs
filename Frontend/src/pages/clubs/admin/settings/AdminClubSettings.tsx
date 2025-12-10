import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { ToggleSwitch } from '@/components/ui';
import './AdminClubSettings.css';

interface ClubSettings {
    ID: string;
    ClubID: string;
    FinesEnabled: boolean;
    ShiftsEnabled: boolean;
    TeamsEnabled: boolean;
    MembersListVisible: boolean;
    CreatedAt: string;
    CreatedBy: string;
    UpdatedAt: string;
    UpdatedBy: string;
}

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
                const response = await api.get(`/api/v2/ClubSettings?$filter=ClubID eq '${id}'`);
                const settingsData = response.data.value?.[0];
                if (settingsData) {
                    setSettings(settingsData);
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

    const updateSettings = async (newSettings: Partial<Pick<ClubSettings, 'FinesEnabled' | 'ShiftsEnabled' | 'TeamsEnabled' | 'MembersListVisible'>>) => {
        if (!settings) return;

        try {
            setSaving(true);
            const completeSettings = {
                FinesEnabled: newSettings.FinesEnabled ?? settings.FinesEnabled,
                ShiftsEnabled: newSettings.ShiftsEnabled ?? settings.ShiftsEnabled,
                TeamsEnabled: newSettings.TeamsEnabled ?? settings.TeamsEnabled,
                MembersListVisible: newSettings.MembersListVisible ?? settings.MembersListVisible,
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
            
            // Update local state
            const updatedSettings = {
                finesEnabled: newSettings.FinesEnabled ?? settings.FinesEnabled,
                shiftsEnabled: newSettings.ShiftsEnabled ?? settings.ShiftsEnabled,
                teamsEnabled: newSettings.TeamsEnabled ?? settings.TeamsEnabled,
                membersListVisible: newSettings.MembersListVisible ?? settings.MembersListVisible,
            };
            setSettings({ ...settings, ...updatedSettings });
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

    const handleMembersListToggle = async (checked: boolean) => {
        if (!settings) return;
        await updateSettings({ MembersListVisible: checked });
    };

    if (loading) return <div>{t('clubs.loading.settings')}</div>;
    if (error) return <div className="error">{error}</div>;
    if (!settings) return <div>{t('clubs.errors.settingsNotFound')}</div>;

    return (
        <div className="club-settings">
            <h3>{t('clubs.clubSettings')}</h3>
            <p>{t('clubs.configureFeatures')}</p>
            
            {error && <div className="error">{error}</div>}
            
            <div className="settings-section">
                <div className="setting-item">
                    <div className="setting-info">
                        <h4>{t('clubs.fines')}</h4>
                        <p>{t('clubs.finesDescription')}</p>
                    </div>
                    <ToggleSwitch
                        checked={settings.FinesEnabled}
                        onChange={handleFinesToggle}
                        disabled={saving}
                    />
                </div>

                <div className="setting-item">
                    <div className="setting-info">
                        <h4>{t('clubs.shifts')}</h4>
                        <p>{t('clubs.shiftsDescription')}</p>
                    </div>
                    <ToggleSwitch
                        checked={settings.ShiftsEnabled}
                        onChange={handleShiftsToggle}
                        disabled={saving}
                    />
                </div>

                <div className="setting-item">
                    <div className="setting-info">
                        <h4>{t('clubs.teams')}</h4>
                        <p>{t('clubs.teamsDescription')}</p>
                    </div>
                    <ToggleSwitch
                        checked={settings.TeamsEnabled}
                        onChange={handleTeamsToggle}
                        disabled={saving}
                    />
                </div>

                <div className="setting-item">
                    <div className="setting-info">
                        <h4>{t('clubs.membersList')}</h4>
                        <p>{t('clubs.membersListDescription')}</p>
                    </div>
                    <ToggleSwitch
                        checked={settings.MembersListVisible}
                        onChange={handleMembersListToggle}
                        disabled={saving}
                    />
                </div>
            </div>

            {saving && <div className="saving-indicator">Saving...</div>}
        </div>
    );
};

export default AdminClubSettings;