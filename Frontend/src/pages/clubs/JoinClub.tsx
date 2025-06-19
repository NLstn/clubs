import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import api from '../../utils/api';

interface Club {
  id: string;
  name: string;
  description: string;
}

const JoinClub: React.FC = () => {
  const { clubId } = useParams<{ clubId: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [club, setClub] = useState<Club | null>(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');
  const [isJoining, setIsJoining] = useState(false);

  useEffect(() => {
    if (!isAuthenticated) {
      // Store current URL for redirect after login
      localStorage.setItem('loginRedirect', window.location.pathname);
      // Redirect to login with return URL
      navigate(`/login?redirect=${encodeURIComponent(window.location.pathname)}`);
      return;
    }

    fetchClubInfo();
  }, [clubId, isAuthenticated, navigate]);

  const fetchClubInfo = async () => {
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
    } finally {
      setLoading(false);
    }
  };

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
    } catch (error: any) {
      console.error('Error joining club:', error);
      if (error.response?.status === 409) {
        setMessage('You are already a member of this club.');
      } else if (error.response?.status === 401) {
        setMessage('Please log in to join this club.');
        navigate('/login');
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
    return (
      <div className="join-club-container">
        <div className="join-club-box">
          <h1>Club Not Found</h1>
          <p>The club invitation link is invalid or the club no longer exists.</p>
          <button onClick={() => navigate('/')} className="button-cancel">
            Go to Dashboard
          </button>
        </div>
      </div>
    );
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
    </div>
  );
};

export default JoinClub;