import { useNavigate } from 'react-router-dom';
import { useDashboardData } from '../hooks/useDashboardData';
import Layout from '../components/layout/Layout';

const Dashboard = () => {
    const navigate = useNavigate();
    const { activities, loading: dashboardLoading, error: dashboardError } = useDashboardData();

    const formatDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            return dateTime.toLocaleDateString();
        } catch {
            return timestamp;
        }
    };

    return (
        <Layout title="Dashboard">
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
                                                <div className="activity-type-badge">{activity.type}</div>
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => navigate(`/clubs/${activity.club_id}`)}
                                                >
                                                    {activity.club_name}
                                                </span>
                                            </div>
                                            <h4 className="activity-title">{activity.title}</h4>
                                            {activity.content && (
                                                <p className="activity-content">{activity.content}</p>
                                            )}
                                            <small className="activity-meta">
                                                {activity.creator_name && (
                                                    <>Created by {activity.creator_name} • </>
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