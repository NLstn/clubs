import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { Table, TableColumn, Button } from '@/components/ui';
import api from '../../utils/api';
import './Profile.css';

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
      // OData v2: Query Invites entity - backend filters to current user's invites
      const response = await api.get('/api/v2/Invites?$expand=Club');
      if (response.status === 200) {
        interface ODataInvite { ID: string; ClubName?: string; Club?: { Name: string; }; }
        const data = response.data.value || [];
        // Map OData response to match expected format
        const mappedInvites = data.map((invite: ODataInvite) => ({
          id: invite.ID,
          clubName: invite.Club?.Name || invite.ClubName || 'Unknown Club'
        }));
        setInvites(mappedInvites);
      }
    } catch (error) {
      console.error('Error fetching invitations:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAccept = async (inviteId: string, clubName: string) => {
    try {
      // OData v2: Use Accept action on Invite entity
      await api.post(`/api/v2/Invites('${inviteId}')/Accept`);
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
      // OData v2: Use Reject action on Invite entity
      await api.post(`/api/v2/Invites('${inviteId}')/Reject`);
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
          <Button
            variant="accept"
            size="sm"
            onClick={() => handleAccept(invite.id, invite.clubName)}
          >
            Accept
          </Button>
          <Button
            variant="cancel"
            size="sm"
            onClick={() => handleDecline(invite.id)}
          >
            Decline
          </Button>
        </div>
      )
    }
  ];

  return (
    <Layout title="Club Invitations">
      <ProfileContentLayout title="Pending Invitations">
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
          columns={columns}
          data={invites}
          keyExtractor={(invite) => invite.id}
          emptyMessage="You have no pending invitations."
          loading={loading}
          loadingMessage="Loading invitations..."
        />
      </ProfileContentLayout>
    </Layout>
  );
};

export default ProfileInvites;