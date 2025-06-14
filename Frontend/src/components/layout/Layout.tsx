import React from 'react';
import Header from './Header';
import CookieConsent from '../CookieConsent';

interface LayoutProps {
  children: React.ReactNode;
  title?: string;
}

const Layout: React.FC<LayoutProps> = ({ children, title }) => {
  return (
    <div className="layout">
      <Header title={title} />
      <main className="main-content">{children}</main>
      <CookieConsent />
    </div>
  );
};

export default Layout;
