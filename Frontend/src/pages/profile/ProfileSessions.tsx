import { useState, useCallback } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from '../../hooks/useAuth';
import { ODataTable, ODataTableColumn, Button, ConfirmDialog } from '@/components/ui';
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
  const [message, setMessage] = useState('');
  const [refreshKey, setRefreshKey] = useState(0);
  const [sessionToDelete, setSessionToDelete] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const refreshSessions = useCallback(() => {
    setRefreshKey(prev => prev + 1);
  }, []);

  const handleDeleteSession = async () => {
    if (!sessionToDelete) return;

    setIsDeleting(true);
    try {
      // Use OData v2 DELETE with entity key
      await api.delete(`/api/v2/UserSessions('${sessionToDelete}')`);
      setMessage('Session deleted successfully');
      setTimeout(() => setMessage(''), 3000);
      refreshSessions(); // Refresh the list
    } catch (error) {
      console.error('Error deleting session:', error);
      setMessage('Failed to delete session');
      setTimeout(() => setMessage(''), 3000);
    } finally {
      setIsDeleting(false);
      setSessionToDelete(null); // Close modal regardless of success or failure
    }
  };

  const formatUserAgent = (userAgent: string) => {
    // Simple user agent parsing to make it more readable
    // Check Edge first since it contains "Chrome" in its user agent
    if (userAgent.includes('Edg/')) return 'Edge';
    if (userAgent.includes('Chrome')) return 'Chrome';
    if (userAgent.includes('Firefox')) return 'Firefox';
    if (userAgent.includes('Safari') && !userAgent.includes('Chrome')) return 'Safari';
    return userAgent.length > 50 ? userAgent.substring(0, 50) + '...' : userAgent;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const tableColumns: ODataTableColumn<Session>[] = [
    {
      key: 'browser',
      header: 'Browser',
      render: (session) => formatUserAgent(session.UserAgent)
    },
    {
      key: 'IPAddress',
      header: 'IP Address',
      render: (session) => session.IPAddress,
      sortable: true,
      sortField: 'IPAddress'
    },
    {
      key: 'CreatedAt',
      header: 'Created',
      render: (session) => formatDate(session.CreatedAt),
      sortable: true,
      sortField: 'CreatedAt'
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
            onClick={() => setSessionToDelete(session.ID)}
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

        <ODataTable
          key={refreshKey}
          endpoint="/api/v2/UserSessions"
          columns={tableColumns}
          keyExtractor={(session) => session.ID}
          pageSize={10}
          initialSortField="CreatedAt"
          initialSortDirection="desc"
          emptyMessage="No active sessions found."
          loadingMessage="Loading sessions..."
        />

        <ConfirmDialog
          isOpen={sessionToDelete !== null}
          onClose={() => setSessionToDelete(null)}
          onConfirm={handleDeleteSession}
          title="Delete Session"
          message="Are you sure you want to delete this session? You will be logged out from that device."
          confirmText="Delete"
          cancelText="Cancel"
          variant="danger"
          isLoading={isDeleting}
        />
      </ProfileContentLayout>
    </Layout>
  );
};

export default ProfileSessions;