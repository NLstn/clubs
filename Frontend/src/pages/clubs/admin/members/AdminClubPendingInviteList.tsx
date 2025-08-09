import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Table, TableColumn } from '@/components/ui';
import api from "../../../../utils/api";

interface Invite {
    id: string;
    email: string;
}

const AdminClubPendingInviteList = () => {

    const { id } = useParams();
    const [invites, setInvites] = useState<Invite[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Define table columns
    const columns: TableColumn<Invite>[] = [
        {
            key: 'email',
            header: 'Email',
            render: (invite) => invite.email
        }
    ];

    useEffect(() => {
        const fetchInvites = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await api.get(`/api/v1/clubs/${id}/invites`);
                setInvites(response.data);
            } catch (error) {
                console.error("Error fetching invites:", error);
                setError("Failed to load pending invites");
            } finally {
                setLoading(false);
            }
        };
        fetchInvites();
    }, [id]);

    return (
        <div>
            <Table
                columns={columns}
                data={invites}
                keyExtractor={(invite) => invite.id}
                loading={loading}
                error={error}
                emptyMessage="No pending invites found"
                loadingMessage="Loading pending invites..."
                errorMessage="Failed to load pending invites"
            />
        </div>
    )
}

export default AdminClubPendingInviteList;