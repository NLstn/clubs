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

export interface ActivityItem {
  id: string;
  type: string; // "news", "event", "role_changed", "member_promoted", "member_demoted"
  title: string;
  content?: string;
  club_name: string;
  club_id: string;
  created_at: string;
  updated_at: string;
  actor?: string;        // User ID who created/initiated the activity
  actor_name?: string;   // Name of the user who created/initiated the activity
  metadata?: {
    start_time?: string;
    end_time?: string;
    user_rsvp?: {
      response: string;
    };
    old_role?: string;
    new_role?: string;
    club_name?: string;
    affected_user_id?: string; // User ID of the person whose role was changed
    [key: string]: unknown;
  }; // For extensibility
}

export const useDashboardData = () => {
  const { api } = useAuth();
  const [news, setNews] = useState<DashboardNews[]>([]);
  const [events, setEvents] = useState<DashboardEvent[]>([]);
  const [activities, setActivities] = useState<ActivityItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const [newsResponse, eventsResponse, activitiesResponse] = await Promise.all([
        api.get('/api/v1/dashboard/news'),
        api.get('/api/v1/dashboard/events'),
        api.get('/api/v1/dashboard/activities')
      ]);
      
      setNews(newsResponse.data || []);
      setEvents(eventsResponse.data || []);
      setActivities(activitiesResponse.data || []);
    } catch (error) {
      console.error('Error fetching dashboard data:', error);
      setError('Failed to load dashboard data');
      setNews([]);
      setEvents([]);
      setActivities([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      
      try {
        const [newsResponse, eventsResponse, activitiesResponse] = await Promise.all([
          api.get('/api/v1/dashboard/news'),
          api.get('/api/v1/dashboard/events'),
          api.get('/api/v1/dashboard/activities')
        ]);
        
        setNews(newsResponse.data || []);
        setEvents(eventsResponse.data || []);
        setActivities(activitiesResponse.data || []);
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
        setError('Failed to load dashboard data');
        setNews([]);
        setEvents([]);
        setActivities([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [api]);

  return {
    news,
    events,
    activities,
    loading,
    error,
    refetch: fetchDashboardData
  };
};