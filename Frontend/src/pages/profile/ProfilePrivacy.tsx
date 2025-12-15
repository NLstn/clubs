import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from "../../hooks/useAuth";
import { useCurrentUser } from "../../hooks/useCurrentUser";
import { FormGroup, Card } from '@/components/ui';
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
    const { user: currentUser } = useCurrentUser();
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
                // Wait for current user to be loaded
                if (!currentUser?.ID) {
                    return;
                }

                setIsLoading(true);
                
                // OData v2: Fetch privacy settings
                const privacyResponse = await api.get('/api/v2/UserPrivacySettings');
                const privacyData = privacyResponse.data.value || [];
                
                // Process privacy settings into expected format
                interface ODataPrivacySetting { ID: string; ClubID?: string; ShareBirthDate: boolean; }
                const globalSetting = privacyData.find((s: ODataPrivacySetting) => !s.ClubID);
                const clubSettings = privacyData
                    .filter((s: ODataPrivacySetting) => s.ClubID)
                    .map((s: ODataPrivacySetting) => ({
                        clubId: s.ClubID!,
                        shareBirthDate: s.ShareBirthDate
                    }));
                
                setPrivacySettings({
                    global: { shareBirthDate: globalSetting?.ShareBirthDate || false },
                    clubs: clubSettings
                });
                
                // OData v2: Fetch user's clubs via Members navigation
                const userResponse = await api.get(`/api/v2/Users('${currentUser.ID}')?$expand=Members($expand=Club)`);
                const members = userResponse.data.Members || [];
                
                // Extract clubs from members
                interface ODataMember { Club?: { ID: string; Name: string; }; }
                const mappedClubs = members
                    .filter((m: ODataMember) => m.Club) // Only include members with a club
                    .map((m: ODataMember) => ({
                        id: m.Club!.ID,
                        name: m.Club!.Name
                    }));
                setClubs(mappedClubs);
                
            } catch (error) {
                console.error('Error fetching data:', error);
                setMessage('Failed to load privacy settings');
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, [api, currentUser]);

    const updateGlobalSetting = async (shareBirthDate: boolean) => {
        try {
            // OData v2: Find or create global privacy setting
            const privacyResponse = await api.get('/api/v2/UserPrivacySettings?$filter=ClubID eq null');
            const privacyData = privacyResponse.data.value || [];
            
            if (privacyData.length > 0) {
                // Update existing global setting
                await api.patch(`/api/v2/UserPrivacySettings('${privacyData[0].ID}')`, { ShareBirthDate: shareBirthDate });
            } else {
                // Create new global setting
                await api.post('/api/v2/UserPrivacySettings', { ShareBirthDate: shareBirthDate });
            }
            
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
            // OData v2: Find or create club-specific privacy setting
            const privacyResponse = await api.get(`/api/v2/UserPrivacySettings?$filter=ClubID eq '${clubId}'`);
            const privacyData = privacyResponse.data.value || [];
            
            if (privacyData.length > 0) {
                // Update existing club setting
                await api.patch(`/api/v2/UserPrivacySettings('${privacyData[0].ID}')`, { ShareBirthDate: shareBirthDate });
            } else {
                // Create new club setting
                await api.post('/api/v2/UserPrivacySettings', { ClubID: clubId, ShareBirthDate: shareBirthDate });
            }
            
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
            <ProfileContentLayout title="Privacy Settings">
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
                                            <Card
                                                key={club.id}
                                                variant="default"
                                                padding="sm"
                                                style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
                                            >
                                                <span style={{ fontWeight: '500' }}>{club.name}</span>
                                                <label style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-xs)' }}>
                                                    <input
                                                        type="checkbox"
                                                        checked={getClubSetting(club.id)}
                                                        onChange={(e) => updateClubSetting(club.id, e.target.checked)}
                                                    />
                                                    <span style={{ fontSize: '0.9rem' }}>Share birth date</span>
                                                </label>
                                            </Card>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <div className="profile-section">
                                <Card variant="light" padding="md">
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
                                </Card>
                            </div>
                        </div>
                    )}
            </ProfileContentLayout>
        </Layout>
    );
};

export default ProfilePrivacy;
