import { useState, useEffect, useCallback } from 'react';
import api from '../utils/api';

export interface Notification {
  id: string;
  userId: string;
  type: string;
  title: string;
  message: string;
  read: boolean;
  createdAt: string;
  clubId?: string;
  eventId?: string;
  fineId?: string;
  inviteId?: string;
}

export interface NotificationPreferences {
  id: string;
  userId: string;
  memberAddedInApp: boolean;
  memberAddedEmail: boolean;
  eventCreatedInApp: boolean;
  eventCreatedEmail: boolean;
  fineAssignedInApp: boolean;
  fineAssignedEmail: boolean;
  newsCreatedInApp: boolean;
  newsCreatedEmail: boolean;
  createdAt: string;
  updatedAt: string;
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
      const params = limit ? `?limit=${limit}` : '';
      const response = await api.get(`/api/v1/notifications${params}`);
      setNotifications(response.data);
    } catch (err) {
      setError('Failed to fetch notifications');
      console.error('Error fetching notifications:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchUnreadCount = useCallback(async () => {
    try {
      const response = await api.get('/api/v1/notifications/count');
      setUnreadCount(response.data.count);
    } catch (err) {
      console.error('Error fetching notification count:', err);
    }
  }, []);

  const markAsRead = useCallback(async (notificationId: string) => {
    try {
      await api.put(`/api/v1/notifications/${notificationId}`);
      
      // Update local state
      setNotifications(prev => 
        prev.map(notification => 
          notification.id === notificationId 
            ? { ...notification, read: true }
            : notification
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
      await api.put('/api/v1/notifications/mark-all-read');
      
      // Update local state
      setNotifications(prev => 
        prev.map(notification => ({ ...notification, read: true }))
      );
      
      // Reset unread count
      setUnreadCount(0);
    } catch (err) {
      console.error('Error marking all notifications as read:', err);
    }
  }, []);

  const deleteNotification = useCallback(async (notificationId: string) => {
    try {
      await api.delete(`/api/v1/notifications/${notificationId}`);
      
      // Update local state
      setNotifications(prev => 
        prev.filter(notification => notification.id !== notificationId)
      );
      
      // Update unread count if the deleted notification was unread
      const deletedNotification = notifications.find(n => n.id === notificationId);
      if (deletedNotification && !deletedNotification.read) {
        setUnreadCount(prev => Math.max(0, prev - 1));
      }
    } catch (err) {
      console.error('Error deleting notification:', err);
    }
  }, [notifications]);

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
      const response = await api.get('/api/v1/notification-preferences');
      setPreferences(response.data);
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
      const response = await api.put('/api/v1/notification-preferences', updates);
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