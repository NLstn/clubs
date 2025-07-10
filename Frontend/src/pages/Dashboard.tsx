import { useNavigate } from 'react-router-dom';
import { useDashboardData } from '../hooks/useDashboardData';
import Layout from '../components/layout/Layout';
import { addRecentClub } from '../utils/recentClubs';

const Dashboard = () => {
    const navigate = useNavigate();
    const { activities, loading: dashboardLoading, error: dashboardError } = useDashboardData();

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
                                                <div className="activity-type-badge">{activity.type.replace('_', ' ')}</div>
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => handleClubClick(activity.club_id, activity.club_name)}
                                                >
                                                    {activity.club_name}
                                                </span>
                                            </div>
                                            <h4 className="activity-title">{activity.title}</h4>
                                            {activity.content && (
                                                <p className="activity-content">{activity.content}</p>
                                            )}
                                            <small className="activity-meta">
                                                {activity.actor_name && (
                                                    <>
                                                        {(activity.type === 'member_promoted' || activity.type === 'member_demoted' || activity.type === 'role_changed') 
                                                            ? `Promoted by ${activity.actor_name}` 
                                                            : `Created by ${activity.actor_name}`
                                                        } • 
                                                    </>
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