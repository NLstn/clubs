import { useState, useEffect } from 'react';
import { useAuth } from './useAuth';
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';

export interface CurrentUser {
  ID: string;
  Email: string;
  FirstName?: string;
  LastName?: string;
  BirthDate?: string;
  SetupCompleted?: boolean;
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
        const response = await api.get<ODataCollectionResponse<CurrentUser>>('/api/v2/Users?$select=ID,Email,FirstName,LastName,BirthDate,SetupCompleted');
        // OData returns collection with 'value' array
        const users = parseODataCollection(response.data);
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
