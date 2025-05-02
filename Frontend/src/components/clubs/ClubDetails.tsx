import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../layout/Layout';

interface Club {
    id: string;
    name: string;
    description: string;
}

interface Events {
    id: string;
    name: string;
    date: string;
    description: string;
    begin_time: string;
    end_time: string;
}

const ClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const [events, setEvents] = useState<Events[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [isAdmin, setIsAdmin] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [clubResponse, eventsResponse, adminResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}`),
                    api.get(`/api/v1/clubs/${id}/events`),
                    api.get(`/api/v1/clubs/${id}/isAdmin`)
                ]);
                setClub(clubResponse.data);
                setEvents(eventsResponse.data);
                setIsAdmin(adminResponse.data.isAdmin);
                setLoading(false);
            } catch {
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchData();
    }, [id]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <Layout title={club.name}>
            <div>
                <h2>{club.name}</h2>
                <div className="club-info">
                    <p>{club.description}</p>
                    <h3>Events</h3>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Date</th>
                                <th>Description</th>
                                <th>Begin Time</th>
                                <th>End Time</th>
                            </tr>
                        </thead>
                        <tbody>
                            {events.map((event) => (
                                <tr key={event.id}>
                                    <td>{event.name}</td>
                                    <td>{event.date}</td>
                                    <td>{event.description}</td>
                                    <td>{event.begin_time}</td>
                                    <td>{event.end_time}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    {isAdmin && (
                        <button 
                            className="button"
                            onClick={() => navigate(`/clubs/${id}/admin`)}
                            style={{ marginTop: '20px' }}
                        >
                            Manage Club
                        </button>
                    )}
                </div>
            </div>
        </Layout>
    );
};

export default ClubDetails;