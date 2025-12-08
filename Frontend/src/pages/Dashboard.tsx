import { useNavigate } from 'react-router-dom';
import { useDashboardData, ActivityItem } from '../hooks/useDashboardData';
import { useCurrentUser } from '../hooks/useCurrentUser';
import { useT } from '../hooks/useTranslation';
import Layout from '../components/layout/Layout';
import { Button } from '../components/ui';
import { addRecentClub } from '../utils/recentClubs';
import './Dashboard.css';
import '../styles/events.css';

const Dashboard = () => {
    const { t } = useT();
    const navigate = useNavigate();
    const { activities, loading: dashboardLoading, error: dashboardError } = useDashboardData();
    const { user: currentUser } = useCurrentUser();

    const translateRole = (role: string | undefined): string => {
        if (!role) return 'Unknown Role';
        return t(`clubs.roles.${role}`) || role;
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
        // Determine the message based on actor_name
        return activity.actor_name 
            ? `by ${activity.actor_name}` 
            : '';
    };

    const getPersonalizedTitle = (activity: ActivityItem) => {
        // For role change activities, generate title based on type and personalization
        if (activity.type === 'member_promoted' || activity.type === 'member_demoted' || activity.type === 'role_changed') {
            if (!currentUser || !activity.metadata?.affected_user_id) {
                // Fallback to generic titles
                if (activity.type === 'member_promoted') return 'Member promoted';
                if (activity.type === 'member_demoted') return 'Member demoted';
                return 'Role changed';
            }

            const isCurrentUser = activity.metadata.affected_user_id === currentUser.ID;
            const { new_role } = activity.metadata;
            
            if (isCurrentUser) {
                if (activity.type === 'member_promoted') {
                    return `You got promoted to ${translateRole(new_role)}!`;
                } else if (activity.type === 'member_demoted') {
                    return `Your role changed to ${translateRole(new_role)}`;
                } else {
                    return `Your role changed to ${translateRole(new_role)}`;
                }
            } else {
                // For other users
                if (activity.type === 'member_promoted') return 'Member promoted';
                if (activity.type === 'member_demoted') return 'Member demoted';
                return 'Role changed';
            }
        }
        
        // For non-role activities, return the stored title or generate a default
        return activity.title || 'Activity';
    };

    const getPersonalizedContent = (activity: ActivityItem) => {
        // For role change activities, generate content based on type and personalization
        if (activity.type === 'member_promoted' || activity.type === 'member_demoted' || activity.type === 'role_changed') {
            if (!currentUser || !activity.metadata?.affected_user_id) {
                // Fallback to generic content
                const { old_role, new_role } = activity.metadata || {};
                return `Role changed from ${translateRole(old_role)} to ${translateRole(new_role)}`;
            }

            const isCurrentUser = activity.metadata.affected_user_id === currentUser.ID;
            const { old_role, new_role } = activity.metadata;
            
            if (isCurrentUser) {
                if (activity.actor_name) {
                    return `${activity.actor_name} changed your role from ${translateRole(old_role)} to ${translateRole(new_role)}.`;
                } else {
                    return `Your role was changed from ${translateRole(old_role)} to ${translateRole(new_role)}.`;
                }
            } else {
                // For other users
                return `Role changed from ${translateRole(old_role)} to ${translateRole(new_role)}`;
            }
        }
        
        // For non-role activities, return the stored content
        return activity.content;
    };

    const handleEventClick = (activity: ActivityItem) => {
        if (activity.type === 'event') {
            addRecentClub(activity.club_id, activity.club_name);
            navigate(`/clubs/${activity.club_id}/events/${activity.id}`);
        }
    };

    const renderActivityContent = (activity: ActivityItem) => {
        if (activity.type === 'event') {
            return (
                <div className="event-activity">
                    {(() => {
                        const personalizedContent = getPersonalizedContent(activity);
                        return personalizedContent && (
                            <p className="activity-content">{personalizedContent}</p>
                        );
                    })()}
                    {activity.metadata?.user_rsvp && (
                        <p className="user-rsvp">
                            <strong>Your RSVP:</strong> 
                            <span className={`rsvp-status ${(activity.metadata.user_rsvp as { response: string }).response}`}>
                                {(activity.metadata.user_rsvp as { response: string }).response === 'yes' ? ' Yes' : ' No'}
                            </span>
                        </p>
                    )}
                </div>
            );
        } else if (activity.type === 'news') {
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
                                        <div key={`${activity.type}-${activity.id}`} className="activity-item">
                                            <div className="activity-header">
                                                {activity.type === 'event' ? (
                                                    <Button 
                                                        size="sm"
                                                        variant="secondary"
                                                        className="activity-type-badge clickable-badge"
                                                        onClick={() => handleEventClick(activity)}
                                                    >
                                                        {activity.type.replace(/_/g, ' ')}
                                                    </Button>
                                                ) : (
                                                    <div className="activity-type-badge">
                                                        {activity.type.replace(/_/g, ' ')}
                                                    </div>
                                                )}
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => handleClubClick(activity.club_id, activity.club_name)}
                                                >
                                                    {activity.club_name}
                                                </span>
                                            </div>
                                            <h4 className="activity-title">{getPersonalizedTitle(activity)}</h4>
                                            {renderActivityContent(activity)}
                                            <small className="activity-meta">
                                                {(activity.type === 'member_promoted' || activity.type === 'member_demoted' || activity.type === 'role_changed') ? (
                                                    <>
                                                        {getRoleChangeMessage(activity) && `${getRoleChangeMessage(activity)} • `}
                                                    </>
                                                ) : (
                                                    activity.actor_name && activity.type !== 'event' && (
                                                        <>
                                                            {`Created by ${activity.actor_name}`} • 
                                                        </>
                                                    )
                                                )}
                                                Posted on {formatDateTime(activity.created_at)}
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