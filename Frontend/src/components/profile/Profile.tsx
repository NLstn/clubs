import { useState, useEffect } from 'react';
import Layout from "../layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useAuth } from "../../context/AuthContext";

interface UserProfile {
    name: string;
    email: string;
}

const Profile = () => {
    const { api } = useAuth();
    const [profile, setProfile] = useState<UserProfile>({
        name: '',
        email: ''
    });
    const [isEditing, setIsEditing] = useState(false);
    const [editedName, setEditedName] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [message, setMessage] = useState('');

    useEffect(() => {
        const fetchUserProfile = async () => {
            try {
                const response = await api.get('/api/v1/auth/me');
                console.log('User profile response:', response.data);
                setProfile({
                    name: response.data.Name || 'User',
                    email: response.data.Email || ''
                });
                setEditedName(response.data.Name || 'User');
            } catch (error) {
                console.error('Error fetching user profile:', error);
                setProfile({
                    name: 'Demo User',
                    email: 'user@example.com'
                });
                setEditedName('Demo User');
            } finally {
                setIsLoading(false);
            }
        };

        fetchUserProfile();
    }, []);

    const handleEdit = () => {
        setIsEditing(true);
    };

    const handleSave = async () => {
        setIsLoading(true);
        try {
            await api.put('/api/v1/auth/me', {
                name: editedName
            });

            setProfile({
                ...profile,
                name: editedName
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
        setEditedName(profile.name);
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
                                <label>Name:</label>
                                {isEditing ? (
                                    <input
                                        type="text"
                                        value={editedName}
                                        onChange={(e) => setEditedName(e.target.value)}
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
                                        <span>{profile.name}</span>
                                        <button onClick={handleEdit} style={{ padding: '5px 10px' }}>Edit</button>
                                    </div>
                                )}
                            </div>
                            
                            <div className="form-group">
                                <label>Email:</label>
                                <div style={{ padding: '8px', border: '1px solid transparent' }}>
                                    {profile.email} <span style={{ color: 'var(--color-text-secondary)', fontSize: '0.9em' }}>(cannot be changed)</span>
                                </div>
                            </div>

                            {isEditing && (
                                <div className="form-actions" style={{ marginTop: '20px' }}>
                                    <button onClick={handleSave} className="button-accept">Save</button>
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