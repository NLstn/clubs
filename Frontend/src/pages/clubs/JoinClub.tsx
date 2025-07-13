import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import api from '../../utils/api';
import ClubNotFound from './ClubNotFound';
import { removeRecentClub } from '../../utils/recentClubs';

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
  const [isJoining, setIsJoining] = useState(false);

  const fetchClubInfo = useCallback(async () => {
    if (!clubId) {
      setMessage('Invalid club invitation link');
      setLoading(false);
      return;
    }

    try {
      const response = await api.get(`/api/v1/clubs/${clubId}/info`);
      if (response.status === 200) {
        setClub(response.data);
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

    setIsJoining(true);
    setMessage('');

    try {
      const response = await api.post(`/api/v1/clubs/${clubId}/join`);
      if (response.status === 201) {
        setMessage('Join request sent successfully! An admin will review your request.');
        // Redirect to profile invites page after a delay
        setTimeout(() => {
          navigate('/profile/invites');
        }, 3000);
      }
    } catch (error: unknown) {
      console.error('Error joining club:', error);
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
    } finally {
      setIsJoining(false);
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
              <button 
                onClick={() => navigate(`/clubs/${clubId}`)} 
                className="button-accept"
              >
                Go to Club
              </button>
              <button 
                onClick={() => navigate('/')} 
                className="button-cancel"
              >
                Back to Dashboard
              </button>
            </div>
          </div>
        ) : club.hasPendingInvite ? (
          <div>
            <div className="message error">
              You have already been invited to this club! Please check your profile invitations page to accept or decline the invitation.
            </div>
            <div className="join-actions">
              <button 
                onClick={() => navigate('/profile/invites')} 
                className="button-accept"
              >
                View Invitations
              </button>
              <button 
                onClick={() => navigate('/')} 
                className="button-cancel"
              >
                Back to Dashboard
              </button>
            </div>
          </div>
        ) : club.hasPendingRequest ? (
          <div>
            <div className="message error">
              You have already sent a join request for this club. Please wait for an admin to review your request.
            </div>
            <div className="join-actions">
              <button 
                onClick={() => navigate('/')} 
                className="button-cancel"
              >
                Back to Dashboard
              </button>
            </div>
          </div>
        ) : (
          <div>
            <div className="join-actions">
              <button 
                onClick={handleJoinClub} 
                disabled={isJoining}
                className="button-accept"
              >
                {isJoining ? 'Sending Request...' : 'Request to Join'}
              </button>
              <button 
                onClick={() => navigate('/')} 
                className="button-cancel"
              >
                Cancel
              </button>
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