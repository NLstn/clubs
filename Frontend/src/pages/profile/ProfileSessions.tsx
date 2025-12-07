import { useState, useEffect, useCallback } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useAuth } from '../../hooks/useAuth';
import { Table, TableColumn, Button } from '@/components/ui';
import './Profile.css';

interface Session {
  id: string;
  userAgent: string;
  ipAddress: string;
  createdAt: string;
  isCurrent: boolean;
}

const ProfileSessions = () => {
  const { api } = useAuth();
  const [sessions, setSessions] = useState<Session[]>([]);
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(true);

  const fetchSessions = useCallback(async () => {
    try {
      setIsLoading(true);
      const refreshToken = localStorage.getItem('refresh_token');
      const headers: Record<string, string> = {};
      if (refreshToken) {
        headers['X-Refresh-Token'] = refreshToken;
      }
      const response = await api.get('/api/v1/me/sessions', { headers });
      if (response.status === 200) {
        setSessions(response.data || []);
      }
    } catch (error) {
      console.error('Error fetching sessions:', error);
      setMessage('Failed to load sessions');
    } finally {
      setIsLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchSessions();
  }, [fetchSessions]);

  const handleDeleteSession = async (sessionId: string) => {
    try {
      await api.delete(`/api/v1/me/sessions/${sessionId}`);
      setSessions(sessions.filter(session => session.id !== sessionId));
      setMessage('Session deleted successfully');
      setTimeout(() => setMessage(''), 3000);
    } catch (error) {
      console.error('Error deleting session:', error);
      setMessage('Failed to delete session');
    }
  };

  const formatUserAgent = (userAgent: string) => {
    // Simple user agent parsing to make it more readable
    if (userAgent.includes('Chrome')) return 'Chrome';
    if (userAgent.includes('Firefox')) return 'Firefox';
    if (userAgent.includes('Safari')) return 'Safari';
    if (userAgent.includes('Edge')) return 'Edge';
    return userAgent.length > 50 ? userAgent.substring(0, 50) + '...' : userAgent;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const tableColumns: TableColumn<Session>[] = [
    {
      key: 'browser',
      header: 'Browser',
      render: (session) => formatUserAgent(session.userAgent)
    },
    {
      key: 'ipAddress',
      header: 'IP Address',
      render: (session) => session.ipAddress
    },
    {
      key: 'created',
      header: 'Created',
      render: (session) => formatDate(session.createdAt)
    },
    {
      key: 'status',
      header: 'Status',
      render: (session) => (
        session.isCurrent ? (
          <span style={{ 
            color: 'var(--color-success-text)',
            fontWeight: 'bold'
          }}>
            Current Session
          </span>
        ) : (
          <span style={{ color: 'var(--color-text-secondary)' }}>
            Active
          </span>
        )
      )
    },
    {
      key: 'actions',
      header: 'Actions',
      render: (session) => (
        session.isCurrent ? (
          <span style={{ color: 'var(--color-text-secondary)' }}>
            Cannot delete current session
          </span>
        ) : (
          <Button
            variant="cancel"
            size="sm"
            onClick={() => handleDeleteSession(session.id)}
          >
            Delete
          </Button>
        )
      )
    }
  ];

  return (
    <Layout title="Active Sessions">
      <div style={{
        display: 'flex',
        minHeight: 'calc(100vh - 90px)',
        width: '100%',
        position: 'relative'
      }}>
        <ProfileSidebar />
        <div style={{
          flex: '1 1 auto',
          padding: '20px',
          maxWidth: 'calc(100% - 200px)'
        }}>
          <h2>Active Sessions</h2>
          
          {message && (
            <div className={message.includes('Failed') ? 'error' : 'success'} 
                 style={{ 
                   padding: '10px', 
                   marginBottom: '20px',
                   backgroundColor: message.includes('Failed') ? 'var(--color-error-bg)' : 'var(--color-success-bg)',
                   color: message.includes('Failed') ? 'var(--color-error-text)' : 'var(--color-success-text)',
                   borderRadius: '4px'
                 }}>
              {message}
            </div>
          )}

          <Table
            columns={tableColumns}
            data={sessions}
            keyExtractor={(session) => session.id}
            emptyMessage="No active sessions found."
            loading={isLoading}
            loadingMessage="Loading sessions..."
          />
        </div>
      </div>
    </Layout>
  );
};

export default ProfileSessions;