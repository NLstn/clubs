import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useT } from '../../hooks/useTranslation';
import { useNotifications } from '../../hooks/useNotifications';
import logo from '../../assets/logo.png';
import RecentClubsDropdown from './RecentClubsDropdown';
import NotificationDropdown from './NotificationDropdown';
import './Header.css';

interface HeaderProps {
  title?: string;
  isClubAdmin?: boolean;
  clubId?: string;
  showRecentClubs?: boolean;
}

const Header: React.FC<HeaderProps> = ({ title, isClubAdmin, clubId, showRecentClubs = false }) => {
  const { t } = useT();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const { logout } = useAuth();
  const navigate = useNavigate();
  const dropdownRef = useRef<HTMLDivElement>(null);
  
  // Use notification hooks
  const {
    notifications,
    unreadCount,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    refreshNotifications,
  } = useNotifications();

  const handleLogout = async () => {
    logout();
    navigate('/login');
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <header className="header">
      <img 
        src={logo} 
        alt="Logo" 
        className="headerLogo" 
        onClick={() => navigate('/')}
        style={{ cursor: 'pointer', height: '40px' }} 
      />
      <h1>{title || t('navigation.clubs')}</h1>
      <div className="header-actions">
        {showRecentClubs && <RecentClubsDropdown />}
        <NotificationDropdown
          notifications={notifications}
          unreadCount={unreadCount}
          onMarkAsRead={markAsRead}
          onMarkAllAsRead={markAllAsRead}
          onRefresh={refreshNotifications}
          onDeleteNotification={deleteNotification}
        />
        <div className="userSection" ref={dropdownRef}>
          <div 
            className="userIcon" 
            onClick={() => setIsDropdownOpen(!isDropdownOpen)}
          >
            {'U'}
          </div>
          
          {isDropdownOpen && (
            <div className="dropdown">
              {isClubAdmin && clubId && (
                <button
                  className="dropdownItem"
                  onClick={() => navigate(`/clubs/${clubId}/admin`)}
                >
                  {t('navigation.adminPanel')}
                </button>
              )}
              {!showRecentClubs && (
                <button
                  className="dropdownItem"
                  onClick={() => navigate('/clubs')}
                >
                  {t('navigation.myClubs')}
                </button>
              )}
              <button
                className="dropdownItem"
                onClick={() => navigate('/profile')}
              >
                {t('navigation.profile')}
              </button>
              <button 
                className="dropdownItem" 
                onClick={() => navigate('/createClub')}
              >
                {t('navigation.createNewClub')}
              </button>
              <button 
                className="dropdownItem" 
                onClick={handleLogout}
              >
                {t('navigation.logout')}
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;
