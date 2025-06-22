import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getRecentClubs, RecentClub } from '../../utils/recentClubs';
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

  const handleClubClick = (clubId: string) => {
    setIsOpen(false);
    navigate(`/clubs/${clubId}`);
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
        <span className="clubs-icon">🏛️</span>
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