import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { Button, ButtonState } from '../../components/ui';
import api from '../../utils/api';
import ClubNotFound from './ClubNotFound';
import { removeRecentClub } from '../../utils/recentClubs';
import './JoinClub.css';

interface Club {
  id: string;
  name: string;
  description: string;
  isMember?: boolean;
  hasPendingRequest?: boolean;
  hasPendingInvite?: boolean;
}

const JoinClub: React.FC = () => {
  const { clubId } = useParams<{ clubId: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [club, setClub] = useState<Club | null>(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');
  const [joinButtonState, setJoinButtonState] = useState<ButtonState>('idle');

  const fetchClubInfo = useCallback(async () => {
    if (!clubId) {
      setMessage('Invalid club invitation link');
      setLoading(false);
      return;
    }

    try {
      // OData v2: Fetch club with basic info (note: /info might need custom implementation)
      // For now, fetch club directly - backend should filter appropriately
      const response = await api.get(`/api/v2/Clubs('${clubId}')`);
      if (response.status === 200) {
        const clubData = response.data;
        setClub({
          id: clubData.ID,
          name: clubData.Name,
          description: clubData.Description,
          isMember: clubData.isMember,
          hasPendingRequest: clubData.hasPendingRequest,
          hasPendingInvite: clubData.hasPendingInvite
        });
      }
    } catch (error) {
      console.error('Error fetching club info:', error);
      setMessage('Club not found or invitation link is invalid');
      // Remove from recent clubs if it doesn't exist
      if (clubId) {
        removeRecentClub(clubId);
      }
    } finally {
      setLoading(false);
    }
  }, [clubId]);

  useEffect(() => {
    if (!isAuthenticated) {
      // Store current URL for redirect after login
      localStorage.setItem('loginRedirect', window.location.pathname);
      // Redirect to login with return URL
      navigate(`/login?redirect=${encodeURIComponent(window.location.pathname)}`);
      return;
    }

    fetchClubInfo();
  }, [clubId, isAuthenticated, navigate, fetchClubInfo]);

  const handleJoinClub = async () => {
    if (!clubId) return;

    setJoinButtonState('loading');
    setMessage('');

    try {
      // OData v2: Use Join action on Club entity
      const response = await api.post(`/api/v2/Clubs('${clubId}')/Join`);
      if (response.status === 201 || response.status === 200) {
        setJoinButtonState('success');
        setMessage('Join request sent successfully! An admin will review your request.');
        // Redirect to profile invites page after a delay
        setTimeout(() => {
          navigate('/profile/invites');
        }, 2000);
      }
    } catch (error: unknown) {
      console.error('Error joining club:', error);
      setJoinButtonState('error');
      
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { status?: number; data?: string } };
        if (axiosError.response?.status === 409) {
          // Handle different 409 conflict scenarios based on error message
          const errorMessage = axiosError.response?.data || 'Conflict occurred';
          if (errorMessage.includes('already a member')) {
            setMessage('You are already a member of this club.');
          } else if (errorMessage.includes('pending join request')) {
            setMessage('You already have a pending join request for this club.');
          } else if (errorMessage.includes('pending invitation')) {
            setMessage('You already have a pending invitation for this club. Please check your profile invitations page.');
          } else {
            setMessage('You are already a member of this club.');
          }
        } else if (axiosError.response?.status === 401) {
          setMessage('Please log in to join this club.');
          navigate('/login');
        } else {
          setMessage('Failed to join club. Please try again.');
        }
      } else {
        setMessage('Failed to join club. Please try again.');
      }
      
      setTimeout(() => setJoinButtonState('idle'), 3000);
    }
  };

  if (loading) {
    return (
      <div className="join-club-container">
        <div className="join-club-box">
          <h1>Loading...</h1>
          <p>Please wait while we load the club information.</p>
        </div>
      </div>
    );
  }

  if (!club) {
    return <ClubNotFound clubId={clubId} title="Club Invitation Not Found" message="The club invitation link is invalid or the club no longer exists." />;
  }

  return (
    <div className="join-club-container">
      <div className="join-club-box">
        <h1>Join Club</h1>
        <div className="club-info">
          <h2>{club.name}</h2>
          <p>{club.description}</p>
        </div>

        {message && (
          <div className={`message ${message.includes('successfully') ? 'success' : 'error'}`}>
            {message}
          </div>
        )}

        {club.isMember ? (
          <div>
            <div className="message success">
              You are already a member of this club!
            </div>
            <div className="join-actions">
              <Button 
                onClick={() => navigate(`/clubs/${clubId}`)} 
                variant="accept"
              >
                Go to Club
              </Button>
              <Button 
                onClick={() => navigate('/')} 
                variant="cancel"
              >
                Back to Dashboard
              </Button>
            </div>
          </div>
        ) : club.hasPendingInvite ? (
          <div>
            <div className="message error">
              You have already been invited to this club! Please check your profile invitations page to accept or decline the invitation.
            </div>
            <div className="join-actions">
              <Button 
                onClick={() => navigate('/profile/invites')} 
                variant="accept"
              >
                View Invitations
              </Button>
              <Button 
                onClick={() => navigate('/')} 
                variant="cancel"
              >
                Back to Dashboard
              </Button>
            </div>
          </div>
        ) : club.hasPendingRequest ? (
          <div>
            <div className="message error">
              You have already sent a join request for this club. Please wait for an admin to review your request.
            </div>
            <div className="join-actions">
              <Button 
                onClick={() => navigate('/')} 
                variant="cancel"
              >
                Back to Dashboard
              </Button>
            </div>
          </div>
        ) : (
          <div>
            <div className="join-actions">
              <Button 
                onClick={handleJoinClub} 
                variant="accept"
                state={joinButtonState}
                successMessage="Request sent!"
                errorMessage="Failed to join"
              >
                Request to Join
              </Button>
              <Button 
                onClick={() => navigate('/')} 
                variant="cancel"
                disabled={joinButtonState === 'loading'}
              >
                Cancel
              </Button>
            </div>

            <div className="join-info">
              <p><strong>Note:</strong> Your join request will be sent to the club administrators for approval.</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default JoinClub;