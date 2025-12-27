import { useState, useEffect } from 'react';
import Layout from "../../components/layout/Layout";
import SimpleSettingsLayout from '../../components/layout/SimpleSettingsLayout';
import { useAuth } from "../../hooks/useAuth";
import { useCurrentUser } from "../../hooks/useCurrentUser";
import { FormGroup, Card } from '@/components/ui';
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import './Profile.css';

interface Club {
    id: string;
    name: string;
    memberId: string; // Added to track member ID for privacy settings
}

interface PrivacySettings {
    global: {
        shareBirthDate: boolean;
        id?: string; // Track ID for updates
    };
    clubs: Array<{
        memberId: string;
        shareBirthDate: boolean;
        id?: string; // Track ID for updates/deletes
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
                setIsLoading(true);
                
                // Wait for current user to be loaded
                if (!currentUser?.ID) {
                    setIsLoading(false);
                    return;
                }
                
                // OData v2: Fetch global privacy settings
                interface ODataPrivacySetting { ID: string; ShareBirthDate: boolean; }
                const privacyResponse = await api.get<ODataCollectionResponse<ODataPrivacySetting>>('/api/v2/UserPrivacySettings?$select=ID,ShareBirthDate');
                const privacyData = parseODataCollection(privacyResponse.data);
                const globalSetting = privacyData[0]; // Should only be one global setting per user
                
                // OData v2: Fetch user's members with clubs and privacy settings
                const encodedUserId = encodeURIComponent(currentUser.ID);
                interface ODataMember {
                    ID: string;
                    Club: { ID: string; Name: string; };
                    PrivacySettings?: { ID: string; ShareBirthDate: boolean; };
                }
                const membersResponse = await api.get<ODataCollectionResponse<ODataMember>>(
                    `/api/v2/Users('${encodedUserId}')/Members?$select=ID&$filter=Club/Deleted eq false&$expand=Club($select=ID,Name),PrivacySettings($select=ID,ShareBirthDate)`
                );
                const members = parseODataCollection(membersResponse.data);
                
                // Process members data
                
                const mappedClubs = members.map((m: ODataMember) => ({
                    id: m.Club.ID,
                    name: m.Club.Name,
                    memberId: m.ID
                }));
                setClubs(mappedClubs);
                
                const clubSettings = members
                    .filter((m: ODataMember) => m.PrivacySettings)
                    .map((m: ODataMember) => ({
                        memberId: m.ID,
                        shareBirthDate: m.PrivacySettings!.ShareBirthDate,
                        id: m.PrivacySettings!.ID
                    }));
                
                setPrivacySettings({
                    global: { 
                        shareBirthDate: globalSetting?.ShareBirthDate || false,
                        id: globalSetting?.ID
                    },
                    clubs: clubSettings
                });
                
            } catch (error) {
                console.error('Error fetching data:', error);
                setMessage('Failed to load privacy settings');
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, [api, currentUser?.ID]);

    const updateGlobalSetting = async (shareBirthDate: boolean) => {
        try {
            if (privacySettings.global.id) {
                // Update existing global setting
                await api.patch(`/api/v2/UserPrivacySettings('${privacySettings.global.id}')`, { 
                    ShareBirthDate: shareBirthDate 
                });
                setPrivacySettings(prev => ({
                    ...prev,
                    global: { ...prev.global, shareBirthDate }
                }));
            } else {
                // Create new global setting
                const response = await api.post('/api/v2/UserPrivacySettings', { 
                    ShareBirthDate: shareBirthDate 
                });
                setPrivacySettings(prev => ({
                    ...prev,
                    global: { shareBirthDate, id: response.data.ID }
                }));
            }
            
            setMessage('Global privacy settings updated successfully!');
            setTimeout(() => setMessage(''), 3000);
        } catch (error) {
            console.error('Error updating global privacy settings:', error);
            setMessage('Failed to update global privacy settings');
        }
    };

    const updateClubSetting = async (memberId: string, shareBirthDate: boolean) => {
        try {
            const existingSetting = privacySettings.clubs.find(c => c.memberId === memberId);
            const globalSetting = privacySettings.global.shareBirthDate;
            
            // If new setting matches global setting, delete the override (if it exists)
            if (shareBirthDate === globalSetting) {
                if (existingSetting?.id) {
                    await api.delete(`/api/v2/MemberPrivacySettings('${existingSetting.id}')`);
                    setPrivacySettings(prev => ({
                        ...prev,
                        clubs: prev.clubs.filter(c => c.memberId !== memberId)
                    }));
                }
            } else {
                // Setting differs from global, create or update override
                if (existingSetting?.id) {
                    // Update existing override
                    await api.patch(`/api/v2/MemberPrivacySettings('${existingSetting.id}')`, { 
                        ShareBirthDate: shareBirthDate 
                    });
                    setPrivacySettings(prev => ({
                        ...prev,
                        clubs: prev.clubs.map(c => 
                            c.memberId === memberId 
                                ? { ...c, shareBirthDate } 
                                : c
                        )
                    }));
                } else {
                    // Create new override
                    const response = await api.post('/api/v2/MemberPrivacySettings', { 
                        MemberID: memberId,
                        ShareBirthDate: shareBirthDate 
                    });
                    setPrivacySettings(prev => ({
                        ...prev,
                        clubs: [...prev.clubs, { 
                            memberId, 
                            shareBirthDate, 
                            id: response.data.ID 
                        }]
                    }));
                }
            }
            
            setMessage('Club privacy settings updated successfully!');
            setTimeout(() => setMessage(''), 3000);
        } catch (error) {
            console.error('Error updating club privacy settings:', error);
            setMessage('Failed to update club privacy settings');
        }
    };

    const getClubSetting = (memberId: string): boolean => {
        const clubSetting = privacySettings.clubs.find(c => c.memberId === memberId);
        return clubSetting ? clubSetting.shareBirthDate : privacySettings.global.shareBirthDate;
    };

    return (
        <Layout title="Privacy Settings">
            <SimpleSettingsLayout title="Privacy Settings">
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
                                                        checked={getClubSetting(club.memberId)}
                                                        onChange={(e) => updateClubSetting(club.memberId, e.target.checked)}
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
            </SimpleSettingsLayout>
        </Layout>
    );
};

export default ProfilePrivacy;
