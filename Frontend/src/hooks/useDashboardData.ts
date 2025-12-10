import { useState, useEffect } from 'react';
import { useAuth } from './useAuth';

// TimelineItem represents a unified timeline entry from the backend
export interface TimelineItem {
  ID: string;
  ClubID: string;
  ClubName: string;
  Type: string; // "activity", "event", "news"
  Title: string;
  Content?: string;
  Timestamp: string; // Unified timestamp for sorting
  CreatedAt: string;
  UpdatedAt: string;
  
  // Event-specific fields (only populated for Type="event")
  StartTime?: string;
  EndTime?: string;
  Location?: string;
  
  // Activity-specific fields (only populated for Type="activity")
  Actor?: string;
  ActorName?: string;
  
  // Metadata for additional information
  Metadata?: {
    [key: string]: unknown;
  };
  
  // RSVP information for events
  UserRSVP?: {
    ID: string;
    EventID: string;
    UserID: string;
    Response: string;
    CreatedAt: string;
    UpdatedAt: string;
  };
}

// Legacy interfaces for backward compatibility
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
  ID: string;
  Type: string; // "news", "event", "role_changed", "member_promoted", "member_demoted"
  Title: string;
  Content?: string;
  ClubName: string;
  ClubID: string;
  CreatedAt: string;
  UpdatedAt: string;
  Actor?: string;        // User ID who created/initiated the activity
  ActorName?: string;   // Name of the user who created/initiated the activity
  Metadata?: {
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

// Helper function to convert TimelineItem to legacy format
function convertTimelineToLegacy(timeline: TimelineItem[]): {
  news: DashboardNews[];
  events: DashboardEvent[];
  activities: ActivityItem[];
} {
  const news: DashboardNews[] = [];
  const events: DashboardEvent[] = [];
  const activities: ActivityItem[] = [];

  timeline.forEach(item => {
    if (item.Type === 'news') {
      news.push({
        id: item.ID,
        title: item.Title,
        content: item.Content || '',
        created_at: item.CreatedAt,
        updated_at: item.UpdatedAt,
        club_name: item.ClubName,
        club_id: item.ClubID,
      });
    } else if (item.Type === 'event') {
      events.push({
        id: item.ID,
        name: item.Title,
        start_time: item.StartTime || item.Timestamp,
        end_time: item.EndTime || item.Timestamp,
        club_name: item.ClubName,
        club_id: item.ClubID,
        user_rsvp: item.UserRSVP ? {
          response: item.UserRSVP.Response
        } : undefined,
      });
    } else if (item.Type === 'activity') {
      activities.push({
        ID: item.ID,
        Type: item.Type,
        Title: item.Title,
        Content: item.Content,
        ClubName: item.ClubName,
        ClubID: item.ClubID,
        CreatedAt: item.CreatedAt,
        UpdatedAt: item.UpdatedAt,
        Actor: item.Actor,
        ActorName: item.ActorName,
        Metadata: item.Metadata as ActivityItem['Metadata'],
      });
    }
  });

  return { news, events, activities };
}

export const useDashboardData = () => {
  const { api } = useAuth();
  const [timeline, setTimeline] = useState<TimelineItem[]>([]);
  const [news, setNews] = useState<DashboardNews[]>([]);
  const [events, setEvents] = useState<DashboardEvent[]>([]);
  const [activities, setActivities] = useState<ActivityItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      // OData v2 API: Use the unified Timeline entity
      const response = await api.get('/api/v2/TimelineItems');
      
      // OData entities return data in the 'value' property
      const timelineData = response.data.value || [];
      setTimeline(timelineData);
      
      // Convert to legacy format for backward compatibility
      const { news: newsData, events: eventsData, activities: activitiesData } = convertTimelineToLegacy(timelineData);
      setNews(newsData);
      setEvents(eventsData);
      setActivities(activitiesData);
    } catch (error) {
      console.error('Error fetching dashboard data:', error);
      setError('Failed to load dashboard data');
      setTimeline([]);
      setNews([]);
      setEvents([]);
      setActivities([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDashboardData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return {
    timeline, // New unified timeline data
    news,      // Legacy format for backward compatibility
    events,    // Legacy format for backward compatibility
    activities,// Legacy format for backward compatibility
    loading,
    error,
    refetch: fetchDashboardData
  };
};