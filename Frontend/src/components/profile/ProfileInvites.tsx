import { useState, useEffect } from 'react';
import Layout from "../layout/Layout";
import ProfileSidebar from "./ProfileSidebar";

interface Invitation {
  id: string;
  clubId: string;
  clubName: string;
  email: string;
}

const ProfileInvites = () => {
  const [invites, setInvites] = useState<Invitation[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetchInvitations();
  }, []);

  const fetchInvitations = async () => {
    setIsLoading(true);
    try {
      // In a real implementation, fetch invitations for the current user
      // const response = await api.get('/api/v1/invitations');
      
      // For demo purposes, using mock data
      setTimeout(() => {
        setInvites([
          { 
            id: '1', 
            clubId: '101', 
            clubName: 'Tennis Club', 
            email: 'user@example.com' 
          },
          { 
            id: '2', 
            clubId: '102', 
            clubName: 'Chess Masters', 
            email: 'user@example.com' 
          },
        ]);
        setIsLoading(false);
      }, 500);
    } catch (error) {
      console.error('Error fetching invitations:', error);
      setMessage('Failed to load invitations');
      setIsLoading(false);
    }
  };

  const handleAccept = async (inviteId: string, clubName: string) => {
    try {
      // In a real implementation, call the API to accept the invitation
      // await api.post(`/api/v1/invitations/${inviteId}/accept`);
      
      // Update the UI by removing the accepted invitation
      setInvites(invites.filter(invite => invite.id !== inviteId));
      
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
      // In a real implementation, call the API to decline the invitation
      // await api.post(`/api/v1/invitations/${inviteId}/decline`);
      
      // Update the UI by removing the declined invitation
      setInvites(invites.filter(invite => invite.id !== inviteId));
      
      // Show success message
      setMessage('Invitation declined');
      setTimeout(() => setMessage(''), 3000);
    } catch (error) {
      console.error('Error declining invitation:', error);
      setMessage('Failed to decline invitation');
    }
  };

  return (
    <Layout title="Club Invitations" showBackButton={true}>
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
          
          {isLoading ? (
            <p>Loading invitations...</p>
          ) : invites.length === 0 ? (
            <p>You have no pending invitations.</p>
          ) : (
            <div className="invitations-list" style={{ maxWidth: '800px' }}>
              <table>
                <thead>
                  <tr>
                    <th>Club</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {invites.map((invite) => (
                    <tr key={invite.id}>
                      <td>{invite.clubName}</td>
                      <td style={{ display: 'flex', gap: '10px' }}>
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

export default ProfileInvites;