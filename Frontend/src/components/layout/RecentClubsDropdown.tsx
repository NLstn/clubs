import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getRecentClubs, RecentClub, removeRecentClub } from '../../utils/recentClubs';
import api from '../../utils/api';
import './RecentClubsDropdown.css';

const RecentClubsDropdown: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  // Initialize state directly from getRecentClubs()
  const [recentClubs, setRecentClubs] = useState<RecentClub[]>(() => getRecentClubs());
  const navigate = useNavigate();
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleClubClick = async (clubId: string) => {
    setIsOpen(false);
    
    try {
      // First try to check if the club exists
      await api.get(`/api/v1/clubs/${clubId}`);
      navigate(`/clubs/${clubId}`);
    } catch {
      // If club doesn't exist, remove it from recent clubs
      console.warn(`Club ${clubId} not found, removing from recent clubs`);
      removeRecentClub(clubId);
      setRecentClubs(getRecentClubs()); // Refresh the list
      
      // Still navigate to the club page, which will show the ClubNotFound component
      navigate(`/clubs/${clubId}`);
    }
  };

  const handleViewAllClubs = () => {
    setIsOpen(false);
    navigate('/clubs');
  };

  const handleRemoveClub = (e: React.MouseEvent, clubId: string) => {
    e.stopPropagation(); // Prevent triggering the club click
    removeRecentClub(clubId);
    setRecentClubs(getRecentClubs()); // Refresh the list
  };

  return (
    <div className="recent-clubs-dropdown" ref={dropdownRef}>
      <div 
        className="recent-clubs-trigger"
        onClick={() => setIsOpen(!isOpen)}
        title="Recent Clubs"
      >
        <span className="clubs-text">Recent clubs</span>
      </div>

      {isOpen && (
        <div className="recent-clubs-menu">
          <div className="recent-clubs-header">Recent Clubs</div>
          
          {recentClubs.length > 0 ? (
            <>
              {recentClubs.map((club) => (
                <div
                  key={club.id}
                  className="recent-club-item-container"
                >
                  <button
                    className="recent-club-item"
                    onClick={() => handleClubClick(club.id)}
                  >
                    {club.name}
                  </button>
                  <button
                    className="recent-club-remove"
                    onClick={(e) => handleRemoveClub(e, club.id)}
                    title="Remove from recent clubs"
                    aria-label={`Remove ${club.name} from recent clubs`}
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2m3 0v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6h14zM10 11v6M14 11v6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                  </button>
                </div>
              ))}
              <div className="recent-clubs-divider"></div>
            </>
          ) : (
            <div className="no-recent-clubs">No recent clubs</div>
          )}
          
          <button
            className="view-all-clubs"
            onClick={handleViewAllClubs}
          >
            View All Clubs
          </button>
        </div>
      )}
    </div>
  );
};

export default RecentClubsDropdown;