import { useState, useEffect, useCallback } from 'react';
import api from '../utils/api';
import { parseODataCollection, type ODataCollectionResponse } from '../utils/odata';

export interface Notification {
  ID: string;
  UserID: string;
  Type: string;
  Title: string;
  Message: string;
  Read: boolean;
  CreatedAt: string;
  ClubID?: string;
  EventID?: string;
  FineID?: string;
  InviteID?: string;
}

export interface NotificationPreferences {
  ID: string;
  UserID: string;
  MemberAddedInApp: boolean;
  MemberAddedEmail: boolean;
  EventCreatedInApp: boolean;
  EventCreatedEmail: boolean;
  FineAssignedInApp: boolean;
  FineAssignedEmail: boolean;
  NewsCreatedInApp: boolean;
  NewsCreatedEmail: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

export const useNotifications = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchNotifications = useCallback(async (limit?: number) => {
    try {
      setLoading(true);
      setError(null);
      // OData v2: Query Notifications with optional top limit, ordered by creation date
      const topParam = limit ? `&$top=${limit}` : '';
      const response = await api.get<ODataCollectionResponse<ODataNotification>>(`/api/v2/Notifications?$orderby=CreatedAt desc${topParam}`);
      interface ODataNotification {
        ID: string;
        UserID: string;
        Type: string;
        Title: string;
        Message: string;
        Read?: boolean;
        IsRead?: boolean;
        CreatedAt: string;
        ClubID?: string;
        EventID?: string;
        FineID?: string;
        InviteID?: string;
      }
      const notificationsData = parseODataCollection(response.data);
      // Map OData response to match expected format
      const mappedNotifications = notificationsData.map((n: ODataNotification) => ({
        ID: n.ID,
        UserID: n.UserID,
        Type: n.Type,
        Title: n.Title,
        Message: n.Message,
        Read: n.Read ?? n.IsRead ?? false,
        CreatedAt: n.CreatedAt,
        ClubID: n.ClubID,
        EventID: n.EventID,
        FineID: n.FineID,
        InviteID: n.InviteID
      }));
      setNotifications(mappedNotifications);
    } catch (err) {
      setError('Failed to fetch notifications');
      console.error('Error fetching notifications:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchUnreadCount = useCallback(async () => {
    try {
      // OData v2: Use $count to get unread notifications count
      const response = await api.get('/api/v2/Notifications/$count?$filter=Read eq false');
      // OData $count returns a plain number
      const count = typeof response.data === 'number' ? response.data : parseInt(response.data, 10);
      setUnreadCount(count);
    } catch (err) {
      console.error('Error fetching notification count:', err);
    }
  }, []);

  const markAsRead = useCallback(async (notificationId: string) => {
    try {
      // OData v2: Use MarkAsRead action on Notification entity
      await api.post(`/api/v2/Notifications('${notificationId}')/MarkAsRead`);
      
      // Update local state
      setNotifications(prev =>
        prev.map(notification =>
          notification.ID === notificationId ? { ...notification, Read: true } : notification
        )
      );
      
      // Update unread count
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch (err) {
      console.error('Error marking notification as read:', err);
    }
  }, []);

  const markAllAsRead = useCallback(async () => {
    try {
      // OData v2: Use MarkAllNotificationsRead unbound action
      await api.post('/api/v2/MarkAllNotificationsRead');
      
      // Update local state
      setNotifications(prev =>
        prev.map(notification => ({ ...notification, Read: true }))
      );
      
      // Reset unread count
      setUnreadCount(0);
    } catch (err) {
      console.error('Error marking all notifications as read:', err);
    }
  }, []);

  const deleteNotification = useCallback(async (notificationId: string) => {
    try {
      // OData v2: DELETE notification entity
      await api.delete(`/api/v2/Notifications('${notificationId}')`);
      
      // Update local state
      setNotifications(prev => prev.filter(n => n.ID !== notificationId));

      // Update unread count if notification was unread
      const notification = notifications.find(n => n.ID === notificationId);
      if (notification && !notification.Read) {
        setUnreadCount(prev => Math.max(0, prev - 1));
      }
    } catch (err) {
      console.error('Error deleting notification:', err);
    }
  }, [notifications]);

  const fetchPreferences = useCallback(async () => {
    try {
      // OData v2: Query UserNotificationPreferences for current user
      const response = await api.get<ODataCollectionResponse<NotificationPreferences>>('/api/v2/UserNotificationPreferences');
      const prefsData = parseODataCollection(response.data);
      return prefsData.length > 0 ? prefsData[0] : null;
    } catch (err) {
      console.error('Error fetching notification preferences:', err);
      return null;
    }
  }, []);

  const updatePreferences = useCallback(async (updates: Partial<NotificationPreferences>) => {
    try {
      // First get current user preferences
      const currentPrefs = await fetchPreferences();
      if (!currentPrefs) {
        throw new Error('Preferences not found');
      }
      // OData v2: PATCH to update preferences
      const response = await api.patch(`/api/v2/UserNotificationPreferences('${currentPrefs.ID}')`, updates);
      return response.data;
    } catch (err) {
      console.error('Error updating notification preferences:', err);
      throw err;
    }
  }, [fetchPreferences]);

  const refreshNotifications = useCallback(() => {
    fetchNotifications();
    fetchUnreadCount();
  }, [fetchNotifications, fetchUnreadCount]);

  // Initial load
  useEffect(() => {
    refreshNotifications();
  }, [refreshNotifications]);

  // Periodic refresh for unread count (every 30 seconds)
  useEffect(() => {
    const interval = setInterval(() => {
      fetchUnreadCount();
    }, 30000);

    return () => clearInterval(interval);
  }, [fetchUnreadCount]);

  return {
    notifications,
    unreadCount,
    loading,
    error,
    fetchNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    refreshNotifications,
    fetchPreferences,
    updatePreferences
  };
};

export const useNotificationPreferences = () => {
  const [preferences, setPreferences] = useState<NotificationPreferences | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchPreferences = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      // OData v2: Query UserNotificationPreferences for current user
      const response = await api.get<ODataCollectionResponse<NotificationPreferences>>('/api/v2/UserNotificationPreferences');
      const prefsData = parseODataCollection(response.data);
      if (prefsData.length > 0) {
        setPreferences(prefsData[0]);
      }
    } catch (err) {
      setError('Failed to fetch notification preferences');
      console.error('Error fetching notification preferences:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const updatePreferences = useCallback(async (updates: Partial<NotificationPreferences>) => {
    try {
      setLoading(true);
      setError(null);
      // First get current preferences to get ID
      const currentResponse = await api.get<ODataCollectionResponse<NotificationPreferences>>('/api/v2/UserNotificationPreferences');
      const currentPrefs = parseODataCollection(currentResponse.data)[0];
      if (!currentPrefs) {
        throw new Error('Preferences not found');
      }
      // OData v2: PATCH to update preferences
      const response = await api.patch(`/api/v2/UserNotificationPreferences('${currentPrefs.ID}')`, updates);
      setPreferences(response.data);
      return response.data;
    } catch (err) {
      setError('Failed to update notification preferences');
      console.error('Error updating notification preferences:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Initial load
  useEffect(() => {
    fetchPreferences();
  }, [fetchPreferences]);

  return {
    preferences,
    loading,
    error,
    fetchPreferences,
    updatePreferences,
  };
};