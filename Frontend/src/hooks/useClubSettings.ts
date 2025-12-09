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
                // OData v2: Query ClubSettings filtered by club ID
                const response = await api.get(`/api/v2/ClubSettings?$filter=ClubID eq '${clubId}'`);
                const settingsData = response.data.value || [];
                if (settingsData.length > 0) {
                    const s = settingsData[0];
                    setSettings({
                        id: s.ID,
                        clubId: s.ClubID,
                        finesEnabled: s.FinesEnabled,
                        shiftsEnabled: s.ShiftsEnabled,
                        teamsEnabled: s.TeamsEnabled,
                        membersListVisible: s.MembersListVisible,
                        createdAt: s.CreatedAt,
                        createdBy: s.CreatedBy,
                        updatedAt: s.UpdatedAt,
                        updatedBy: s.UpdatedBy
                    });
                } else {
                    // No settings found, use defaults
                    throw new Error('Settings not found');
                }
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
            // OData v2: Query ClubSettings filtered by club ID
            const response = await api.get(`/api/v2/ClubSettings?$filter=ClubID eq '${clubId}'`);
            const settingsData = response.data.value || [];
            if (settingsData.length > 0) {
                const s = settingsData[0];
                setSettings({
                    id: s.ID,
                    clubId: s.ClubID,
                    finesEnabled: s.FinesEnabled,
                    shiftsEnabled: s.ShiftsEnabled,
                    teamsEnabled: s.TeamsEnabled,
                    membersListVisible: s.MembersListVisible,
                    createdAt: s.CreatedAt,
                    createdBy: s.CreatedBy,
                    updatedAt: s.UpdatedAt,
                    updatedBy: s.UpdatedBy
                });
            } else {
                // No settings found, use defaults
                throw new Error('Settings not found');
            }
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