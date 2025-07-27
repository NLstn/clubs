import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../utils/api';
import { useT } from '../../hooks/useTranslation';
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
        return t(`clubs.roles.${role}`) || role;
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

    // Don't render anything if loading, error, or no members due to permission restrictions
    if (loading) return <div className="loading-text">Loading members...</div>;
    if (error) return <div className="error-text">{error}</div>;
    if (members.length === 0) return null;

    return (
        <div className="content-section">
            <h3>{t('clubs.members') || 'Members'}</h3>
            <div className="members-table-container">
                <table className="readonly-members-table">
                    <thead>
                        <tr>
                            <th>{t('common.name') || 'Name'}</th>
                            <th>{t('common.role') || 'Role'}</th>
                            <th>{t('clubs.joined') || 'Joined'}</th>
                            <th>{t('clubs.birthDate') || 'Birth Date'}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {members.map((member) => (
                            <tr key={member.id}>
                                <td>{member.name}</td>
                                <td>
                                    <span className={`role-badge ${member.role.toLowerCase()}`}>
                                        {translateRole(member.role)}
                                    </span>
                                </td>
                                <td>{member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : 'N/A'}</td>
                                <td>{member.birthDate ? new Date(member.birthDate).toLocaleDateString() : 'Not shared'}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                <div className="member-count">
                    {t('clubs.totalMembers', { count: members.length }) || `Total: ${members.length} members`}
                </div>
            </div>
        </div>
    );
};

export default ReadonlyMemberList;
