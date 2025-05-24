import { useEffect, useState } from "react";
import api from "../../utils/api";
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";

interface Fine {
    id: string;
    clubName: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
}

const ProfileFines = () => {
    const [fines, setFines] = useState<Fine[]>([]);
    const [error, setError] = useState('');

    useEffect(() => {
        fetchFines();
    }, []);

    const fetchFines = async () => {
        try {
            const response = await api.get('/api/v1/me/fines');
            if (response.status === 200) {
                const data = response.data;
                setFines(data);
            }
        } catch (error) {
            setError('Error fetching fines: ' + error);
        }
    };

    return (
        <Layout title="Club Invitations">
            <div style={{
                display: 'flex',
                minHeight: 'calc(100vh - 90px)',
                width: '100%',
                position: 'relative'
            }}>
                <ProfileSidebar />
                <div style={{
                    flex: '1 1 auto',
                    padding: '20px',
                    maxWidth: 'calc(100% - 200px)'
                }}>
                    <h1>Fines</h1>
                    {error && <div className="error">{error}</div>}
                    <table>
                        <thead>
                            <tr>
                                <th>Club Name</th>
                                <th>Amount</th>
                                <th>Status</th>
                                <th>Date</th>
                                <th>Reason</th>
                            </tr>
                        </thead>
                        <tbody>
                            {fines.map((fine) => (
                                <tr key={fine.id}>
                                    <td>{fine.clubName}</td>
                                    <td>{fine.amount}</td>
                                    <td>{new Date(fine.createdAt).toLocaleDateString()}</td>
                                    <td>{fine.reason}</td>
                                    <td>{fine.paid ? 'Paid' : 'Unpaid'}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </Layout>
    )
}

export default ProfileFines;