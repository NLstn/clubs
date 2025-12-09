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
        // OData v2: Query Users entity - backend hooks filter to current user
        const response = await api.get('/api/v2/Users');
        // OData returns collection with 'value' array
        const users = response.data.value || [];
        if (users.length > 0) {
          setUser(users[0]);
        } else {
          throw new Error('User not found');
        }
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
