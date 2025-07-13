import { useState, useEffect, useCallback } from 'react';
import api from '../utils/api';

export const useOwnerCount = (clubId: string) => {
    const [ownerCount, setOwnerCount] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const fetchOwnerCount = useCallback(async () => {
        if (!clubId) return;
        
        setLoading(true);
        setError(null);
        
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/ownerCount`);
            setOwnerCount(response.data.ownerCount);
        } catch (err) {
            setError('Failed to fetch owner count');
            console.error('Error fetching owner count:', err);
        } finally {
            setLoading(false);
        }
    }, [clubId]);

    useEffect(() => {
        fetchOwnerCount();
    }, [fetchOwnerCount]);

    return { ownerCount, loading, error, refetch: fetchOwnerCount };
};
