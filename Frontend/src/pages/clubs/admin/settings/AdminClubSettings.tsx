import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';

interface ClubSettings {
    id: string;
    clubId: string;
    finesEnabled: boolean;
    shiftsEnabled: boolean;
    createdAt: string;
    createdBy: string;
    updatedAt: string;
    updatedBy: string;
}

interface AdminClubSettingsProps {
    onSettingsUpdate?: () => void;
}

const AdminClubSettings = ({ onSettingsUpdate }: AdminClubSettingsProps) => {
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
                setError('Failed to load club settings');
            } finally {
                setLoading(false);
            }
        };

        fetchSettings();
    }, [id]);

    const updateSettings = async (newSettings: Partial<Pick<ClubSettings, 'finesEnabled' | 'shiftsEnabled'>>) => {
        if (!settings) return;

        try {
            setSaving(true);
            const completeSettings = {
                finesEnabled: newSettings.finesEnabled ?? settings.finesEnabled,
                shiftsEnabled: newSettings.shiftsEnabled ?? settings.shiftsEnabled,
            };
            await api.post(`/api/v1/clubs/${id}/settings`, completeSettings);
            setSettings({ ...settings, ...completeSettings });
            setError(null);
            // Notify parent component that settings have been updated
            onSettingsUpdate?.();
        } catch (err: unknown) {
            console.error('Error updating club settings:', err);
            setError('Failed to update settings');
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

    if (loading) return <div>Loading settings...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!settings) return <div>Settings not found</div>;

    return (
        <div className="club-settings">
            <h3>Club Settings</h3>
            <p>Configure which features are enabled for your club.</p>
            
            {error && <div className="error">{error}</div>}
            
            <div className="settings-section">
                <div className="setting-item">
                    <div className="setting-info">
                        <h4>Fines</h4>
                        <p>Allow club admins to create and manage fines for members.</p>
                    </div>
                    <label className="toggle-switch">
                        <input
                            type="checkbox"
                            checked={settings.finesEnabled}
                            onChange={handleFinesToggle}
                            disabled={saving}
                        />
                        <span className="slider"></span>
                    </label>
                </div>

                <div className="setting-item">
                    <div className="setting-info">
                        <h4>Shifts in Events</h4>
                        <p>Enable shift scheduling and management for club events.</p>
                    </div>
                    <label className="toggle-switch">
                        <input
                            type="checkbox"
                            checked={settings.shiftsEnabled}
                            onChange={handleShiftsToggle}
                            disabled={saving}
                        />
                        <span className="slider"></span>
                    </label>
                </div>
            </div>

            {saving && <div className="saving-indicator">Saving...</div>}
        </div>
    );
};

export default AdminClubSettings;