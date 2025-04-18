import React from 'react';
import Header from './Header';

interface LayoutProps {
  children: React.ReactNode;
  title: string;
  showBackButton?: boolean;
}

const Layout: React.FC<LayoutProps> = ({ children, title, showBackButton = true }) => {
  return (
    <div>
      <Header title={title} showBackButton={showBackButton} />
      <main style={{ padding: '20px', marginTop: '90px' }}>
        {children}
      </main>
    </div>
  );
};

export default Layout;
