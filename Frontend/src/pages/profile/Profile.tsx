import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useAuth } from "../../hooks/useAuth";
import LanguageSwitcher from "../../components/LanguageSwitcher";

interface UserProfile {
    firstName: string;
    lastName: string;
    email: string;
}

const Profile = () => {
    const { api } = useAuth();
    const [profile, setProfile] = useState<UserProfile>({
        firstName: '',
        lastName: '',
        email: ''
    });
    const [isEditing, setIsEditing] = useState(false);
    const [editedFirstName, setEditedFirstName] = useState('');
    const [editedLastName, setEditedLastName] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [message, setMessage] = useState('');

    useEffect(() => {
        const fetchUserProfile = async () => {
            try {
                const response = await api.get('/api/v1/me');
                setProfile({
                    firstName: response.data.FirstName || '',
                    lastName: response.data.LastName || '',
                    email: response.data.Email || ''
                });
                setEditedFirstName(response.data.FirstName || '');
                setEditedLastName(response.data.LastName || '');
            } catch (error) {
                console.error('Error fetching user profile:', error);
                setProfile({
                    firstName: 'Demo',
                    lastName: 'User',
                    email: 'user@example.com'
                });
                setEditedFirstName('Demo');
                setEditedLastName('User');
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
            await api.put('/api/v1/me', {
                firstName: editedFirstName,
                lastName: editedLastName
            });

            setProfile({
                ...profile,
                firstName: editedFirstName,
                lastName: editedLastName
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
        setIsEditing(false);
    };

    return (
        <Layout title="User Profile">
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
                    <h2>My Profile</h2>
                    
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
                        <p>Loading profile...</p>
                    ) : (
                        <div className="profile-container" style={{ maxWidth: '600px' }}>
                            <div className="form-group">
                                <label>First Name:</label>
                                {isEditing ? (
                                    <input
                                        type="text"
                                        value={editedFirstName}
                                        onChange={(e) => setEditedFirstName(e.target.value)}
                                        style={{ width: '100%', padding: '8px', marginBottom: '10px' }}
                                    />
                                ) : (
                                    <div style={{ 
                                        padding: '8px',
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        border: '1px solid transparent',
                                        marginBottom: '10px'
                                    }}>
                                        <span>{profile.firstName}</span>
                                        <button onClick={handleEdit} style={{ padding: '5px 10px' }}>Edit</button>
                                    </div>
                                )}
                            </div>
                            
                            <div className="form-group">
                                <label>Last Name:</label>
                                {isEditing ? (
                                    <input
                                        type="text"
                                        value={editedLastName}
                                        onChange={(e) => setEditedLastName(e.target.value)}
                                        style={{ width: '100%', padding: '8px' }}
                                    />
                                ) : (
                                    <div style={{ 
                                        padding: '8px',
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        border: '1px solid transparent'
                                    }}>
                                        <span>{profile.lastName}</span>
                                        {!isEditing && <div style={{ width: '80px' }}></div>}
                                    </div>
                                )}
                            </div>
                            
                            <div className="form-group">
                                <label>Email:</label>
                                <div style={{ padding: '8px', border: '1px solid transparent' }}>
                                    {profile.email} <span style={{ color: 'var(--color-text-secondary)', fontSize: '0.9em' }}>(cannot be changed)</span>
                                </div>
                            </div>

                            <div className="form-group">
                                <LanguageSwitcher />
                            </div>

                            {isEditing && (
                                <div className="form-actions" style={{ marginTop: '20px' }}>
                                    <button 
                                        onClick={handleSave} 
                                        className="button-accept"
                                        disabled={!editedFirstName.trim() || !editedLastName.trim()}
                                    >
                                        Save
                                    </button>
                                    <button onClick={handleCancel} className="button-cancel">Cancel</button>
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