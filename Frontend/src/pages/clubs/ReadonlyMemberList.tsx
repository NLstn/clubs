import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../utils/api';
import { useT } from '../../hooks/useTranslation';
import { Table, TableColumn } from '@/components/ui';
import './ReadonlyMemberList.css';

interface Member {
    id: string;
    name: string;
    role: string;
    joinedAt: string;
    userId?: string;
    birthDate?: string;
}

const ReadonlyMemberList = () => {
    const { t } = useT();
    const { id: clubId } = useParams();
    const [members, setMembers] = useState<Member[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const translateRole = (role: string): string => {
        return t(`clubs.roles.${role}`);
    };

    const sortMembersByRole = (members: Member[]): Member[] => {
        const roleOrder: { [key: string]: number } = { 
            'owner': 0, 
            'admin': 1, 
            'member': 2 
        };
        
        return [...members].sort((a, b) => {
            const aOrder = roleOrder[a.role.toLowerCase()] ?? 999;
            const bOrder = roleOrder[b.role.toLowerCase()] ?? 999;
            return aOrder - bOrder;
        });
    };

    useEffect(() => {
        const fetchMembers = async () => {
            if (!clubId) {
                setLoading(false);
                return;
            }

            try {
                const response = await api.get(`/api/v1/clubs/${clubId}/members`);
                const sortedMembers = sortMembersByRole(response.data);
                setMembers(sortedMembers);
                setError(null);
            } catch (err) {
                console.error('Error fetching members:', err);
                // If it's a 403 error, the member list might be disabled for regular members
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 403) {
                        // Don't show anything if member list is disabled
                        setMembers([]);
                        setError(null);
                        setLoading(false);
                        return;
                    }
                }
                setError('Failed to fetch members');
            } finally {
                setLoading(false);
            }
        };

        fetchMembers();
    }, [clubId]);

    // Define table columns
    const columns: TableColumn<Member>[] = [
        {
            key: 'name',
            header: t('common.name'),
            render: (member) => <span>{member.name}</span>
        },
        {
            key: 'role',
            header: t('common.role'),
            render: (member) => (
                <span className={`role-badge ${member.role.toLowerCase()}`}>
                    {translateRole(member.role)}
                </span>
            )
        },
        {
            key: 'joinedAt',
            header: t('clubs.joined'),
            render: (member) => (
                <span>{member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : 'N/A'}</span>
            ),
            className: 'hide-mobile'
        },
        {
            key: 'birthDate',
            header: t('clubs.birthDate'),
            render: (member) => (
                <span>{member.birthDate ? new Date(member.birthDate).toLocaleDateString() : 'Not shared'}</span>
            ),
            className: 'hide-small'
        }
    ];

    // Don't render anything if no members due to permission restrictions and not loading/error
    if (!loading && !error && members.length === 0) return null;

    return (
        <div className="content-section">
            <h3>{t('clubs.members')}</h3>
            <Table
                columns={columns}
                data={members}
                keyExtractor={(member) => member.id}
                loading={loading}
                error={error}
                emptyMessage="No members available"
                loadingMessage="Loading members..."
                errorMessage="Failed to load members"
                footer={
                    members.length > 0 ? (
                        <div>
                            {t('clubs.totalMembers', { count: members.length })}
                        </div>
                    ) : null
                }
            />
        </div>
    );
};

export default ReadonlyMemberList;
