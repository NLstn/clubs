import React from 'react';
import Header from './Header';
import CookieConsent from '../CookieConsent';

interface LayoutProps {
  children: React.ReactNode;
  title?: string;
  showRecentClubs?: boolean;
}

const Layout: React.FC<LayoutProps> = ({ children, title, showRecentClubs = false }) => {
  return (
    <div className="layout">
      <Header title={title} showRecentClubs={showRecentClubs} />
      <main className="main-content">{children}</main>
      <CookieConsent />
    </div>
  );
};

export default Layout;
