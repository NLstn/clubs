import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { addRecentClub } from '../../utils/recentClubs';
import './NotificationDropdown.css';

interface Notification {
  ID: string;
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

interface NotificationDropdownProps {
  notifications: Notification[];
  unreadCount: number;
  onMarkAsRead: (notificationId: string) => void;
  onMarkAllAsRead: () => void;
  onRefresh: () => void;
  onDeleteNotification: (notificationId: string) => void;
}

const NotificationDropdown: React.FC<NotificationDropdownProps> = ({
  notifications,
  unreadCount,
  onMarkAsRead,
  onMarkAllAsRead,
  onRefresh,
  onDeleteNotification,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleToggle = () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      onRefresh(); // Refresh notifications when opening
    }
  };

  const handleNotificationClick = (notification: Notification) => {
    if (!notification.Read) {
      onMarkAsRead(notification.ID);
    }

    // Navigate to appropriate page based on notification type
    if (notification.Type === 'invite_received') {
      setIsOpen(false); // Close the dropdown
      navigate('/settings/invites');
    } else if (notification.Type === 'member_added' && notification.ClubID) {
      setIsOpen(false); // Close the dropdown
      // Add to recent clubs when navigating from notification
      addRecentClub(notification.ClubID, notification.Title.replace('Welcome to ', ''));
      navigate(`/clubs/${notification.ClubID}`);
    } else if (notification.Type === 'fine_assigned' && notification.ClubID) {
      setIsOpen(false); // Close the dropdown
      // Extract club name from notification title/message if available
      addRecentClub(notification.ClubID, notification.Title.split(' - ')[0] || 'Club');
      navigate(`/clubs/${notification.ClubID}`);
    } else if (notification.Type === 'event_created' && notification.ClubID) {
      setIsOpen(false); // Close the dropdown
      // Extract club name from notification message if available
      addRecentClub(notification.ClubID, notification.Title.split(' - ')[0] || 'Club');
      navigate(`/clubs/${notification.ClubID}`);
    } else if (notification.Type === 'news_created' && notification.ClubID) {
      setIsOpen(false); // Close the dropdown
      // Extract club name from notification message if available
      addRecentClub(notification.ClubID, notification.Title.split(' - ')[0] || 'Club');
      navigate(`/clubs/${notification.ClubID}`);
    } else if (notification.Type === 'join_request_received' && notification.ClubID) {
      setIsOpen(false); // Close the dropdown
      // Extract club name from notification title/message if available
      const clubName = notification.Message.match(/join (.+)$/)?.[1] || 'Club';
      addRecentClub(notification.ClubID, clubName);
      // Navigate to admin members page with join requests modal open
      navigate(`/clubs/${notification.ClubID}/admin/members?openJoinRequests=true`);
    }
  };

  const handleDeleteNotification = (e: React.MouseEvent, notificationId: string) => {
    e.stopPropagation(); // Prevent triggering the notification click
    onDeleteNotification(notificationId);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60));

    if (diffInHours < 1) {
      return 'Just now';
    } else if (diffInHours < 24) {
      return `${diffInHours}h ago`;
    } else {
      const diffInDays = Math.floor(diffInHours / 24);
      return `${diffInDays}d ago`;
    }
  };

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'member_added':
        return '👋';
      case 'invite_received':
        return '✉️';
      case 'event_created':
        return '📅';
      case 'fine_assigned':
        return '💰';
      case 'news_created':
        return '📰';
      case 'join_request_received':
        return '🙋';
      default:
        return '📢';
    }
  };

  return (
    <div className="notification-dropdown" ref={dropdownRef}>
      <button
        className="notification-trigger"
        onClick={handleToggle}
        aria-label={`Notifications ${unreadCount > 0 ? `(${unreadCount} unread)` : ''}`}
      >
        <span className="notification-icon">🔔</span>
        {unreadCount > 0 && (
          <span className="notification-badge">{unreadCount > 99 ? '99+' : unreadCount}</span>
        )}
      </button>

      {isOpen && (
        <div className="notification-menu">
          <div className="notification-header">
            <h3>Notifications</h3>
            {unreadCount > 0 && (
              <button
                className="mark-all-read-btn"
                onClick={onMarkAllAsRead}
              >
                Mark all read
              </button>
            )}
          </div>

          <div className="notification-list">
            {notifications.length === 0 ? (
              <div className="no-notifications">
                <p>No notifications yet</p>
              </div>
            ) : (
              notifications.map((notification) => (
                <div
                  key={notification.ID}
                  className={`notification-item ${!notification.Read ? 'unread' : ''}`}
                  onClick={() => handleNotificationClick(notification)}
                >
                  <div className="notification-content">
                    <div className="notification-icon-wrapper">
                      <span className="notification-type-icon">
                        {getNotificationIcon(notification.Type)}
                      </span>
                      {!notification.Read && <div className="unread-dot"></div>}
                    </div>
                    <div className="notification-text">
                      <div className="notification-title" title={notification.Title}>{notification.Title}</div>
                      <div className="notification-message" title={notification.Message}>{notification.Message}</div>
                      <div className="notification-time">{formatDate(notification.CreatedAt)}</div>
                    </div>
                    <button
                      className="notification-delete-btn"
                      onClick={(e) => handleDeleteNotification(e, notification.ID)}
                      aria-label="Delete notification"
                      title="Delete notification"
                    >
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2m3 0v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6h14zM10 11v6M14 11v6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                      </svg>
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default NotificationDropdown;