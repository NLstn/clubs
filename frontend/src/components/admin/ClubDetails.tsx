import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

import './ClubDetails.css';

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

const ClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const [members, setMembers] = useState<Member[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [newMember, setNewMember] = useState({ name: '', email: '' });

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [clubResponse, membersResponse] = await Promise.all([
                    axios.get(`${import.meta.env.VITE_API_HOST}/api/v1/clubs/${id}`),
                    axios.get(`${import.meta.env.VITE_API_HOST}/api/v1/clubs/${id}/members`)
                ]);
                setClub(clubResponse.data);
                setMembers(membersResponse.data);
                setLoading(false);
            } catch (error) {
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchData();
    }, [id]);

    const addMember = async () => {
        try {
            const response = await axios.post(
                `${import.meta.env.VITE_API_HOST}/api/v1/clubs/${id}/members`,
                {
                    name: newMember.name,
                    email: newMember.email
                }
            );
            setMembers([...members, response.data]);
            setNewMember({ name: '', email: '' }); // Reset form
        } catch (error) {
            setError('Failed to add member');
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <div className="club-details">
            <button onClick={() => navigate(-1)} className="back-button">
                ← Back
            </button>
            <h2>{club.name}</h2>
            <div className="club-info">
                <p>{club.description}</p>
                <h3>Members</h3>
                <table className="members-table">
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
                    <button 
                        onClick={addMember}
                        disabled={!newMember.name || !newMember.email}
                    >
                        Add Member
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ClubDetails;