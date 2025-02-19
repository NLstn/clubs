import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

interface Club {
    id: number;
    name: string;
    description: string;
}

const ClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchClubDetails = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/api/v1/clubs/${id}`);
                setClub(response.data);
                setLoading(false);
            } catch (error) {
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchClubDetails();
    }, [id]);

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
            </div>
        </div>
    );
};

export default ClubDetails;