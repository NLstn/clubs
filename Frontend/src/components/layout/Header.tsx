import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import logo from '../../assets/logo.png';
import './Header.css';

interface HeaderProps {
  title?: string;
  isClubAdmin?: boolean;
  clubId?: string;
}

const Header: React.FC<HeaderProps> = ({ title, isClubAdmin, clubId }) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const { logout } = useAuth();
  const navigate = useNavigate();
  const dropdownRef = useRef<HTMLDivElement>(null);

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
      <h1>{title || 'Clubs'}</h1>
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
                Admin Panel
              </button>
            )}
            <button
              className="dropdownItem"
              onClick={() => navigate('/clubs')}
            >
              My Clubs
            </button>
            <button
              className="dropdownItem"
              onClick={() => navigate('/profile')}
            >
              Profile
            </button>
            <button 
              className="dropdownItem" 
              onClick={() => navigate('/createClub')}
            >
              Create New Club
            </button>
            <button 
              className="dropdownItem" 
              onClick={handleLogout}
            >
              Logout
            </button>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;
