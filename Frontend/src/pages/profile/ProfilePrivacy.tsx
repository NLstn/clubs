import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useAuth } from "../../hooks/useAuth";
import { FormGroup } from '@/components/ui';
import './Profile.css';

interface Club {
    id: string;
    name: string;
}

interface PrivacySettings {
    global: {
        shareBirthDate: boolean;
    };
    clubs: Array<{
        clubId: string;
        shareBirthDate: boolean;
    }>;
}

const ProfilePrivacy = () => {
    const { api } = useAuth();
    const [privacySettings, setPrivacySettings] = useState<PrivacySettings>({
        global: { shareBirthDate: false },
        clubs: []
    });
    const [clubs, setClubs] = useState<Club[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [message, setMessage] = useState('');

    useEffect(() => {
        const fetchData = async () => {
            try {
                setIsLoading(true);
                
                // Fetch privacy settings
                const privacyResponse = await api.get('/api/v1/me/privacy/clubs');
                setPrivacySettings(privacyResponse.data);
                
                // Fetch user's clubs
                const clubsResponse = await api.get('/api/v1/clubs');
                setClubs(clubsResponse.data || []);
                
            } catch (error) {
                console.error('Error fetching data:', error);
                setMessage('Failed to load privacy settings');
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, [api]);

    const updateGlobalSetting = async (shareBirthDate: boolean) => {
        try {
            await api.put('/api/v1/me/privacy', { shareBirthDate });
            setPrivacySettings(prev => ({
                ...prev,
                global: { shareBirthDate }
            }));
            setMessage('Global privacy settings updated successfully!');
            setTimeout(() => setMessage(''), 3000);
        } catch (error) {
            console.error('Error updating global privacy settings:', error);
            setMessage('Failed to update global privacy settings');
        }
    };

    const updateClubSetting = async (clubId: string, shareBirthDate: boolean) => {
        try {
            await api.put('/api/v1/me/privacy', { shareBirthDate, clubId });
            
            setPrivacySettings(prev => {
                const updatedClubs = prev.clubs.filter(c => c.clubId !== clubId);
                if (shareBirthDate !== prev.global.shareBirthDate) {
                    updatedClubs.push({ clubId, shareBirthDate });
                }
                return {
                    ...prev,
                    clubs: updatedClubs
                };
            });
            
            setMessage('Club privacy settings updated successfully!');
            setTimeout(() => setMessage(''), 3000);
        } catch (error) {
            console.error('Error updating club privacy settings:', error);
            setMessage('Failed to update club privacy settings');
        }
    };

    const getClubSetting = (clubId: string): boolean => {
        const clubSetting = privacySettings.clubs ? privacySettings.clubs.find(c => c.clubId === clubId) : undefined;
        return clubSetting ? clubSetting.shareBirthDate : privacySettings.global.shareBirthDate;
    };

    return (
        <Layout title="Privacy Settings">
            <div className="profile-layout">
                <ProfileSidebar />
                <div className="profile-content">
                    <div className="profile-header">
                        <h2>Privacy Settings</h2>
                        <p>Control who can see your personal information</p>
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
                            <p>Loading privacy settings...</p>
                        </div>
                    ) : (
                        <div className="profile-container" style={{ maxWidth: '700px' }}>
                            <div className="profile-section">
                                <h3 className="profile-section-title">Global Settings</h3>
                                <p style={{ color: 'var(--color-text-secondary)', marginBottom: 'var(--space-md)' }}>
                                    These settings apply to all clubs unless overridden by club-specific settings below.
                                </p>
                                
                                <FormGroup>
                                    <label style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-sm)' }}>
                                        <input
                                            type="checkbox"
                                            checked={privacySettings.global.shareBirthDate}
                                            onChange={(e) => updateGlobalSetting(e.target.checked)}
                                        />
                                        <span>Share my birth date with club members</span>
                                    </label>
                                </FormGroup>
                            </div>

                            <div className="profile-section">
                                <h3 className="profile-section-title">Club-Specific Settings</h3>
                                <p style={{ color: 'var(--color-text-secondary)', marginBottom: 'var(--space-md)' }}>
                                    Override global settings for specific clubs. If not set, the global setting will be used.
                                </p>
                                
                                {clubs.length === 0 ? (
                                    <p style={{ color: 'var(--color-text-secondary)' }}>
                                        You're not a member of any clubs yet.
                                    </p>
                                ) : (
                                    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-sm)' }}>
                                        {clubs.map((club) => (
                                            <div key={club.id} style={{
                                                display: 'flex',
                                                justifyContent: 'space-between',
                                                alignItems: 'center',
                                                padding: 'var(--space-sm)',
                                                backgroundColor: 'var(--color-card-bg)',
                                                borderRadius: 'var(--border-radius)',
                                                border: '1px solid var(--color-border)'
                                            }}>
                                                <span style={{ fontWeight: '500' }}>{club.name}</span>
                                                <label style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-xs)' }}>
                                                    <input
                                                        type="checkbox"
                                                        checked={getClubSetting(club.id)}
                                                        onChange={(e) => updateClubSetting(club.id, e.target.checked)}
                                                    />
                                                    <span style={{ fontSize: '0.9rem' }}>Share birth date</span>
                                                </label>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <div className="profile-section">
                                <div style={{
                                    padding: 'var(--space-md)',
                                    backgroundColor: 'var(--color-background-light)',
                                    borderRadius: 'var(--border-radius)',
                                    border: '1px solid var(--color-border)'
                                }}>
                                    <h4 style={{ margin: '0 0 var(--space-sm) 0', color: 'var(--color-text-primary)' }}>
                                        Privacy Information
                                    </h4>
                                    <ul style={{ 
                                        margin: 0, 
                                        paddingLeft: 'var(--space-md)',
                                        color: 'var(--color-text-secondary)',
                                        fontSize: '0.9rem'
                                    }}>
                                        <li>Your birth date will only be visible to other members of clubs where you've enabled sharing</li>
                                        <li>Club administrators can always see member information for administrative purposes</li>
                                        <li>You can change these settings at any time</li>
                                        <li>Global settings provide a default for all clubs, but you can override them per club</li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </Layout>
    );
};

export default ProfilePrivacy;
