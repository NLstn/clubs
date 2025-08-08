import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { Table, TableColumn } from '@/components/ui';
import api from '../../utils/api';

interface Invitation {
  id: string;
  clubName: string;
}

const ProfileInvites = () => {
  const [invites, setInvites] = useState<Invitation[]>([]);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchInvitations();
  }, []);

  const fetchInvitations = async () => {
    setLoading(true);
    try {
      const response = await api.get('/api/v1/invites');
      if (response.status === 200) {
  const data = response.data;
  // Normalize to array if backend returns null/undefined
  setInvites(Array.isArray(data) ? data : []);
      }
    } catch (error) {
      console.error('Error fetching invitations:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAccept = async (inviteId: string, clubName: string) => {
    try {
      await api.post(`/api/v1/invites/${inviteId}/accept`);
      setInvites((prev) => (Array.isArray(prev) ? prev.filter(invite => invite.id !== inviteId) : []));

      // Show success message
      setMessage(`You've joined ${clubName}!`);
      setTimeout(() => setMessage(''), 3000);
    } catch (error) {
      console.error('Error accepting invitation:', error);
      setMessage('Failed to accept invitation');
    }
  };

  const handleDecline = async (inviteId: string) => {
    try {
      await api.post(`/api/v1/invites/${inviteId}/reject`);
      setInvites((prev) => (Array.isArray(prev) ? prev.filter(invite => invite.id !== inviteId) : []));

      // Show success message
      setMessage('Invitation declined');
      setTimeout(() => setMessage(''), 3000);
    } catch (error) {
      console.error('Error declining invitation:', error);
      setMessage('Failed to decline invitation');
    }
  };

  const columns: TableColumn<Invitation>[] = [
    {
      key: 'clubName',
      header: 'Club',
      render: (invite) => invite.clubName
    },
    {
      key: 'actions',
      header: 'Actions',
      render: (invite) => (
        <div style={{ display: 'flex', gap: '10px' }}>
          <button
            className="button-accept"
            onClick={() => handleAccept(invite.id, invite.clubName)}
          >
            Accept
          </button>
          <button
            className="button-cancel"
            onClick={() => handleDecline(invite.id)}
          >
            Decline
          </button>
        </div>
      )
    }
  ];

  return (
    <Layout title="Club Invitations">
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
          <h2>Pending Invitations</h2>

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

          <div className="invitations-list" style={{ maxWidth: '800px' }}>
            <Table
              columns={columns}
              data={invites}
              keyExtractor={(invite) => invite.id}
              emptyMessage="You have no pending invitations."
              loading={loading}
              loadingMessage="Loading invitations..."
            />
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default ProfileInvites;