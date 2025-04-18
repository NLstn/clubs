import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../layout/Layout';
import InviteMember from './InviteMember';

interface Member {
    id: string;
    name: string;
    email: string;
}

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
    const [club, setClub] = useState<Club | null>(null);
    const [members, setMembers] = useState<Member[]>([]);
    const [events, setEvents] = useState<Events[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [newMember, setNewMember] = useState({ name: '', email: '' });
    const [isModalOpen, setIsModalOpen] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [clubResponse, membersResponse, eventsResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}`),
                    api.get(`/api/v1/clubs/${id}/members`),
                    api.get(`/api/v1/clubs/${id}/events`)
                ]);
                setClub(clubResponse.data);
                setMembers(membersResponse.data);
                setEvents(eventsResponse.data);
                setLoading(false);
            } catch (error) {
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchData();
    }, [id]);

    const deleteMember = async (memberId: string) => {
        try {
            await api.delete(`/api/v1/clubs/${id}/members/${memberId}`);
            setMembers(members.filter(member => member.id !== memberId));
        } catch (error) {
            setError('Failed to delete member');
        }
    };

    const sendInvite = async (email: string) => {
        try {
            await api.post(`/api/v1/clubs/${id}/joinRequest`, { email });
            setIsModalOpen(false);
        } catch (error) {
            setError('Failed to send invite');
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <Layout title={club.name}>
            <div>
                <h2>{club.name}</h2>
                <div className="club-info">
                    <p>{club.description}</p>
                    <h3>Members</h3>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Email</th>
                            </tr>
                        </thead>
                        <tbody>
                            {members.map((member) => (
                                <tr key={member.id}>
                                    <td>{member.name}</td>
                                    <td>{member.email}</td>
                                    <td className="delete-cell">
                                        <button
                                            onClick={() => deleteMember(member.id)}
                                            className="delete-button"
                                            aria-label="Delete member"
                                        >
                                            ×
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    <div className="add-member-form">
                        <input
                            type="text"
                            value={newMember.name}
                            onChange={(e) => setNewMember({ ...newMember, name: e.target.value })}
                            placeholder="Member name"
                        />
                        <input
                            type="email"
                            value={newMember.email}
                            onChange={(e) => setNewMember({ ...newMember, email: e.target.value })}
                            placeholder="Member email"
                        />
                    </div>
                    <button onClick={() => setIsModalOpen(true)}>Invite Member</button>
                    <InviteMember
                        isOpen={isModalOpen}
                        onClose={() => setIsModalOpen(false)}
                        onSubmit={sendInvite}
                    />
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
                </div>
            </div>
        </Layout>
    );
};

export default ClubDetails;