import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { ToggleSwitch } from '../../../../components/ToggleSwitch';
import './AdminClubSettings.css';

interface ClubSettings {
    id: string;
    clubId: string;
    finesEnabled: boolean;
    shiftsEnabled: boolean;
    teamsEnabled: boolean;
    membersListVisible: boolean;
    createdAt: string;
    createdBy: string;
    updatedAt: string;
    updatedBy: string;
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
                const response = await api.get(`/api/v1/clubs/${id}/settings`);
                setSettings(response.data);
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

    const updateSettings = async (newSettings: Partial<Pick<ClubSettings, 'finesEnabled' | 'shiftsEnabled' | 'teamsEnabled' | 'membersListVisible'>>) => {
        if (!settings) return;

        try {
            setSaving(true);
            const completeSettings = {
                finesEnabled: newSettings.finesEnabled ?? settings.finesEnabled,
                shiftsEnabled: newSettings.shiftsEnabled ?? settings.shiftsEnabled,
                teamsEnabled: newSettings.teamsEnabled ?? settings.teamsEnabled,
                membersListVisible: newSettings.membersListVisible ?? settings.membersListVisible,
            };
            await api.post(`/api/v1/clubs/${id}/settings`, completeSettings);
            setSettings({ ...settings, ...completeSettings });
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

    const handleFinesToggle = async () => {
        if (!settings) return;
        await updateSettings({ finesEnabled: !settings.finesEnabled });
    };

    const handleShiftsToggle = async () => {
        if (!settings) return;
        await updateSettings({ shiftsEnabled: !settings.shiftsEnabled });
    };

    const handleTeamsToggle = async () => {
        if (!settings) return;
        await updateSettings({ teamsEnabled: !settings.teamsEnabled });
    };

    const handleMembersListToggle = async () => {
        if (!settings) return;
        await updateSettings({ membersListVisible: !settings.membersListVisible });
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
                        checked={settings.finesEnabled}
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
                        checked={settings.shiftsEnabled}
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
                        checked={settings.teamsEnabled}
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
                        checked={settings.membersListVisible}
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