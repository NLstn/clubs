import React from 'react';
import Header from './Header';

interface LayoutProps {
  children: React.ReactNode;
  title: string;
}

const Layout: React.FC<LayoutProps> = ({ children, title }) => {
  return (
    <div>
      <Header title={title} />
      <main style={{ padding: '20px', marginTop: '90px' }}>
        {children}
      </main>
    </div>
  );
};

export default Layout;
