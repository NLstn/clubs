import { useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { useCurrentUser } from '../hooks/useCurrentUser';
import { useT } from '../hooks/useTranslation';
import Layout from '../components/layout/Layout';
import { Button } from '../components/ui';
import { addRecentClub } from '../utils/recentClubs';
import './Dashboard.css';
import '../styles/events.css';

// ActivityItem represents timeline entries from the backend
export interface ActivityItem {
    ID: string;
    Type: string; // "news", "event", "role_changed", "member_promoted", "member_demoted"
    Title: string;
    Content?: string;
    ClubName: string;
    ClubID: string;
    CreatedAt: string;
    UpdatedAt: string;
    Actor?: string;        // User ID who created/initiated the activity
    ActorName?: string;   // Name of the user who created/initiated the activity
    Metadata?: {
        start_time?: string;
        end_time?: string;
        user_rsvp?: {
            response: string;
        };
        old_role?: string;
        new_role?: string;
        club_name?: string;
        affected_user_id?: string; // User ID of the person whose role was changed
        [key: string]: unknown;
    }; // For extensibility
}

// TimelineItem represents the OData response format
interface TimelineItem {
    ID: string;
    ClubID: string;
    ClubName: string;
    Type: string;
    Title: string;
    Content?: string;
    Timestamp: string;
    CreatedAt: string;
    UpdatedAt: string;
    StartTime?: string;
    EndTime?: string;
    Location?: string;
    Actor?: string;
    ActorName?: string;
    Metadata?: {
        [key: string]: unknown;
    };
    UserRSVP?: {
        ID: string;
        EventID: string;
        UserID: string;
        Response: string;
        CreatedAt: string;
        UpdatedAt: string;
    };
}

// Convert TimelineItem to ActivityItem
function convertToActivity(item: TimelineItem): ActivityItem {
    return {
        ID: item.ID,
        Type: item.Type,
        Title: item.Title,
        Content: item.Content,
        ClubName: item.ClubName,
        ClubID: item.ClubID,
        CreatedAt: item.CreatedAt,
        UpdatedAt: item.UpdatedAt,
        Actor: item.Actor,
        ActorName: item.ActorName,
        Metadata: {
            ...item.Metadata,
            start_time: item.StartTime,
            end_time: item.EndTime,
            user_rsvp: item.UserRSVP ? {
                response: item.UserRSVP.Response
            } : undefined,
        }
    };
}

const Dashboard = () => {
    const { t } = useT();
    const navigate = useNavigate();
    const { api } = useAuth();
    const { user: currentUser } = useCurrentUser();
    const [activities, setActivities] = useState<ActivityItem[]>([]);
    const [dashboardLoading, setDashboardLoading] = useState(true);
    const [dashboardError, setDashboardError] = useState<string | null>(null);

    useEffect(() => {
        const fetchDashboardData = async () => {
            setDashboardLoading(true);
            setDashboardError(null);
            
            try {
                const response = await api.get('/api/v2/TimelineItems');
                const timelineData = response.data.value || [];
                const activitiesData = timelineData.map((item: TimelineItem) => convertToActivity(item));
                setActivities(activitiesData);
            } catch (error) {
                console.error('Error fetching dashboard data:', error);
                setDashboardError('Failed to load dashboard data');
                setActivities([]);
            } finally {
                setDashboardLoading(false);
            }
        };

        fetchDashboardData();
    }, [api]);

    const translateRole = (role: string | undefined): string => {
        if (!role) return t('common.unknownRole');
        return t(`clubs.roles.${role}`);
    };

    const handleClubClick = (clubId: string, clubName: string) => {
        addRecentClub(clubId, clubName);
        navigate(`/clubs/${clubId}`);
    };

    const formatDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            return dateTime.toLocaleDateString();
        } catch {
            return timestamp;
        }
    };

    const getRoleChangeMessage = (activity: ActivityItem) => {
        // Determine the message based on ActorName
        return activity.ActorName 
            ? `by ${activity.ActorName}` 
            : '';
    };

    const getPersonalizedTitle = (activity: ActivityItem) => {
        // For role change activities, generate title based on type and personalization
        if (activity.Type === 'member_promoted' || activity.Type === 'member_demoted' || activity.Type === 'role_changed') {
            if (!currentUser || !activity.Metadata?.affected_user_id) {
                // Fallback to generic titles
                if (activity.Type === 'member_promoted') return 'Member promoted';
                if (activity.Type === 'member_demoted') return 'Member demoted';
                return 'Role changed';
            }

            const isCurrentUser = activity.Metadata.affected_user_id === currentUser.ID;
            const { new_role } = activity.Metadata;
            
            if (isCurrentUser) {
                if (activity.Type === 'member_promoted') {
                    return `You got promoted to ${translateRole(new_role)}!`;
                } else if (activity.Type === 'member_demoted') {
                    return `Your role changed to ${translateRole(new_role)}`;
                } else {
                    return `Your role changed to ${translateRole(new_role)}`;
                }
            } else {
                // For other users
                if (activity.Type === 'member_promoted') return 'Member promoted';
                if (activity.Type === 'member_demoted') return 'Member demoted';
                return 'Role changed';
            }
        }
        
        // For non-role activities, return the stored title or generate a default
        return activity.Title || 'Activity';
    };

    const getPersonalizedContent = (activity: ActivityItem) => {
        // For role change activities, generate content based on type and personalization
        if (activity.Type === 'member_promoted' || activity.Type === 'member_demoted' || activity.Type === 'role_changed') {
            if (!currentUser || !activity.Metadata?.affected_user_id) {
                // Fallback to generic content
                const { old_role, new_role } = activity.Metadata || {};
                return `Role changed from ${translateRole(old_role)} to ${translateRole(new_role)}`;
            }

            const isCurrentUser = activity.Metadata.affected_user_id === currentUser.ID;
            const { old_role, new_role } = activity.Metadata;
            
            if (isCurrentUser) {
                if (activity.ActorName) {
                    return `${activity.ActorName} changed your role from ${translateRole(old_role)} to ${translateRole(new_role)}.`;
                } else {
                    return `Your role was changed from ${translateRole(old_role)} to ${translateRole(new_role)}.`;
                }
            } else {
                // For other users
                return `Role changed from ${translateRole(old_role)} to ${translateRole(new_role)}`;
            }
        }
        
        // For non-role activities, return the stored content
        return activity.Content;
    };

    const handleEventClick = (activity: ActivityItem) => {
        if (activity.Type === 'event') {
            addRecentClub(activity.ClubID, activity.ClubName);
            navigate(`/clubs/${activity.ClubID}/events/${activity.ID}`);
        }
    };

    const renderActivityContent = (activity: ActivityItem) => {
        if (activity.Type === 'event') {
            return (
                <div className="event-activity">
                    {(() => {
                        const personalizedContent = getPersonalizedContent(activity);
                        return personalizedContent && (
                            <p className="activity-content">{personalizedContent}</p>
                        );
                    })()}
                    {activity.Metadata?.user_rsvp && (
                        <p className="user-rsvp">
                            <strong>Your RSVP:</strong> 
                            <span className={`rsvp-status ${(activity.Metadata.user_rsvp as { response: string }).response}`}>
                                {(activity.Metadata.user_rsvp as { response: string }).response === 'yes' ? ' Yes' : ' No'}
                            </span>
                        </p>
                    )}
                </div>
            );
        } else if (activity.Type === 'news') {
            return (
                <div className="news-activity">
                    {(() => {
                        const personalizedContent = getPersonalizedContent(activity);
                        return personalizedContent && (
                            <p className="activity-content">{personalizedContent}</p>
                        );
                    })()}
                </div>
            );
        } else {
            // Regular content for non-event, non-news activities
            const personalizedContent = getPersonalizedContent(activity);
            return personalizedContent && (
                <p className="activity-content">{personalizedContent}</p>
            );
        }
    };

    return (
        <Layout title="Dashboard" showRecentClubs={true}>
            <div>
                {dashboardError && <p className="error">{dashboardError}</p>}

                {/* Activity Feed Section */}
                {dashboardLoading ? (
                    <div>Loading dashboard...</div>
                ) : (
                    <>
                        <div className="dashboard-section">
                            <h2>Activity Feed</h2>
                            {activities.length > 0 ? (
                                <div className="activity-feed">
                                    {activities.map(activity => (
                                        <div key={`${activity.Type}-${activity.ID}`} className="activity-item">
                                            <div className="activity-header">
                                                {activity.Type === 'event' ? (
                                                    <Button 
                                                        size="sm"
                                                        variant="secondary"
                                                        className="activity-type-badge clickable-badge"
                                                        onClick={() => handleEventClick(activity)}
                                                    >
                                                        {activity.Type.replace(/_/g, ' ')}
                                                    </Button>
                                                ) : (
                                                    <div className="activity-type-badge">
                                                        {activity.Type.replace(/_/g, ' ')}
                                                    </div>
                                                )}
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => handleClubClick(activity.ClubID, activity.ClubName)}
                                                >
                                                    {activity.ClubName}
                                                </span>
                                            </div>
                                            <h4 className="activity-title">{getPersonalizedTitle(activity)}</h4>
                                            {renderActivityContent(activity)}
                                            <small className="activity-meta">
                                                {(activity.Type === 'member_promoted' || activity.Type === 'member_demoted' || activity.Type === 'role_changed') ? (
                                                    <>
                                                        {getRoleChangeMessage(activity) && `${getRoleChangeMessage(activity)} • `}
                                                    </>
                                                ) : (
                                                    activity.ActorName && activity.Type !== 'event' && (
                                                        <>
                                                            {`Created by ${activity.ActorName}`} • 
                                                        </>
                                                    )
                                                )}
                                                Posted on {formatDateTime(activity.CreatedAt)}
                                            </small>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className="empty-state">
                                    <p>No recent activities from your clubs.</p>
                                </div>
                            )}
                        </div>
                    </>
                )}
            </div>
        </Layout>
    );
};

export default Dashboard;