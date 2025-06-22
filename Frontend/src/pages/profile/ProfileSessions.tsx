import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useAuth } from "../../hooks/useAuth";

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

  useEffect(() => {
    fetchSessions();
  }, []);

  const fetchSessions = async () => {
    try {
      setIsLoading(true);
      const response = await api.get('/api/v1/me/sessions');
      if (response.status === 200) {
        setSessions(response.data || []);
      }
    } catch (error) {
      console.error('Error fetching sessions:', error);
      setMessage('Failed to load sessions');
    } finally {
      setIsLoading(false);
    }
  };

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

          {isLoading ? (
            <p>Loading sessions...</p>
          ) : sessions.length === 0 ? (
            <p>No active sessions found.</p>
          ) : (
            <div style={{ overflowX: 'auto' }}>
              <table style={{ 
                width: '100%', 
                borderCollapse: 'collapse',
                backgroundColor: 'var(--color-card-bg)',
                borderRadius: '8px',
                overflow: 'hidden'
              }}>
                <thead>
                  <tr style={{ backgroundColor: 'var(--color-background-light)' }}>
                    <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid var(--color-border)' }}>
                      Browser
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid var(--color-border)' }}>
                      IP Address
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid var(--color-border)' }}>
                      Created
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid var(--color-border)' }}>
                      Status
                    </th>
                    <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid var(--color-border)' }}>
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {sessions.map((session) => (
                    <tr key={session.id} style={{ borderBottom: '1px solid var(--color-border)' }}>
                      <td style={{ padding: '12px' }}>
                        {formatUserAgent(session.userAgent)}
                      </td>
                      <td style={{ padding: '12px' }}>
                        {session.ipAddress}
                      </td>
                      <td style={{ padding: '12px' }}>
                        {formatDate(session.createdAt)}
                      </td>
                      <td style={{ padding: '12px' }}>
                        {session.isCurrent ? (
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
                        )}
                      </td>
                      <td style={{ padding: '12px' }}>
                        {session.isCurrent ? (
                          <span style={{ color: 'var(--color-text-secondary)' }}>
                            Cannot delete current session
                          </span>
                        ) : (
                          <button
                            className="button-cancel"
                            onClick={() => handleDeleteSession(session.id)}
                            style={{ fontSize: '12px', padding: '6px 12px' }}
                          >
                            Delete
                          </button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default ProfileSessions;