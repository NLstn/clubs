import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getRecentClubs, RecentClub, removeRecentClub } from '../../utils/recentClubs';
import api from '../../utils/api';
import './RecentClubsDropdown.css';

const RecentClubsDropdown: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [recentClubs, setRecentClubs] = useState<RecentClub[]>([]);
  const navigate = useNavigate();
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Load recent clubs on component mount
  useEffect(() => {
    setRecentClubs(getRecentClubs());
  }, []);

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
                <button
                  key={club.id}
                  className="recent-club-item"
                  onClick={() => handleClubClick(club.id)}
                >
                  {club.name}
                </button>
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