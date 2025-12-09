import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from "../../hooks/useAuth";
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
    const [isLoading, setIsLoading] = useState(true);
    const [message, setMessage] = useState('');

    useEffect(() => {
        const fetchUserProfile = async () => {
            try {
                // OData v2: Query Users entity - backend filters to current user
                const response = await api.get('/api/v2/Users');
                const users = response.data.value || [];
                const userData = users[0] || {};
                const birthDate = userData.BirthDate ? userData.BirthDate.split('T')[0] : undefined;
                setProfile({
                    firstName: userData.FirstName || '',
                    lastName: userData.LastName || '',
                    email: userData.Email || '',
                    birthDate: birthDate
                });
                setEditedFirstName(userData.FirstName || '');
                setEditedLastName(userData.LastName || '');
                setEditedBirthDate(birthDate || '');
            } catch (error) {
                console.error('Error fetching user profile:', error);
                setProfile({
                    firstName: 'Demo',
                    lastName: 'User',
                    email: 'user@example.com'
                });
                setEditedFirstName('Demo');
                setEditedLastName('User');
                setEditedBirthDate('');
            } finally {
                setIsLoading(false);
            }
        };

        fetchUserProfile();
    }, [api]);

    const handleEdit = () => {
        setIsEditing(true);
    };

    const handleSave = async () => {
        setIsLoading(true);
        try {
            // First fetch current user to get ID
            const userResponse = await api.get('/api/v2/Users');
            const users = userResponse.data.value || [];
            const userId = users[0]?.ID;
            
            if (!userId) {
                throw new Error('User ID not found');
            }
            
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
            await api.patch(`/api/v2/Users('${userId}')`, updateData);

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
            setIsLoading(false);
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
                                disabled={!editedFirstName.trim() || !editedLastName.trim()}
                            >
                                Save Changes
                            </Button>
                            <Button onClick={handleCancel} variant="cancel">Cancel</Button>
                        </>
                    )
                }
            >
                {message && (
                    <div className={message.includes('Failed') ? 'error-message' : 'success-message'}>
                        {message}
                    </div>
                )}

                {isLoading ? (
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