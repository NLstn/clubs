import { useState, useEffect } from 'react';
import { useAuth } from './useAuth';

export interface DashboardNews {
  id: string;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
  club_name: string;
  club_id: string;
}

export interface DashboardEvent {
  id: string;
  name: string;
  start_time: string;
  end_time: string;
  club_name: string;
  club_id: string;
  user_rsvp?: {
    response: string;
  };
}

export const useDashboardData = () => {
  const { api } = useAuth();
  const [news, setNews] = useState<DashboardNews[]>([]);
  const [events, setEvents] = useState<DashboardEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const [newsResponse, eventsResponse] = await Promise.all([
        api.get('/api/v1/dashboard/news'),
        api.get('/api/v1/dashboard/events')
      ]);
      
      setNews(newsResponse.data || []);
      setEvents(eventsResponse.data || []);
    } catch (error) {
      console.error('Error fetching dashboard data:', error);
      setError('Failed to load dashboard data');
      setNews([]);
      setEvents([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      
      try {
        const [newsResponse, eventsResponse] = await Promise.all([
          api.get('/api/v1/dashboard/news'),
          api.get('/api/v1/dashboard/events')
        ]);
        
        setNews(newsResponse.data || []);
        setEvents(eventsResponse.data || []);
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
        setError('Failed to load dashboard data');
        setNews([]);
        setEvents([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [api]);

  return {
    news,
    events,
    loading,
    error,
    refetch: fetchDashboardData
  };
};