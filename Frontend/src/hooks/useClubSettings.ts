import { useState, useEffect } from 'react';
import api from '../utils/api';
import { parseODataCollection, type ODataCollectionResponse } from '../utils/odata';

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
                // OData v2: Query ClubSettings filtered by club ID
                const response = await api.get<ODataCollectionResponse<ClubSettings>>(`/api/v2/ClubSettings?$filter=ClubID eq '${clubId}'`);
                const settingsData = parseODataCollection(response.data);
                if (settingsData.length > 0) {
                    setSettings(settingsData[0]);
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
            // OData v2: Query ClubSettings filtered by club ID
            const response = await api.get<ODataCollectionResponse<ClubSettings>>(`/api/v2/ClubSettings?$filter=ClubID eq '${clubId}'`);
            const settingsData = parseODataCollection(response.data);
            if (settingsData.length > 0) {
                setSettings(settingsData[0]);
            } else{
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