import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import './Header.css';

interface HeaderProps {
  title: string;
  showBackButton?: boolean;
}

const Header: React.FC<HeaderProps> = ({ title, showBackButton = true }) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const { logout } = useAuth();
  const navigate = useNavigate();
  const dropdownRef = useRef<HTMLDivElement>(null);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const handleBack = () => {
    navigate(-1);
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
      {showBackButton && (
        <button className="backButton" onClick={handleBack}>
          ← Back
        </button>
      )}
      <h1>{title}</h1>
      <div className="userSection" ref={dropdownRef}>
        <div 
          className="userIcon" 
          onClick={() => setIsDropdownOpen(!isDropdownOpen)}
        >
          {'U'}
        </div>
        
        {isDropdownOpen && (
          <div className="dropdown">
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
