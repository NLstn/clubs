import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from "../../hooks/useAuth";
import { useCurrentUser } from "../../hooks/useCurrentUser";
import { Input, Button, FormGroup } from '@/components/ui';
import './Profile.css';

interface UserProfile {
    firstName: string;
    lastName: string;
    email: string;
    birthDate?: string;
}

const Profile = () => {
    const { api } = useAuth();
    const { user: currentUser, loading: userLoading, error: userError } = useCurrentUser();
    const [profile, setProfile] = useState<UserProfile>({
        firstName: '',
        lastName: '',
        email: '',
        birthDate: undefined
    });
    const [isEditing, setIsEditing] = useState(false);
    const [editedFirstName, setEditedFirstName] = useState('');
    const [editedLastName, setEditedLastName] = useState('');
    const [editedBirthDate, setEditedBirthDate] = useState('');
    const [isSaving, setIsSaving] = useState(false);
    const [message, setMessage] = useState('');

    useEffect(() => {
        if (currentUser) {
            const birthDate = currentUser.BirthDate ? currentUser.BirthDate.split('T')[0] : undefined;
            setProfile({
                firstName: currentUser.FirstName || '',
                lastName: currentUser.LastName || '',
                email: currentUser.Email || '',
                birthDate: birthDate
            });
            setEditedFirstName(currentUser.FirstName || '');
            setEditedLastName(currentUser.LastName || '');
            setEditedBirthDate(birthDate || '');
        } else if (userError) {
            console.error('Error loading user profile:', userError);
            setProfile({
                firstName: 'Demo',
                lastName: 'User',
                email: 'user@example.com'
            });
            setEditedFirstName('Demo');
            setEditedLastName('User');
            setEditedBirthDate('');
        }
    }, [currentUser, userError]);

    const handleEdit = () => {
        setIsEditing(true);
    };

    const handleSave = async () => {
        if (!currentUser?.ID) {
            setMessage('User ID not found');
            return;
        }

        setIsSaving(true);
        try {
            const updateData: {
                firstName: string;
                lastName: string;
                birthDate?: string;
            } = {
                firstName: editedFirstName,
                lastName: editedLastName
            };
            
            // Add birth date if it's set
            if (editedBirthDate) {
                updateData.birthDate = editedBirthDate + 'T00:00:00Z';
            }
            
            // OData v2: PATCH to update user entity
            await api.patch(`/api/v2/Users('${currentUser.ID}')`, updateData);

            setProfile({
                ...profile,
                firstName: editedFirstName,
                lastName: editedLastName,
                birthDate: editedBirthDate || undefined
            });
            setMessage('Profile updated successfully!');
            setTimeout(() => setMessage(''), 3000);
        } catch (error) {
            console.error('Error updating profile:', error);
            setMessage('Failed to update profile');
        } finally {
            setIsSaving(false);
            setIsEditing(false);
        }
    };

    const handleCancel = () => {
        setEditedFirstName(profile.firstName);
        setEditedLastName(profile.lastName);
        setEditedBirthDate(profile.birthDate || '');
        setIsEditing(false);
    };

    return (
        <Layout title="User Profile">
            <ProfileContentLayout
                title={`${profile.firstName || 'User'} ${profile.lastName || ''}`}
                subtitle={profile.email}
                headerContent={
                    <div className="profile-avatar-section">
                        <div className="profile-avatar-placeholder">
                            <span className="avatar-placeholder-text">
                                {(profile.firstName?.[0] || 'U')}
                            </span>
                        </div>
                    </div>
                }
                actions={
                    !isEditing ? (
                        <Button onClick={handleEdit} variant="primary">Edit Profile</Button>
                    ) : (
                        <>
                            <Button 
                                onClick={handleSave}
                                variant="accept"
                                disabled={!editedFirstName.trim() || !editedLastName.trim() || isSaving}
                            >
                                {isSaving ? 'Saving...' : 'Save Changes'}
                            </Button>
                            <Button onClick={handleCancel} variant="cancel" disabled={isSaving}>Cancel</Button>
                        </>
                    )
                }
            >
                {message && (
                    <div className={message.includes('Failed') ? 'error-message' : 'success-message'}>
                        {message}
                    </div>
                )}

                {userLoading ? (
                    <div style={{ textAlign: 'center', padding: 'var(--space-xl)', color: 'var(--color-text-secondary)' }}>
                        <p>Loading profile...</p>
                    </div>
                ) : (
                    <div className="profile-content-sections">
                        <div className="content-section">
                            <h3>Personal Information</h3>
                            <FormGroup>
                                <label htmlFor="firstName">First Name</label>
                                {isEditing ? (
                                    <Input
                                        id="firstName"
                                        type="text"
                                        value={editedFirstName}
                                        onChange={(e) => setEditedFirstName(e.target.value)}
                                        placeholder="Enter your first name"
                                    />
                                ) : (
                                    <div className="form-field-display">
                                        <span>{profile.firstName || 'Not set'}</span>
                                    </div>
                                )}
                            </FormGroup>

                            <FormGroup>
                                <label htmlFor="lastName">Last Name</label>
                                {isEditing ? (
                                    <Input
                                        id="lastName"
                                        type="text"
                                        value={editedLastName}
                                        onChange={(e) => setEditedLastName(e.target.value)}
                                        placeholder="Enter your last name"
                                    />
                                ) : (
                                    <div className="form-field-display">
                                        <span>{profile.lastName || 'Not set'}</span>
                                    </div>
                                )}
                            </FormGroup>

                            <FormGroup>
                                <label htmlFor="birthDate">Birth Date</label>
                                {isEditing ? (
                                    <Input
                                        id="birthDate"
                                        type="date"
                                        value={editedBirthDate}
                                        onChange={(e) => setEditedBirthDate(e.target.value)}
                                    />
                                ) : (
                                    <div className="form-field-display">
                                        <span>{profile.birthDate ? new Date(profile.birthDate).toLocaleDateString() : 'Not set'}</span>
                                    </div>
                                )}
                            </FormGroup>

                            <FormGroup>
                                <label htmlFor="email">Email Address</label>
                                <div className="form-field-display form-field-readonly">
                                    <span>{profile.email}</span>
                                    <span className="field-note">Cannot be changed</span>
                                </div>
                            </FormGroup>
                        </div>
                    </div>
                )}
            </ProfileContentLayout>
        </Layout>
    );
}

export default Profile;