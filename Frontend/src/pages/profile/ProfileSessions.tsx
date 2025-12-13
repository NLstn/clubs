import { useState, useEffect, useCallback } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from '../../hooks/useAuth';
import { Table, TableColumn, Button } from '@/components/ui';
import './Profile.css';

interface Session {
  ID: string;
  UserAgent: string;
  IPAddress: string;
  CreatedAt: string;
  IsCurrent: boolean;
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
      // Use OData v2 API - filter and ordering handled by backend
      const response = await api.get('/api/v2/UserSessions?$orderby=CreatedAt desc', { headers });
      if (response.status === 200) {
        setSessions(response.data?.value || []);
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
      // Use OData v2 DELETE with entity key
      await api.delete(`/api/v2/UserSessions('${sessionId}')`);
      setSessions(sessions.filter(session => session.ID !== sessionId));
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
      render: (session) => formatUserAgent(session.UserAgent)
    },
    {
      key: 'ipAddress',
      header: 'IP Address',
      render: (session) => session.IPAddress
    },
    {
      key: 'created',
      header: 'Created',
      render: (session) => formatDate(session.CreatedAt)
    },
    {
      key: 'status',
      header: 'Status',
      render: (session) => (
        session.IsCurrent ? (
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
        session.IsCurrent ? (
          <span style={{ color: 'var(--color-text-secondary)' }}>
            Cannot delete current session
          </span>
        ) : (
          <Button
            variant="cancel"
            size="sm"
            onClick={() => handleDeleteSession(session.ID)}
          >
            Delete
          </Button>
        )
      )
    }
  ];

  return (
    <Layout title="Active Sessions">
      <ProfileContentLayout title="Active Sessions">
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
          keyExtractor={(session) => session.ID}
          emptyMessage="No active sessions found."
          loading={isLoading}
          loadingMessage="Loading sessions..."
        />
      </ProfileContentLayout>
    </Layout>
  );
};

export default ProfileSessions;