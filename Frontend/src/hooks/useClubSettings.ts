import { useState, useEffect } from 'react';
import api from '../utils/api';

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

interface UseClubSettingsResult {
    settings: ClubSettings | null;
    loading: boolean;
    error: string | null;
    refetch: () => Promise<void>;
}

export const useClubSettings = (clubId: string | undefined): UseClubSettingsResult => {
    const [settings, setSettings] = useState<ClubSettings | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchSettings = async () => {
            if (!clubId) {
                setLoading(false);
                return;
            }

            try {
                setLoading(true);
                // OData v2: Use navigation property from Club to Settings
                const response = await api.get<ClubSettings>(`/api/v2/Clubs('${clubId}')/Settings`);
                if (response.data) {
                    setSettings(response.data);
                } else {
                    // No settings found, use defaults
                    throw new Error('Settings not found');
                }
                setError(null);
            } catch (err: unknown) {
                console.error('Error fetching club settings:', err);
                // If settings don't exist, assume defaults (both features enabled)
                setSettings({
                ID: '',
                ClubID: clubId,
                FinesEnabled: true,
                ShiftsEnabled: true,
                TeamsEnabled: true,
                MembersListVisible: true,
                CreatedAt: '',
                CreatedBy: '',
                UpdatedAt: '',
                UpdatedBy: ''
            });
                setError(null);
            } finally {
                setLoading(false);
            }
        };

        fetchSettings();
    }, [clubId]);

    const refetch = async () => {
        if (!clubId) return;

        try {
            setLoading(true);
            // OData v2: Use navigation property from Club to Settings
            const response = await api.get<ClubSettings>(`/api/v2/Clubs('${clubId}')/Settings`);
            if (response.data) {
                setSettings(response.data);
            } else {
                // No settings found, use defaults
                throw new Error('Settings not found');
            }
            setError(null);
        } catch (err: unknown) {
            console.error('Error fetching club settings:', err);
            // If settings don't exist, assume defaults (all features enabled)
            setSettings({
                ID: '',
                ClubID: clubId,
                FinesEnabled: true,
                ShiftsEnabled: true,
                TeamsEnabled: true,
                MembersListVisible: true,
                CreatedAt: '',
                CreatedBy: '',
                UpdatedAt: '',
                UpdatedBy: ''
            });
            setError(null);
        } finally {
            setLoading(false);
        }
    };

    return {
        settings,
        loading,
        error,
        refetch
    };
};