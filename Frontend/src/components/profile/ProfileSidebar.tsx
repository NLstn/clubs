import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

interface ProfileNavItem {
  label: string;
  path: string;
}

const navItems: ProfileNavItem[] = [
  { label: 'Profile', path: '/profile' },
  { label: 'Invites', path: '/profile/invites' }
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
    <div className="profile-sidebar" style={{
      width: '200px',
      flexShrink: 0,
      padding: '20px 0',
      borderRight: '1px solid var(--color-border)',
      height: '100%',
      marginLeft: 0
    }}>
      <h3 style={{ padding: '0 20px', marginTop: 0 }}>Navigation</h3>
      <ul style={{
        listStyle: 'none',
        padding: 0,
        margin: 0
      }}>
        {navItems.map((item) => (
          <li 
            key={item.path}
            onClick={() => handleNavigation(item.path)}
            style={{
              padding: '12px 20px',
              cursor: 'pointer',
              backgroundColor: activeItem === item.path ? 'var(--color-background-light)' : 'transparent',
              borderLeft: activeItem === item.path ? '4px solid var(--color-primary)' : '4px solid transparent',
              transition: 'all 0.2s'
            }}
          >
            {item.label}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ProfileSidebar;