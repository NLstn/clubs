import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../../components/layout/Layout';
import { Button, Card } from '../../components/ui';
import api from '../../utils/api';
import { addRecentClub } from '../../utils/recentClubs';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import { useT } from '../../hooks/useTranslation';
import './ClubList.css';

interface Club {
    id: string;
    name: string;
    description: string;
    user_role: string;
    created_at: string;
    deleted?: boolean;
    user_teams?: Team[];
}

interface Team {
    id: string;
    name: string;
    description: string;
    createdAt: string;
    clubId: string;
}

// OData response types for User -> Members -> Club navigation
interface ODataMemberWithClub {
    Role: string;
    Club: ODataClubWithTeams; // Club is guaranteed by $expand query
}

interface ODataClubWithTeams {
    ID: string;
    Name: string;
    Description: string;
    CreatedAt: string;
    Deleted?: boolean;
    Teams?: ODataTeam[];
}

interface ODataTeam {
    ID: string;
    Name: string;
    Description: string;
    CreatedAt: string;
    ClubID: string;
}

const ClubList = () => {
    const { t } = useT();
    const { user: currentUser } = useCurrentUser();
    const [clubs, setClubs] = useState<Club[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    const fetchClubs = useCallback(async () => {
        try {
            // Wait for current user to be loaded
            if (!currentUser?.ID) {
                return;
            }

            // OData v2 API: Get user's clubs via Members navigation
            // This uses the User -> Members -> Club navigation pattern
            const response = await api.get(`/api/v2/Users('${currentUser.ID}')?$expand=Members($expand=Club($expand=Teams))`);
            
            // Extract members with their clubs from the response
            const members: ODataMemberWithClub[] = response.data.Members || [];
            
            // Transform to Club array with role information
            const transformedClubs = members.map((m) => ({
                id: m.Club.ID,
                name: m.Club.Name,
                description: m.Club.Description,
                created_at: m.Club.CreatedAt,
                deleted: m.Club.Deleted,
                user_role: m.Role, // Role comes from the Member entity
                // Map expanded Teams
                user_teams: m.Club.Teams?.map((team) => ({
                    id: team.ID,
                    name: team.Name,
                    description: team.Description,
                    createdAt: team.CreatedAt,
                    clubId: team.ClubID
                })) || []
            }));
            
            setClubs(transformedClubs);
        } catch (err: Error | unknown) {
            console.error('Error fetching clubs:', err);
            setError(t('clubList.failedToFetch'));
        } finally {
            setLoading(false);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [currentUser]);

    useEffect(() => {
        fetchClubs();
    }, [fetchClubs]);

    const handleClubClick = (clubId: string, clubName: string) => {
        // Add to recent clubs when clicking
        addRecentClub(clubId, clubName);
        navigate(`/clubs/${clubId}`);
    };

    const handleTeamClick = (e: React.MouseEvent, clubId: string, teamId: string) => {
        e.stopPropagation(); // Prevent club card click
        navigate(`/clubs/${clubId}/teams/${teamId}`);
    };

    const adminClubs = clubs.filter(club => club.user_role === 'owner' || club.user_role === 'admin');
    const memberClubs = clubs.filter(club => club.user_role === 'member');

    if (loading) {
        return (
            <Layout title={t('clubList.myClubs')}>
                <div>{t('clubList.loadingClubs')}</div>
            </Layout>
        );
    }

    if (error) {
        return (
            <Layout title={t('clubList.myClubs')}>
                <div className="error">{error}</div>
            </Layout>
        );
    }

    return (
        <Layout title={t('clubList.myClubs')} showRecentClubs={false}>
            <div className="clubs-container">
                {adminClubs.length > 0 && (
                    <div className="clubs-section">
                        <h2>{t('clubList.clubsIManage')}</h2>
                        <div className="clubs-grid">
                            {adminClubs.map(club => (
                                <Card
                                    key={club.id}
                                    variant="light"
                                    padding="lg"
                                    clickable
                                    hover
                                    onClick={() => handleClubClick(club.id, club.name)}
                                    className="club-card"
                                >
                                    <div className="club-header">
                                        <h3>{club.name}</h3>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                    {club.deleted && (
                                        <div className="club-deleted-badge">
                                            {t('clubList.deleted')}
                                        </div>
                                    )}
                                    {club.user_teams && club.user_teams.length > 0 && (
                                        <div className="club-teams-section">
                                            <h4 className="teams-title">{t('clubList.myTeams')}</h4>
                                            <div className="teams-list">
                                                {club.user_teams.map(team => (
                                                    <div 
                                                        key={team.id}
                                                        className="team-badge"
                                                        onClick={(e) => handleTeamClick(e, club.id, team.id)}
                                                        title={team.description}
                                                    >
                                                        {team.name}
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </Card>
                            ))}
                        </div>
                    </div>
                )}

                {memberClubs.length > 0 && (
                    <div className="clubs-section">
                        <h2>{t('clubList.clubsImMemberOf')}</h2>
                        <div className="clubs-grid">
                            {memberClubs.map(club => (
                                <Card
                                    key={club.id}
                                    variant="light"
                                    padding="lg"
                                    clickable
                                    hover
                                    onClick={() => handleClubClick(club.id, club.name)}
                                    className="club-card"
                                >
                                    <div className="club-header">
                                        <h3>{club.name}</h3>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                    {club.user_teams && club.user_teams.length > 0 && (
                                        <div className="club-teams-section">
                                            <h4 className="teams-title">{t('clubList.myTeams')}</h4>
                                            <div className="teams-list">
                                                {club.user_teams.map(team => (
                                                    <div 
                                                        key={team.id}
                                                        className="team-badge"
                                                        onClick={(e) => handleTeamClick(e, club.id, team.id)}
                                                        title={team.description}
                                                    >
                                                        {team.name}
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </Card>
                            ))}
                        </div>
                    </div>
                )}

                {adminClubs.length === 0 && memberClubs.length === 0 && (
                    <div className="empty-state">
                        <h2>{t('clubList.noClubsYet')}</h2>
                        <p>{t('clubList.notMemberYet')}</p>
                        <Button 
                            onClick={() => navigate('/createClub')}
                            variant="primary"
                        >
                            {t('clubList.createFirstClub')}
                        </Button>
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default ClubList;