import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './Profile.css';

interface ProfileNavItem {
  label: string;
  path: string;
}

const navItems: ProfileNavItem[] = [
  { label: 'Profile', path: '/profile' },
  { label: 'Preferences', path: '/profile/preferences' },
  { label: 'Privacy', path: '/profile/privacy' },
  { label: 'Invites', path: '/profile/invites' },
  { label: 'Fines', path: '/profile/fines' },
  { label: 'Shifts', path: '/profile/shifts' },
  { label: 'Sessions', path: '/profile/sessions' },
  { label: 'API Keys', path: '/profile/api-keys' },
  { label: 'Notifications', path: '/profile/notifications' }
];

const ProfileSidebar = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [activeItem, setActiveItem] = useState<string>(location.pathname);

  const handleNavigation = (path: string) => {
    setActiveItem(path);
    navigate(path);
  };

  return (
    <div className="profile-sidebar">
      <h3>Navigation</h3>
      <ul className="profile-nav">
        {navItems.map((item) => (
          <li 
            key={item.path}
            onClick={() => handleNavigation(item.path)}
            className={`profile-nav-item ${activeItem === item.path ? 'active' : ''}`}
          >
            {item.label}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ProfileSidebar;