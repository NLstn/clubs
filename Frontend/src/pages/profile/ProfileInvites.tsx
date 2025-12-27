import { useState, useEffect, useRef } from 'react';
import Layout from "../../components/layout/Layout";
import SimpleSettingsLayout from '../../components/layout/SimpleSettingsLayout';
import { Table, TableColumn, Button, ButtonState } from '@/components/ui';
import api from '../../utils/api';
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import './Profile.css';

interface Invitation {
  id: string;
  clubName: string;
}

const ProfileInvites = () => {
  const [invites, setInvites] = useState<Invitation[]>([]);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [processingInvites, setProcessingInvites] = useState<Record<string, ButtonState>>({});
  const timeoutRefs = useRef<number[]>([]);

  useEffect(() => {
    fetchInvitations();
  }, []);

  useEffect(() => {
    // Cleanup timeouts on unmount
    const timeouts = timeoutRefs.current;
    return () => {
      timeouts.forEach(clearTimeout);
    };
  }, []);

  const fetchInvitations = async () => {
    setLoading(true);
    try {
      // OData v2: Query Invites entity - backend filters to current user's invites
      interface ODataInvite { ID: string; ClubName?: string; Club?: { Name: string; }; }
      const response = await api.get<ODataCollectionResponse<ODataInvite>>('/api/v2/Invites?$select=ID&$expand=Club($select=Name)');
      if (response.status === 200) {
        const data = parseODataCollection(response.data);
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
    setProcessingInvites(prev => ({ ...prev, [`accept-${inviteId}`]: 'loading' }));
    
    try {
      // OData v2: Use Accept action on Invite entity
      await api.post(`/api/v2/Invites('${inviteId}')/Accept`);
      setProcessingInvites(prev => ({ ...prev, [`accept-${inviteId}`]: 'success' }));

      // Show success message
      setMessage(`You've joined ${clubName}!`);
      
      const timeoutId = window.setTimeout(() => {
        setInvites((prev) => (Array.isArray(prev) ? prev.filter(invite => invite.id !== inviteId) : []));
        setMessage('');
        setProcessingInvites(prev => {
          const newState = { ...prev };
          delete newState[`accept-${inviteId}`];
          return newState;
        });
      }, 1500);
      timeoutRefs.current.push(timeoutId);
    } catch (error) {
      console.error('Error accepting invitation:', error);
      setProcessingInvites(prev => ({ ...prev, [`accept-${inviteId}`]: 'error' }));
      setMessage('Failed to accept invitation');
      
      const timeoutId = window.setTimeout(() => {
        setProcessingInvites(prev => {
          const newState = { ...prev };
          delete newState[`accept-${inviteId}`];
          return newState;
        });
      }, 3000);
      timeoutRefs.current.push(timeoutId);
    }
  };

  const handleDecline = async (inviteId: string) => {
    setProcessingInvites(prev => ({ ...prev, [`decline-${inviteId}`]: 'loading' }));
    
    try {
      // OData v2: Use Reject action on Invite entity
      await api.post(`/api/v2/Invites('${inviteId}')/Reject`);
      setProcessingInvites(prev => ({ ...prev, [`decline-${inviteId}`]: 'success' }));

      // Show success message
      setMessage('Invitation declined');
      
      const timeoutId = window.setTimeout(() => {
        setInvites((prev) => (Array.isArray(prev) ? prev.filter(invite => invite.id !== inviteId) : []));
        setMessage('');
        setProcessingInvites(prev => {
          const newState = { ...prev };
          delete newState[`decline-${inviteId}`];
          return newState;
        });
      }, 1500);
      timeoutRefs.current.push(timeoutId);
    } catch (error) {
      console.error('Error declining invitation:', error);
      setProcessingInvites(prev => ({ ...prev, [`decline-${inviteId}`]: 'error' }));
      setMessage('Failed to decline invitation');
      
      const timeoutId = window.setTimeout(() => {
        setProcessingInvites(prev => {
          const newState = { ...prev };
          delete newState[`decline-${inviteId}`];
          return newState;
        });
      }, 3000);
      timeoutRefs.current.push(timeoutId);
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
            state={processingInvites[`accept-${invite.id}`] || 'idle'}
            successMessage="Joined!"
            errorMessage="Failed"
          >
            Accept
          </Button>
          <Button
            variant="cancel"
            size="sm"
            onClick={() => handleDecline(invite.id)}
            state={processingInvites[`decline-${invite.id}`] || 'idle'}
            successMessage="Declined"
            errorMessage="Failed"
          >
            Decline
          </Button>
        </div>
      )
    }
  ];

  return (
    <Layout title="Club Invitations">
      <SimpleSettingsLayout title="Pending Invitations">
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
      </SimpleSettingsLayout>
    </Layout>
  );
};

export default ProfileInvites;