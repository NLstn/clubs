import { useState, useEffect } from 'react';
import { useAuth } from './useAuth';

export interface CurrentUser {
  ID: string;
  Email: string;
  FirstName?: string;
  LastName?: string;
  BirthDate?: string;
}

export const useCurrentUser = () => {
  const { api } = useAuth();
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCurrentUser = async () => {
      try {
        setLoading(true);
        const response = await api.get('/api/v1/me');
        setUser(response.data);
        setError(null);
      } catch (err) {
        console.error('Error fetching current user:', err);
        setError('Failed to fetch user information');
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    fetchCurrentUser();
  }, [api]);

  return { user, loading, error };
};
