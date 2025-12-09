import { useState, useEffect } from 'react';
import Layout from '../../components/layout/Layout';
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from '../../hooks/useAuth';
import { useT } from '../../hooks/useTranslation';
import { Card } from '@/components/ui';

interface UserShift {
    id: string;
    startTime: string;
    endTime: string;
    eventId: string;
    eventName: string;
    location: string;
    clubId: string;
    clubName: string;
    members: string[];
}

function ProfileShifts() {
    const { api } = useAuth();
    const { t } = useT();
    const [shifts, setShifts] = useState<UserShift[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchShifts = async () => {
            try {
                setLoading(true);
                const response = await api.get('/api/v1/me/shifts');
                setShifts(response.data || []);
                setError(null);
            } catch (err) {
                console.error('Error fetching shifts:', err);
                setError('Failed to load shifts');
            } finally {
                setLoading(false);
            }
        };

        fetchShifts();
    }, [api]);

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString();
    };

    const formatTime = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    };

    return (
        <Layout>
            <ProfileContentLayout title={t('My Future Shifts')}>
                {loading && (
                    <div style={{ 
                        textAlign: 'center', 
                        padding: '2rem', 
                        color: 'var(--color-text-secondary)' 
                    }}>
                        Loading shifts...
                    </div>
                )}

                {error && (
                    <div style={{ 
                        background: 'var(--color-error-bg)', 
                        color: 'var(--color-error)', 
                        padding: '1rem', 
                        borderRadius: '4px',
                        marginBottom: '1rem'
                    }}>
                        {error}
                    </div>
                )}

                {!loading && !error && shifts.length === 0 && (
                    <div style={{ 
                        textAlign: 'center', 
                        padding: '3rem', 
                        color: 'var(--color-text-secondary)' 
                    }}>
                        <p>{t('No upcoming shifts found.')}</p>
                        <p style={{ fontSize: '0.9rem', marginTop: '0.5rem' }}>
                            {t('Check back later or contact your club administrators if you expect to see shifts here.')}
                        </p>
                    </div>
                )}

                {!loading && !error && shifts.length > 0 && (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            {shifts.map((shift) => (
                                <Card
                                    key={shift.id}
                                    variant="default"
                                    padding="lg"
                                >
                                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1rem' }}>
                                        <div>
                                            <h3 style={{ 
                                                color: 'var(--color-text-primary)', 
                                                margin: '0 0 0.5rem 0',
                                                fontSize: '1.25rem'
                                            }}>
                                                {shift.eventName}
                                            </h3>
                                            <p style={{ 
                                                color: 'var(--color-text-secondary)', 
                                                margin: '0',
                                                fontSize: '0.9rem'
                                            }}>
                                                {shift.clubName}
                                            </p>
                                        </div>
                                        <div style={{ textAlign: 'right' }}>
                                            <div style={{ 
                                                background: 'var(--color-primary)', 
                                                color: 'white', 
                                                padding: '0.25rem 0.75rem',
                                                borderRadius: '4px',
                                                fontSize: '0.85rem',
                                                fontWeight: '500'
                                            }}>
                                                {formatDate(shift.startTime)}
                                            </div>
                                        </div>
                                    </div>

                                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem' }}>
                                        <div>
                                            <h4 style={{ 
                                                color: 'var(--color-text-primary)', 
                                                margin: '0 0 0.5rem 0',
                                                fontSize: '0.9rem',
                                                fontWeight: '600'
                                            }}>
                                                {t('Time')}
                                            </h4>
                                            <p style={{ 
                                                color: 'var(--color-text-secondary)', 
                                                margin: '0',
                                                fontSize: '0.9rem'
                                            }}>
                                                {formatTime(shift.startTime)} - {formatTime(shift.endTime)}
                                            </p>
                                        </div>

                                        {shift.location && (
                                            <div>
                                                <h4 style={{ 
                                                    color: 'var(--color-text-primary)', 
                                                    margin: '0 0 0.5rem 0',
                                                    fontSize: '0.9rem',
                                                    fontWeight: '600'
                                                }}>
                                                    {t('Location')}
                                                </h4>
                                                <p style={{ 
                                                    color: 'var(--color-text-secondary)', 
                                                    margin: '0',
                                                    fontSize: '0.9rem'
                                                }}>
                                                    {shift.location}
                                                </p>
                                            </div>
                                        )}

                                        <div>
                                            <h4 style={{ 
                                                color: 'var(--color-text-primary)', 
                                                margin: '0 0 0.5rem 0',
                                                fontSize: '0.9rem',
                                                fontWeight: '600'
                                            }}>
                                                {t('Team Members')} ({shift.members.length})
                                            </h4>
                                            <div style={{ 
                                                color: 'var(--color-text-secondary)', 
                                                fontSize: '0.9rem'
                                            }}>
                                                {shift.members.length > 0 ? (
                                                    <div>
                                                        {shift.members.slice(0, 3).map((member, index) => (
                                                            <div key={index} style={{ marginBottom: '0.25rem' }}>
                                                                {member}
                                                            </div>
                                                        ))}
                                                        {shift.members.length > 3 && (
                                                            <div style={{ 
                                                                fontStyle: 'italic',
                                                                color: 'var(--color-text-tertiary)'
                                                            }}>
                                                                +{shift.members.length - 3} more
                                                            </div>
                                                        )}
                                                    </div>
                                                ) : (
                                                    <span style={{ fontStyle: 'italic' }}>
                                                        {t('No other members assigned')}
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                </Card>
                            ))}
                        </div>
                    )}
            </ProfileContentLayout>
        </Layout>
    );
}

export default ProfileShifts;