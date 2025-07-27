import { useState, useEffect } from 'react';
import api from '../utils/api';

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
                const response = await api.get(`/api/v1/clubs/${clubId}/settings`);
                setSettings(response.data);
                setError(null);
            } catch (err: unknown) {
                console.error('Error fetching club settings:', err);
                // If settings don't exist, assume defaults (both features enabled)
                setSettings({
                id: '',
                clubId: clubId,
                finesEnabled: true,
                shiftsEnabled: true,
                teamsEnabled: true,
                membersListVisible: true,
                createdAt: '',
                createdBy: '',
                updatedAt: '',
                updatedBy: ''
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
            const response = await api.get(`/api/v1/clubs/${clubId}/settings`);
            setSettings(response.data);
            setError(null);
        } catch (err: unknown) {
            console.error('Error fetching club settings:', err);
            // If settings don't exist, assume defaults (all features enabled)
            setSettings({
                id: '',
                clubId: clubId,
                finesEnabled: true,
                shiftsEnabled: true,
                teamsEnabled: true,
                membersListVisible: true,
                createdAt: '',
                createdBy: '',
                updatedAt: '',
                updatedBy: ''
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