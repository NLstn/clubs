import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useAuth } from "../../hooks/useAuth";
import LanguageSwitcher from "../../components/LanguageSwitcher";

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
                const response = await api.get('/api/v1/me');
                const birthDate = response.data.BirthDate ? response.data.BirthDate.split('T')[0] : undefined;
                setProfile({
                    firstName: response.data.FirstName || '',
                    lastName: response.data.LastName || '',
                    email: response.data.Email || '',
                    birthDate: birthDate
                });
                setEditedFirstName(response.data.FirstName || '');
                setEditedLastName(response.data.LastName || '');
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
            
            await api.put('/api/v1/me', updateData);

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
            <div className="profile-layout">
                <ProfileSidebar />
                <div className="profile-content">
                    <div className="profile-header">
                        <h2>My Profile</h2>
                    </div>
                    
                    {message && (
                        <div className={message.includes('Failed') ? 'error-message' : 'success-message'}>
                            {message}
                        </div>
                    )}
                    
                    {isLoading ? (
                        <div style={{ 
                            textAlign: 'center', 
                            padding: 'var(--space-xl)',
                            color: 'var(--color-text-secondary)'
                        }}>
                            <p>Loading profile...</p>
                        </div>
                    ) : (
                        <div className="profile-container" style={{ maxWidth: '700px' }}>
                            <div className="profile-section">
                                <div style={{ 
                                    display: 'flex', 
                                    justifyContent: 'space-between', 
                                    alignItems: 'center',
                                    marginBottom: 'var(--space-md)'
                                }}>
                                    <h3 className="profile-section-title" style={{ margin: 0 }}>Personal Information</h3>
                                    {!isEditing && (
                                        <button onClick={handleEdit} className="edit-button">
                                            Edit Profile
                                        </button>
                                    )}
                                </div>
                                
                                <div className="form-group">
                                    <label htmlFor="firstName">First Name</label>
                                    {isEditing ? (
                                        <input
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
                                </div>
                                
                                <div className="form-group">
                                    <label htmlFor="lastName">Last Name</label>
                                    {isEditing ? (
                                        <input
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
                                </div>
                                
                                <div className="form-group">
                                    <label htmlFor="birthDate">Birth Date</label>
                                    {isEditing ? (
                                        <input
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
                                </div>
                                
                                <div className="form-group">
                                    <label htmlFor="email">Email Address</label>
                                    <div className="form-field-display form-field-readonly">
                                        <span>{profile.email}</span>
                                        <span className="field-note">Cannot be changed</span>
                                    </div>
                                </div>
                            </div>

                            <div className="profile-section">
                                <h3 className="profile-section-title">Preferences</h3>
                                <div className="form-group">
                                    <label>Language</label>
                                    <div style={{ marginTop: 'var(--space-xs)' }}>
                                        <LanguageSwitcher />
                                    </div>
                                </div>
                            </div>

                            {isEditing && (
                                <div className="form-actions">
                                    <button 
                                        onClick={handleSave} 
                                        className="button-accept"
                                        disabled={!editedFirstName.trim() || !editedLastName.trim()}
                                    >
                                        Save Changes
                                    </button>
                                    <button onClick={handleCancel} className="button-cancel">
                                        Cancel
                                    </button>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </Layout>
    );
}

export default Profile;