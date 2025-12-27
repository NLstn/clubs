import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './Profile.css';

interface ProfileNavItem {
  label: string;
  path: string;
}

const navItems: ProfileNavItem[] = [
  { label: 'Profile', path: '/settings/profile' },
  { label: 'Preferences', path: '/settings/preferences' },
  { label: 'Privacy', path: '/settings/privacy' },
  { label: 'Invites', path: '/settings/invites' },
  { label: 'Fines', path: '/settings/fines' },
  { label: 'Shifts', path: '/settings/shifts' },
  { label: 'Sessions', path: '/settings/sessions' },
  { label: 'API Keys', path: '/settings/api-keys' },
  { label: 'Notifications', path: '/settings/notifications' }
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
      <ul className="profile-nav" role="tablist">
        {navItems.map((item) => (
          <li 
            key={item.path}
            onClick={() => handleNavigation(item.path)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                handleNavigation(item.path);
              }
            }}
            className={`profile-nav-item ${activeItem === item.path ? 'active' : ''}`}
            tabIndex={0}
            role="tab"
            aria-selected={activeItem === item.path}
          >
            {item.label}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ProfileSidebar;