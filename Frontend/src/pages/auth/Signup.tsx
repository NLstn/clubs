import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import CookieConsent from '../../components/CookieConsent';
import { Input, Button } from '@/components/ui';
import './Login.css';

const Signup: React.FC = () => {
  const { api } = useAuth();
  const navigate = useNavigate();
  const { user: currentUser } = useCurrentUser();
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setMessage('');

    if (!firstName.trim() || !lastName.trim()) {
      setMessage('Both first name and last name are required');
      setIsSubmitting(false);
      return;
    }

    if (!currentUser?.ID) {
      setMessage('User information not available. Please try again.');
      setIsSubmitting(false);
      return;
    }

    try {
      // OData v2: PATCH to update user entity (PascalCase field names required)
      await api.patch(`/api/v2/Users('${currentUser.ID}')`, {
        FirstName: firstName.trim(),
        LastName: lastName.trim()
      });

      // Navigate to dashboard after successful signup
      navigate('/');
    } catch (error) {
      console.error('Error updating profile:', error);
      setMessage('Failed to update profile. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>Complete Your Profile</h1>
        <p>Please provide your first and last name to complete your account setup.</p>

        {message && (
          <div className={`message ${message.includes('Failed') || message.includes('required') ? 'error' : 'success'}`}>
            {message}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <Input
            label="First Name *"
            type="text"
            id="firstName"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            required
            disabled={isSubmitting}
            placeholder="Enter your first name"
          />

          <Input
            label="Last Name *"
            type="text"
            id="lastName"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            required
            disabled={isSubmitting}
            placeholder="Enter your last name"
          />

          <Button type="submit" disabled={isSubmitting || !firstName.trim() || !lastName.trim()} variant="primary" fullWidth>
            {isSubmitting ? 'Saving...' : 'Complete Profile'}
          </Button>
        </form>
      </div>
      <CookieConsent />
    </div>
  );
};

export default Signup;