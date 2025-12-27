import { useState, useEffect } from 'react';
import Layout from '../../components/layout/Layout';
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from '../../hooks/useAuth';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import { useT } from '../../hooks/useTranslation';
import { Card } from '@/components/ui';
import { buildODataQuery, odataExpandWithOptions, ODataFilter, parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';

interface ODataShiftMember { 
    ID: string; 
    UserID: string;
    Shift?: { 
        ID: string; 
        StartTime: string; 
        EndTime: string; 
        Event?: { 
            ID: string; 
            Name: string; 
            Location: string; 
            Club?: { 
                ID: string; 
                Name: string; 
            }; 
        }; 
    }; 
}

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
    const { user: currentUser } = useCurrentUser();
    const { t } = useT();
    const [shifts, setShifts] = useState<UserShift[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchShifts = async () => {
            if (!currentUser?.ID) {
                setLoading(false);
                return;
            }
            
            try {
                setLoading(true);
                // OData v2: Query ShiftMembers with nested expansions for current user's future shifts
                const now = new Date().toISOString();
                
                // Build nested expand with select clauses for optimal payload
                // Filter by current user's ID to only show their shifts
                const query = buildODataQuery({
                    select: ['ID', 'UserID'],
                    expand: odataExpandWithOptions('Shift', {
                        select: ['ID', 'StartTime', 'EndTime'],
                        expand: odataExpandWithOptions('Event', {
                            select: ['ID', 'Name', 'Location'],
                            expand: odataExpandWithOptions('Club', {
                                select: ['ID', 'Name']
                            })
                        })
                    }),
                    filter: ODataFilter.and(
                        ODataFilter.eq('UserID', currentUser.ID),
                        ODataFilter.gt('Shift/StartTime', now)
                    ),
                    orderby: 'Shift/StartTime'
                });
                
                const response = await api.get<ODataCollectionResponse<ODataShiftMember>>(`/api/v2/ShiftMembers${query}`);
                
                const shiftMembers = parseODataCollection(response.data);
                // Map OData response to match expected format
                const mappedShifts = shiftMembers.map((sm: ODataShiftMember) => ({
                    id: sm.Shift?.ID || sm.ID,
                    startTime: sm.Shift?.StartTime || '',
                    endTime: sm.Shift?.EndTime || '',
                    eventId: sm.Shift?.Event?.ID || '',
                    eventName: sm.Shift?.Event?.Name || 'Unknown Event',
                    location: sm.Shift?.Event?.Location || '',
                    clubId: sm.Shift?.Event?.Club?.ID || '',
                    clubName: sm.Shift?.Event?.Club?.Name || 'Unknown Club',
                    members: [] // Simplified for now
                }));
                setShifts(mappedShifts);
                setError(null);
            } catch (err) {
                console.error('Error fetching shifts:', err);
                setError('Failed to load shifts');
            } finally {
                setLoading(false);
            }
        };

        fetchShifts();
    }, [api, currentUser]);

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
            <ProfileContentLayout title={t('shifts.myFutureShifts')}>
                {loading && (
                    <div style={{ 
                        textAlign: 'center', 
                        padding: '2rem', 
                        color: 'var(--color-text-secondary)' 
                    }}>
                        {t('shifts.loadingShifts')}
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
                        <p>{t('shifts.noUpcomingShifts')}</p>
                        <p style={{ fontSize: '0.9rem', marginTop: '0.5rem' }}>
                            {t('shifts.checkBackLater')}
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
                                                {t('shifts.time')}
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
                                                    {t('shifts.location')}
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
                                                {t('shifts.teamMembers')} ({shift.members.length})
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
                                                                {t('shifts.moreMembers', { count: shift.members.length - 3 })}
                                                            </div>
                                                        )}
                                                    </div>
                                                ) : (
                                                    <span style={{ fontStyle: 'italic' }}>
                                                        {t('shifts.noOtherMembers')}
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